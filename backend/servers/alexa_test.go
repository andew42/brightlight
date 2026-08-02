package servers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andew42/brightlight/animations"
)

func namedButtons(names ...string) []animations.Button {

	buttons := make([]animations.Button, len(names))
	for i, n := range names {
		buttons[i] = animations.Button{Key: i + 1, Name: n}
	}
	return buttons
}

func TestCanonicaliseButtonName(t *testing.T) {

	cases := map[string]string{
		"Sweet Shop": "sweet shop",
		"  OFF ":     "off",
		"Two Tone!":  "2 tone",
		"2 tone":     "2 tone",
		"Baby-Bows":  "baby bows",
		"":           "",
		"Three":      "3",
	}
	for in, want := range cases {
		if got := canonicaliseButtonName(in); got != want {
			t.Errorf("canonicaliseButtonName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFindButtonByName(t *testing.T) {

	buttons := namedButtons("OFF", "Full", "Sweet Shop", "Two Tone", "Rainbow", "One")

	cases := map[string]string{
		"off":        "OFF",
		"sweet shop": "Sweet Shop",
		"Sweet Shop": "Sweet Shop",
		"two tone":   "Two Tone", // spoken "two" matches button "Two"
		"2 tone":     "Two Tone",
		"rainbows":   "Rainbow", // small speech to text error forgiven
		"sweet shot": "Sweet Shop",
		"1":          "One", // digit heard, word named button
	}
	for spoken, want := range cases {
		b := findButtonByName(buttons, spoken)
		if b == nil {
			t.Errorf("findButtonByName(%q) = nil, want %q", spoken, want)
		} else if b.Name != want {
			t.Errorf("findButtonByName(%q) = %q, want %q", spoken, b.Name, want)
		}
	}

	// Short names must not fuzzy match each other or gibberish
	for _, spoken := range []string{"on", "disco", ""} {
		if b := findButtonByName(buttons, spoken); b != nil {
			t.Errorf("findButtonByName(%q) = %q, want no match", spoken, b.Name)
		}
	}
}

const testSkillId = "amzn1.ask.skill.11111111-2222-3333-4444-555555555555"

// Build the request envelope Alexa would post for a skill
func testEnvelopeJson(skillId string, requestType string, intentName string, slot string, stamp time.Time) []byte {

	intent := ""
	if intentName != "" {
		intent = fmt.Sprintf(`,"intent":{"name":%q,"slots":{"buttonName":{"name":"buttonName","value":%q}}}`,
			intentName, slot)
	}
	return []byte(fmt.Sprintf(
		`{"version":"1.0",`+
			`"session":{"application":{"applicationId":%q}},`+
			`"context":{"System":{"application":{"applicationId":%q}}},`+
			`"request":{"type":%q,"timestamp":%q%s}}`,
		skillId, skillId, requestType, stamp.Format(time.RFC3339), intent))
}

// Post a body to the skill handler, signed as Amazon would sign it unless a
// signature is given
func postSkillRequest(body []byte, signature string) *httptest.ResponseRecorder {

	if signature == "" {
		signature = testSignature(body)
	}
	handler := getSkillHandler("", testSkillId, testVerifier(validTestChain(), sharedTestPki().roots))

	r := httptest.NewRequest("POST", "/alexa/skill", bytes.NewReader(body))
	r.Header.Set("SignatureCertChainUrl", testChainUrl)
	r.Header.Set("Signature-256", signature)

	w := httptest.NewRecorder()
	handler(w, r)
	return w
}

func TestSkillHandlerAcceptsSignedRequest(t *testing.T) {

	body := testEnvelopeJson(testSkillId, "LaunchRequest", "", "", testNow)
	w := postSkillRequest(body, "")

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var response alexaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if response.Response.OutputSpeech == nil || response.Response.OutputSpeech.Text == "" {
		t.Error("response carried no speech, want a prompt")
	}
}

// The signature only proves a request came from Amazon — anyone can point
// their own skill at this endpoint, so a foreign application id must be
// refused even though it is perfectly signed
func TestSkillHandlerRejectsForeignSkill(t *testing.T) {

	body := testEnvelopeJson("amzn1.ask.skill.somebody-elses-skill",
		"IntentRequest", "RunButtonIntent", "rainbow", testNow)

	if w := postSkillRequest(body, ""); w.Code != 401 {
		t.Errorf("status = %d for a validly signed request from another skill, want 401", w.Code)
	}
}

func TestSkillHandlerRejectsBadRequests(t *testing.T) {

	cases := []struct {
		name      string
		body      []byte
		signature string
	}{
		{
			name:      "unsigned",
			body:      testEnvelopeJson(testSkillId, "LaunchRequest", "", "", testNow),
			signature: "bm90IGEgc2lnbmF0dXJl",
		},
		{
			name:      "body altered after signing",
			body:      testEnvelopeJson(testSkillId, "LaunchRequest", "", "", testNow),
			signature: testSignature(testEnvelopeJson(testSkillId, "LaunchRequest", "", "", testNow.Add(time.Second))),
		},
		{
			name: "replayed",
			body: testEnvelopeJson(testSkillId, "IntentRequest", "RunButtonIntent", "rainbow",
				testNow.Add(-10*time.Minute)),
		},
		{
			name: "no application id",
			body: testEnvelopeJson("", "LaunchRequest", "", "", testNow),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if w := postSkillRequest(c.body, c.signature); w.Code != 401 {
				t.Errorf("status = %d, want 401", w.Code)
			}
		})
	}
}

func TestSkillHandlerRejectsNonPost(t *testing.T) {

	handler := getSkillHandler("", testSkillId, testVerifier(validTestChain(), sharedTestPki().roots))
	w := httptest.NewRecorder()
	handler(w, httptest.NewRequest("GET", "/alexa/skill", nil))

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d for GET, want 405", w.Code)
	}
}

func TestHandleAlexaRequest(t *testing.T) {

	cases := []struct {
		name        string
		requestType string
		intent      string
		slot        string
		wantSpeech  bool
		wantEnd     bool
	}{
		{name: "launch asks what to run", requestType: "LaunchRequest", wantSpeech: true},
		{name: "help explains", requestType: "IntentRequest", intent: "AMAZON.HelpIntent", wantSpeech: true},
		{name: "stop says goodbye", requestType: "IntentRequest", intent: "AMAZON.StopIntent",
			wantSpeech: true, wantEnd: true},
		{name: "empty slot re-prompts", requestType: "IntentRequest", intent: "RunButtonIntent", wantSpeech: true},
		{name: "session end is silent", requestType: "SessionEndedRequest", wantEnd: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			var envelope alexaEnvelope
			body := testEnvelopeJson(testSkillId, c.requestType, c.intent, c.slot, testNow)
			if err := json.Unmarshal(body, &envelope); err != nil {
				t.Fatalf("test envelope is not valid JSON: %v", err)
			}

			response := handleAlexaRequest("", envelope)
			if got := response.Response.OutputSpeech != nil; got != c.wantSpeech {
				t.Errorf("has speech = %v, want %v", got, c.wantSpeech)
			}
			if response.Response.ShouldEndSession != c.wantEnd {
				t.Errorf("shouldEndSession = %v, want %v", response.Response.ShouldEndSession, c.wantEnd)
			}
			if response.Version != "1.0" {
				t.Errorf("version = %q, want 1.0", response.Version)
			}
		})
	}
}

// The application id is read from the session for a spoken interaction and
// from the context otherwise
func TestEnvelopeApplicationId(t *testing.T) {

	var contextOnly alexaEnvelope
	contextOnly.Context.System.Application.ApplicationId = testSkillId
	if contextOnly.applicationId() != testSkillId {
		t.Error("application id not read from context when session is empty")
	}

	var both alexaEnvelope
	both.Session.Application.ApplicationId = testSkillId
	both.Context.System.Application.ApplicationId = "other"
	if both.applicationId() != testSkillId {
		t.Error("session application id should win when both are present")
	}
}
