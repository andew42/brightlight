package servers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/andew42/brightlight/animations"
)

// Alexa voice control. The Alexa service calls this endpoint directly (there
// is no Lambda), so it is the one part of brightlight reachable from the
// internet — a reverse proxy terminates TLS and forwards the single route
// below. It gets its own port and mux so nothing else can ever be exposed by
// a proxy misconfiguration, and requests are only acted on when they carry a
// valid Amazon signature and our skill id.
//
//	POST /alexa/skill   <- Alexa request envelope, speech response back
//
// The feature is opt-in and fails closed: without BRIGHTLIGHT_ALEXA_SKILL_ID
// there is nothing to bind requests to our skill, so the listener does not
// start at all.

// StartAlexaServer Start the Alexa listener if configured.
// uiConfigDir is the filesystem path to the ui-config directory.
func StartAlexaServer(uiConfigDir string) {

	skillId := os.Getenv("BRIGHTLIGHT_ALEXA_SKILL_ID")
	if skillId == "" {
		slog.Info("BRIGHTLIGHT_ALEXA_SKILL_ID not set, Alexa endpoint disabled")
		return
	}

	port := os.Getenv("BRIGHTLIGHT_ALEXA_PORT")
	if port == "" {
		port = "8443"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/alexa/skill", getSkillHandler(uiConfigDir, skillId, newAlexaVerifier()))

	go func() {
		slog.Info("serving Alexa endpoint", "port", port, "skillId", skillId)
		err := http.ListenAndServe(":"+port, mux)
		slog.Error("Alexa server exited", "err", err)
	}()
}

// The button file also contains user segments which we ignore here
type buttonFileDef struct {
	Buttons []animations.Button
}

// Load buttons, preferring user buttons over defaults
func loadButtons(uiConfigDir string) ([]animations.Button, error) {

	content, err := os.ReadFile(filepath.Join(uiConfigDir, "user-buttons.json"))
	if err != nil {
		if content, err = os.ReadFile(filepath.Join(uiConfigDir, "default-buttons.json")); err != nil {
			return nil, err
		}
	}

	var buttonFile buttonFileDef
	if err = json.Unmarshal(content, &buttonFile); err != nil {
		return nil, err
	}
	return buttonFile.Buttons, nil
}

// Spoken numbers arrive as words but buttons may be named with digits
var alexaWordDigits = map[string]string{
	"zero": "0", "one": "1", "two": "2", "three": "3", "four": "4",
	"five": "5", "six": "6", "seven": "7", "eight": "8", "nine": "9", "ten": "10",
}

// Reduce a name to a canonical form for matching: lower case, punctuation
// stripped, whitespace collapsed and number words replaced with digits so
// e.g. "Two Tone!" and "2 tone" both become "2 tone"
func canonicaliseButtonName(name string) string {

	mapped := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return ' '
	}, name)

	words := strings.Fields(mapped)
	for i, w := range words {
		if digit, ok := alexaWordDigits[w]; ok {
			words[i] = digit
		}
	}
	return strings.Join(words, " ")
}

// Levenshtein edit distance, used to forgive small speech to text errors
func editDistance(a string, b string) int {

	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(curr[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}

// Find the button best matching a spoken name, exact canonical match first
// then closest edit distance within a small threshold
func findButtonByName(buttons []animations.Button, spoken string) *animations.Button {

	target := canonicaliseButtonName(spoken)
	if target == "" {
		return nil
	}

	for i := range buttons {
		if canonicaliseButtonName(buttons[i].Name) == target {
			return &buttons[i]
		}
	}

	// Allow up to 2 edits but never more than a quarter of the name so short
	// names like "off" don't match "on"
	best := -1
	bestDistance := 3
	for i := range buttons {
		name := canonicaliseButtonName(buttons[i].Name)
		if name == "" {
			continue
		}
		maxAllowed := min(2, len(name)/4)
		if d := editDistance(name, target); d <= maxAllowed && d < bestDistance {
			best, bestDistance = i, d
		}
	}
	if best < 0 {
		return nil
	}
	return &buttons[best]
}

// The parts of an Alexa request envelope we act on. The application id
// appears under session for a spoken interaction and under context for
// requests that arrive outside one, so both are read.
type alexaEnvelope struct {
	Session struct {
		Application struct {
			ApplicationId string `json:"applicationId"`
		} `json:"application"`
	} `json:"session"`
	Context struct {
		System struct {
			Application struct {
				ApplicationId string `json:"applicationId"`
			} `json:"application"`
		} `json:"System"`
	} `json:"context"`
	Request struct {
		Type      string    `json:"type"`
		Timestamp time.Time `json:"timestamp"`
		Intent    struct {
			Name  string `json:"name"`
			Slots map[string]struct {
				Value string `json:"value"`
			} `json:"slots"`
		} `json:"intent"`
	} `json:"request"`
}

func (e alexaEnvelope) applicationId() string {

	if id := e.Session.Application.ApplicationId; id != "" {
		return id
	}
	return e.Context.System.Application.ApplicationId
}

type alexaSpeech struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type alexaResponse struct {
	Version  string `json:"version"`
	Response struct {
		// A pointer so SessionEndedRequest can answer without speech
		OutputSpeech     *alexaSpeech `json:"outputSpeech,omitempty"`
		ShouldEndSession bool         `json:"shouldEndSession"`
	} `json:"response"`
}

func speak(text string, endSession bool) alexaResponse {

	var r alexaResponse
	r.Version = "1.0"
	r.Response.OutputSpeech = &alexaSpeech{Type: "PlainText", Text: text}
	r.Response.ShouldEndSession = endSession
	return r
}

// Handle POST /alexa/skill, the endpoint registered with the Alexa skill
func getSkillHandler(uiConfigDir string, skillId string, verifier *alexaVerifier) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "method not allowed", 405)
			return
		}

		// The signature covers the exact bytes Amazon sent, so the body is
		// read whole and checked before anything is unmarshalled from it
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64*1024))
		if err != nil {
			slog.Warn("Alexa request body unreadable", "err", err.Error())
			http.Error(w, "bad request", 400)
			return
		}

		if err = verifier.verify(r.Header.Get("SignatureCertChainUrl"),
			r.Header.Get("Signature-256"), body); err != nil {
			slog.Warn("Alexa request failed signature check", "remote", r.RemoteAddr, "err", err.Error())
			http.Error(w, "unauthorized", 401)
			return
		}

		var envelope alexaEnvelope
		if err = json.Unmarshal(body, &envelope); err != nil {
			slog.Warn("Alexa request bad body", "err", err.Error())
			http.Error(w, "bad request", 400)
			return
		}

		// A signature only proves the request came from Amazon; anyone can
		// aim their own skill at this URL, so the application id is what
		// says it came from ours
		if envelope.applicationId() != skillId {
			slog.Warn("Alexa request from unexpected skill", "skillId", envelope.applicationId())
			http.Error(w, "unauthorized", 401)
			return
		}

		if err = checkAlexaTimestamp(envelope.Request.Timestamp, verifier.now()); err != nil {
			slog.Warn("Alexa request timestamp rejected", "err", err.Error())
			http.Error(w, "unauthorized", 401)
			return
		}

		response, _ := json.Marshal(handleAlexaRequest(uiConfigDir, envelope))
		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(response); err != nil {
			slog.Warn("Alexa failed to write response", "err", err.Error())
		}
	}
}

// Turn a verified request into the speech Alexa should say back
func handleAlexaRequest(uiConfigDir string, envelope alexaEnvelope) alexaResponse {

	switch envelope.Request.Type {

	case "LaunchRequest":
		return speak("Which light setting would you like?", false)

	case "IntentRequest":
		switch envelope.Request.Intent.Name {

		case "RunButtonIntent":
			spoken := envelope.Request.Intent.Slots["buttonName"].Value
			if spoken == "" {
				return speak("Which light setting would you like?", false)
			}
			return runSpokenButton(uiConfigDir, spoken)

		case "AMAZON.HelpIntent":
			return speak("Say the name of a light setting, for example rainbow.", false)

		default: // AMAZON.StopIntent, AMAZON.CancelIntent etc.
			return speak("Goodbye.", true)
		}

	default: // SessionEndedRequest must not include speech
		var r alexaResponse
		r.Version = "1.0"
		r.Response.ShouldEndSession = true
		return r
	}
}

// Find the button matching what was heard and run it
func runSpokenButton(uiConfigDir string, spoken string) alexaResponse {

	buttons, err := loadButtons(uiConfigDir)
	if err != nil {
		slog.Error("Alexa failed to load buttons", "err", err.Error())
		return speak("I couldn't read the light settings.", true)
	}

	button := findButtonByName(buttons, spoken)
	if button == nil {
		slog.Info("Alexa no matching button", "heard", spoken)
		return speak("I couldn't find a light setting called "+spoken+".", true)
	}
	slog.Info("Alexa run button", "heard", spoken, "matched", button.Name)

	// Run the animation and keep UI button pads in sync
	animations.RunAnimations(button.Segments)
	updateActiveButtonKey(button.Key)
	return speak("OK, "+button.Name+".", true)
}

// GetAlexaButtonsHandler Handle GET /api/AlexaButtons returning a JSON array
// of button names to paste into the skill's slot values. It lives on the LAN
// server rather than the Alexa listener so nothing but the skill route is
// ever reachable from the internet.
func GetAlexaButtonsHandler(uiConfigDir string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		buttons, err := loadButtons(uiConfigDir)
		if err != nil {
			slog.Error("Buttons failed to load buttons", "err", err.Error())
			http.Error(w, "failed to load buttons", 500)
			return
		}

		names := make([]string, 0, len(buttons))
		for _, b := range buttons {
			if b.Name != "" {
				names = append(names, b.Name)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(names)
		if _, err = w.Write(response); err != nil {
			slog.Warn("Buttons failed to write response", "err", err.Error())
		}
	}
}
