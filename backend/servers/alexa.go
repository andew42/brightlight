package servers

import (
	"crypto/subtle"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/andew42/brightlight/animations"
)

// Alexa voice control endpoints, served over HTTPS on a separate port with
// a separate mux so only these two routes are ever exposed through the
// router. The feature is opt-in: if BRIGHTLIGHT_ALEXA_TOKEN is unset the
// listener is not started at all.
//
//	POST /alexa/RunButton {"button":"sweet shop"} -> runs the button
//	GET  /alexa/Buttons                           -> lists button names

// StartAlexaServer Start the Alexa HTTPS listener if configured.
// uiConfigDir is the filesystem path to the ui-config directory.
func StartAlexaServer(uiConfigDir string) {

	token := os.Getenv("BRIGHTLIGHT_ALEXA_TOKEN")
	if token == "" {
		slog.Info("BRIGHTLIGHT_ALEXA_TOKEN not set, Alexa endpoint disabled")
		return
	}
	if len(token) < 32 {
		slog.Warn("BRIGHTLIGHT_ALEXA_TOKEN is short, use at least 32 random characters")
	}

	certPath := os.Getenv("BRIGHTLIGHT_ALEXA_CERT")
	keyPath := os.Getenv("BRIGHTLIGHT_ALEXA_KEY")
	if certPath == "" || keyPath == "" {
		slog.Error("BRIGHTLIGHT_ALEXA_CERT or BRIGHTLIGHT_ALEXA_KEY not set, Alexa endpoint disabled")
		return
	}

	port := os.Getenv("BRIGHTLIGHT_ALEXA_PORT")
	if port == "" {
		port = "8443"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/alexa/RunButton", requireAlexaToken(token, getRunButtonHandler(uiConfigDir)))
	mux.HandleFunc("/alexa/Buttons", requireAlexaToken(token, getButtonsHandler(uiConfigDir)))

	go func() {
		slog.Info("serving Alexa endpoint", "port", port)
		err := http.ListenAndServeTLS(":"+port, certPath, keyPath, mux)
		slog.Error("Alexa server exited", "err", err)
	}()
}

// Reject requests without the expected bearer token
func requireAlexaToken(token string, handler http.HandlerFunc) http.HandlerFunc {

	expected := []byte("Bearer " + token)
	return func(w http.ResponseWriter, r *http.Request) {
		supplied := []byte(r.Header.Get("Authorization"))
		if subtle.ConstantTimeCompare(supplied, expected) != 1 {
			slog.Warn("Alexa request with bad token", "remote", r.RemoteAddr)
			http.Error(w, "unauthorized", 401)
			return
		}
		handler(w, r)
	}
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

// Handle POST /alexa/RunButton {"button":"name"}
func getRunButtonHandler(uiConfigDir string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "method not allowed", 405)
			return
		}

		var request struct {
			Button string
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024)).Decode(&request); err != nil {
			slog.Warn("RunButton bad body", "err", err.Error())
			http.Error(w, "bad request", 400)
			return
		}

		buttons, err := loadButtons(uiConfigDir)
		if err != nil {
			slog.Error("RunButton failed to load buttons", "err", err.Error())
			http.Error(w, "failed to load buttons", 500)
			return
		}

		button := findButtonByName(buttons, request.Button)
		if button == nil {
			slog.Info("RunButton no matching button", "button", request.Button)
			http.Error(w, "button not found", 404)
			return
		}
		slog.Info("RunButton called", "heard", request.Button, "matched", button.Name)

		// Run the animation and keep UI button pads in sync
		animations.RunAnimations(button.Segments)
		updateActiveButtonKey(button.Key)

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(struct{ Matched string }{button.Name})
		if _, err = w.Write(response); err != nil {
			slog.Warn("RunButton failed to write response", "err", err.Error())
		}
	}
}

// Handle GET /alexa/Buttons returning a JSON array of button names
func getButtonsHandler(uiConfigDir string) http.HandlerFunc {

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
