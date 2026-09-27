package servers

import (
	"encoding/json"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/config"
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

// Load buttons, preferring user buttons over defaults. The defaults go
// through siteVariantPath for the same reason the config server's GET does:
// on a bedroom site the generic default-buttons.json is a different button
// set from the one the UI is showing, so matching against it would answer
// about lights this site does not have.
func loadButtons(uiConfigDir string) ([]animations.Button, error) {

	content, err := os.ReadFile(filepath.Join(uiConfigDir, "user-buttons.json"))
	if err != nil {
		if content, err = os.ReadFile(siteVariantPath(uiConfigDir, "default-buttons.json")); err != nil {
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

		names, _, err := buttonNames(uiConfigDir)
		if err != nil {
			slog.Error("Buttons failed to load buttons", "err", err.Error())
			http.Error(w, "failed to load buttons", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(names)
		if _, err = w.Write(response); err != nil {
			slog.Warn("Buttons failed to write response", "err", err.Error())
		}
	}
}

// buttonNames returns the speakable button names and those that cannot be
// reached by voice. A blank name has nothing to say, and two names sharing a
// canonical form are indistinguishable to findButtonByName, which returns the
// first — so the page says so rather than letting one button look broken.
func buttonNames(uiConfigDir string) (names []string, unreachable []string, err error) {

	buttons, err := loadButtons(uiConfigDir)
	if err != nil {
		return nil, nil, err
	}

	seen := make(map[string]string, len(buttons))
	for _, b := range buttons {
		if b.Name == "" {
			continue
		}
		names = append(names, b.Name)
		canonical := canonicaliseButtonName(b.Name)
		if first, clash := seen[canonical]; clash {
			unreachable = append(unreachable, b.Name+" (sounds like "+first+")")
			continue
		}
		seen[canonical] = b.Name
	}
	return names, unreachable, nil
}

// The developer console's slot-value bulk editor takes one value per line, so
// that is what the textarea holds. Kept deliberately plain: this is served by
// the LAN server and read on a phone while the console is open on a laptop.
var alexaButtonsPage = template.Must(template.New("alexaButtons").Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Alexa slot values</title>
<style>
 :root { color-scheme: light dark }
 body { font: 16px system-ui, sans-serif; margin: 0; padding: 20px; max-width: 46rem }
 h1 { font-size: 1.25rem; margin: 0 0 .25rem }
 p { margin: .5rem 0; line-height: 1.45 }
 .muted { opacity: .7; font-size: .9rem }
 textarea { width: 100%; box-sizing: border-box; min-height: 20rem; margin-top: .75rem;
            font: 14px ui-monospace, monospace; padding: .6rem }
 button { font: inherit; padding: .5rem 1rem; cursor: pointer }
 ol { padding-left: 1.4rem } li { margin: .3rem 0 }
 .warn { border-left: 3px solid #c60; padding: .1rem 0 .1rem .7rem; margin: 1rem 0 }
</style></head><body>
<h1>Alexa slot values</h1>
<p class="muted">{{len .Names}} button names, site <strong>{{.Site}}</strong>, from
{{.Source}}. This page is LAN-only and is never reachable from the internet.</p>

<button id="copy">Copy all</button> <span id="done" class="muted"></span>
<textarea id="values" readonly spellcheck="false">{{range .Names}}{{.}}
{{end}}</textarea>

{{if .Unreachable}}<div class="warn"><strong>Not reachable by voice</strong>
<p class="muted">These collide with another button once punctuation and case are
ignored, so the first one always wins. Rename one to use both by voice.</p>
<ul>{{range .Unreachable}}<li>{{.}}</li>{{end}}</ul></div>{{end}}

<p><strong>To update the skill</strong></p>
<ol>
 <li>Open the skill in the Alexa developer console, then
     <em>Build &rarr; Interaction Model &rarr; Slot Types &rarr; BUTTON_NAME</em>.</li>
 <li>Choose <em>Bulk Edit</em>, replace everything with the list above, and Save.</li>
 <li>Press <em>Build Model</em> and wait for it to finish — until that completes
     Alexa still resolves against the previous names.</li>
</ol>
<p class="muted">Renaming a button takes effect on this server at once; the slot
list is a separate copy held by Amazon, which is why it has to be pasted again.</p>

<script>
 document.getElementById('copy').addEventListener('click', async () => {
   const area = document.getElementById('values');
   try { await navigator.clipboard.writeText(area.value); }
   catch { area.select(); document.execCommand('copy'); }
   document.getElementById('done').textContent = 'copied';
 });
</script>
</body></html>`))

// GetAlexaButtonsPageHandler Handle GET /api/AlexaButtons.html, the same names
// as the JSON endpoint but as a page with a copy button and the console steps,
// so updating the skill after a rename is paste-once rather than hand-typed.
func GetAlexaButtonsPageHandler(uiConfigDir string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		names, unreachable, err := buttonNames(uiConfigDir)
		if err != nil {
			slog.Error("Buttons page failed to load buttons", "err", err.Error())
			http.Error(w, "failed to load buttons", 500)
			return
		}

		source := "user-buttons.json"
		if _, err = os.Stat(filepath.Join(uiConfigDir, "user-buttons.json")); err != nil {
			source = filepath.Base(siteVariantPath(uiConfigDir, "default-buttons.json"))
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err = alexaButtonsPage.Execute(w, struct {
			Names       []string
			Unreachable []string
			Site        string
			Source      string
		}{names, unreachable, config.Site, source}); err != nil {
			slog.Warn("Buttons page failed to write response", "err", err.Error())
		}
	}
}
