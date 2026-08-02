package servers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"sync"
	"testing"
	"time"
)

// Fixed instant so certificate validity and request timestamps are stable
var testNow = time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)

const testChainUrl = "https://s3.amazonaws.com/echo.api/echo-api-cert.pem"

// Stands in for Amazon: an authority we trust in tests and the key that
// signs requests. Generating RSA keys is slow so one set is shared.
type testPki struct {
	caKey   *rsa.PrivateKey
	caCert  *x509.Certificate
	leafKey *rsa.PrivateKey
	roots   *x509.CertPool
}

var sharedTestPki = sync.OnceValue(func() *testPki {

	caKey := mustGenerateKey()
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test root"},
		NotBefore:             testNow.Add(-365 * 24 * time.Hour),
		NotAfter:              testNow.Add(365 * 24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &caKey.PublicKey, caKey)
	if err != nil {
		panic(err)
	}
	caCert, err := x509.ParseCertificate(der)
	if err != nil {
		panic(err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(caCert)
	return &testPki{caKey: caKey, caCert: caCert, leafKey: mustGenerateKey(), roots: roots}
})

func mustGenerateKey() *rsa.PrivateKey {

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return key
}

// Issue a signing certificate and return the PEM chain that would be
// published at the certificate URL
func (p *testPki) chain(hostname string, notBefore time.Time, notAfter time.Time) []byte {

	template := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: hostname},
		DNSNames:     []string{hostname},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, p.caCert, &p.leafKey.PublicKey, p.caKey)
	if err != nil {
		panic(err)
	}

	encode := func(b []byte) []byte {
		return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: b})
	}
	return append(encode(der), encode(p.caCert.Raw)...)
}

// A chain that would be accepted by a correctly behaving verifier
func validTestChain() []byte {

	return sharedTestPki().chain(alexaCertHostname,
		testNow.Add(-time.Hour), testNow.Add(time.Hour))
}

func testVerifier(chain []byte, roots *x509.CertPool) *alexaVerifier {

	return &alexaVerifier{
		roots:  roots,
		now:    func() time.Time { return testNow },
		cached: make(map[string][]*x509.Certificate),
		fetch:  func(string) ([]byte, error) { return chain, nil },
	}
}

func testSignature(body []byte) string {

	hashed := sha256.Sum256(body)
	sig, err := rsa.SignPKCS1v15(rand.Reader, sharedTestPki().leafKey, crypto.SHA256, hashed[:])
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(sig)
}

func TestValidateCertChainUrl(t *testing.T) {

	valid := []string{
		"https://s3.amazonaws.com/echo.api/echo-api-cert-7.pem",
		"https://s3.amazonaws.com:443/echo.api/echo-api-cert-7.pem",
		"https://S3.AMAZONAWS.COM/echo.api/echo-api-cert-7.pem",
	}
	for _, u := range valid {
		if err := validateCertChainUrl(u); err != nil {
			t.Errorf("validateCertChainUrl(%q) = %v, want accepted", u, err)
		}
	}

	invalid := map[string]string{
		"": "empty",
		"http://s3.amazonaws.com/echo.api/cert.pem":           "not https",
		"https://notamazon.com/echo.api/cert.pem":             "wrong host",
		"https://s3.amazonaws.com.evil.com/echo.api/cert.pem": "host is a prefix only",
		"https://s3.amazonaws.com/evil/cert.pem":              "wrong path",
		"https://s3.amazonaws.com/echo.api/../evil/cert.pem":  "path traversal",
		"https://s3.amazonaws.com:8443/echo.api/cert.pem":     "wrong port",
	}
	for u, why := range invalid {
		if err := validateCertChainUrl(u); err == nil {
			t.Errorf("validateCertChainUrl(%q) accepted, want rejected (%s)", u, why)
		}
	}
}

func TestVerifyAcceptsAmazonSignature(t *testing.T) {

	body := []byte(`{"request":{"type":"LaunchRequest"}}`)
	v := testVerifier(validTestChain(), sharedTestPki().roots)

	if err := v.verify(testChainUrl, testSignature(body), body); err != nil {
		t.Errorf("verify of a correctly signed request = %v, want accepted", err)
	}
}

func TestVerifyRejects(t *testing.T) {

	body := []byte(`{"request":{"type":"LaunchRequest"}}`)
	pki := sharedTestPki()

	cases := []struct {
		name        string
		chain       []byte
		roots       *x509.CertPool
		chainUrl    string
		signature   string
		noSignature bool
		body        []byte
	}{
		{
			name:        "missing signature header",
			noSignature: true,
		},
		{
			name:      "signature is not base64",
			signature: "not base64 at all!",
		},
		{
			name:      "signature of a different body",
			signature: testSignature([]byte(`{"request":{"type":"SomethingElse"}}`)),
		},
		{
			name: "body altered after signing",
			body: []byte(`{"request":{"type":"LaunchRequest"} }`),
		},
		{
			name:     "certificate url not at Amazon",
			chainUrl: "https://evil.example.com/echo.api/cert.pem",
		},
		{
			name:  "certificate is not for the Alexa service",
			chain: pki.chain("evil.example.com", testNow.Add(-time.Hour), testNow.Add(time.Hour)),
		},
		{
			name:  "certificate has expired",
			chain: pki.chain(alexaCertHostname, testNow.Add(-2*time.Hour), testNow.Add(-time.Hour)),
		},
		{
			name:  "certificate is not yet valid",
			chain: pki.chain(alexaCertHostname, testNow.Add(time.Hour), testNow.Add(2*time.Hour)),
		},
		{
			name:  "certificate does not chain to a trusted root",
			roots: x509.NewCertPool(),
		},
		{
			name:  "no certificate at the url",
			chain: []byte("this is not a pem file"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			chain := c.chain
			if chain == nil {
				chain = validTestChain()
			}
			roots := c.roots
			if roots == nil {
				roots = pki.roots
			}
			chainUrl := c.chainUrl
			if chainUrl == "" {
				chainUrl = testChainUrl
			}
			signature := c.signature
			if signature == "" && !c.noSignature {
				signature = testSignature(body)
			}
			signed := c.body
			if signed == nil {
				signed = body
			}

			v := testVerifier(chain, roots)
			if err := v.verify(chainUrl, signature, signed); err == nil {
				t.Error("verify accepted the request, want rejected")
			}
		})
	}
}

// A cached chain must stop being accepted once it expires rather than being
// trusted for as long as the process lives
func TestVerifyRejectsCachedCertificateOnceExpired(t *testing.T) {

	body := []byte(`{"request":{"type":"LaunchRequest"}}`)
	v := testVerifier(validTestChain(), sharedTestPki().roots)

	if err := v.verify(testChainUrl, testSignature(body), body); err != nil {
		t.Fatalf("first verify = %v, want accepted", err)
	}

	// Same cached chain, clock moved past the certificate's expiry
	v.now = func() time.Time { return testNow.Add(2 * time.Hour) }
	if err := v.verify(testChainUrl, testSignature(body), body); err == nil {
		t.Error("verify accepted an expired cached certificate, want rejected")
	}
}

func TestCheckAlexaTimestamp(t *testing.T) {

	cases := []struct {
		name   string
		stamp  time.Time
		accept bool
	}{
		{"now", testNow, true},
		{"just inside the window", testNow.Add(-149 * time.Second), true},
		{"slightly ahead", testNow.Add(10 * time.Second), true},
		{"replayed", testNow.Add(-10 * time.Minute), false},
		{"far in the future", testNow.Add(10 * time.Minute), false},
		{"missing", time.Time{}, false},
	}

	for _, c := range cases {
		err := checkAlexaTimestamp(c.stamp, testNow)
		if c.accept && err != nil {
			t.Errorf("%s: checkAlexaTimestamp = %v, want accepted", c.name, err)
		}
		if !c.accept && err == nil {
			t.Errorf("%s: checkAlexaTimestamp accepted, want rejected", c.name)
		}
	}
}
