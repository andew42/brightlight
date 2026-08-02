package servers

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"
)

// Verification of requests from the Alexa service. Amazon signs the raw
// request body with a certificate it publishes at a URL given in the
// SignatureCertChainUrl header, so we validate that URL, fetch and chain
// verify the certificate, confirm it was issued for the Alexa service and
// finally check the body against the Signature-256 header.
//
// A valid signature only proves the request came from Amazon, not that it
// came from our skill — anyone can point their own skill at this endpoint.
// Checking the application id is what binds a request to our skill and that
// is done by the caller.

const (
	alexaCertChainHost = "s3.amazonaws.com"
	alexaCertChainPath = "/echo.api/"
	alexaCertHostname  = "echo-api.amazon.com"

	// Amazon requires requests older than this to be rejected as replays
	alexaTimestampWindow = 150 * time.Second

	// The signing certificate changes rarely so it is cached, but the URL
	// comes from the request, so cap the cache to keep it a fixed cost
	alexaCertCacheLimit = 4
)

// Verifies Alexa request signatures. fetch, now and roots are fields so
// tests can supply a fake certificate chain, a fixed clock and their own
// trust root; production uses HTTP, the wall clock and the system roots.
type alexaVerifier struct {
	fetch func(url string) ([]byte, error)
	now   func() time.Time
	roots *x509.CertPool

	mux    sync.Mutex
	cached map[string][]*x509.Certificate
}

func newAlexaVerifier() *alexaVerifier {

	client := &http.Client{Timeout: 5 * time.Second}

	return &alexaVerifier{
		now:    time.Now,
		cached: make(map[string][]*x509.Certificate),
		fetch: func(u string) ([]byte, error) {
			response, err := client.Get(u)
			if err != nil {
				return nil, err
			}
			defer response.Body.Close()
			if response.StatusCode != 200 {
				return nil, errors.New("unexpected status " + response.Status)
			}
			return io.ReadAll(io.LimitReader(response.Body, 64*1024))
		},
	}
}

// Check that body was signed by the Alexa service, where chainUrl and
// signature are the SignatureCertChainUrl and Signature-256 request headers
func (v *alexaVerifier) verify(chainUrl string, signature string, body []byte) error {

	if signature == "" {
		return errors.New("missing Signature-256 header")
	}
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return errors.New("Signature-256 is not valid base64")
	}

	leaf, err := v.signingCertificate(chainUrl)
	if err != nil {
		return err
	}

	key, ok := leaf.PublicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("signing certificate does not hold an RSA key")
	}

	hashed := sha256.Sum256(body)
	if rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], sig) != nil {
		return errors.New("request body does not match its signature")
	}
	return nil
}

// Fetch (or recall) the certificate the request was signed with, having
// checked it chains to a trusted root and belongs to the Alexa service
func (v *alexaVerifier) signingCertificate(chainUrl string) (*x509.Certificate, error) {

	if err := validateCertChainUrl(chainUrl); err != nil {
		return nil, err
	}

	chain, err := v.certificateChain(chainUrl)
	if err != nil {
		return nil, err
	}

	// Verify on every request rather than at fetch time so a cached chain
	// still stops being accepted the moment it expires
	intermediates := x509.NewCertPool()
	for _, c := range chain[1:] {
		intermediates.AddCert(c)
	}
	leaf := chain[0]
	if _, err = leaf.Verify(x509.VerifyOptions{
		Roots:         v.roots,
		Intermediates: intermediates,
		CurrentTime:   v.now(),
	}); err != nil {
		return nil, errors.New("signing certificate is not trusted: " + err.Error())
	}

	if leaf.VerifyHostname(alexaCertHostname) != nil {
		return nil, errors.New("signing certificate is not for " + alexaCertHostname)
	}
	return leaf, nil
}

func (v *alexaVerifier) certificateChain(chainUrl string) ([]*x509.Certificate, error) {

	v.mux.Lock()
	cached, ok := v.cached[chainUrl]
	v.mux.Unlock()
	if ok {
		return cached, nil
	}

	content, err := v.fetch(chainUrl)
	if err != nil {
		return nil, errors.New("could not fetch signing certificate: " + err.Error())
	}
	chain, err := parseCertificateChain(content)
	if err != nil {
		return nil, err
	}

	v.mux.Lock()
	if len(v.cached) >= alexaCertCacheLimit {
		v.cached = make(map[string][]*x509.Certificate)
	}
	v.cached[chainUrl] = chain
	v.mux.Unlock()
	return chain, nil
}

// Split a PEM bundle into the leaf certificate followed by any intermediates
func parseCertificateChain(content []byte) ([]*x509.Certificate, error) {

	var chain []*x509.Certificate
	for block, rest := pem.Decode(content); block != nil; block, rest = pem.Decode(rest) {
		if block.Type != "CERTIFICATE" {
			continue
		}
		c, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, errors.New("signing certificate could not be parsed")
		}
		chain = append(chain, c)
	}
	if len(chain) == 0 {
		return nil, errors.New("no certificate found at the signing certificate URL")
	}
	return chain, nil
}

// The certificate URL arrives in a request header, so it is only trusted
// once it is known to point at Amazon's own certificate location
func validateCertChainUrl(raw string) error {

	if raw == "" {
		return errors.New("missing SignatureCertChainUrl header")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return errors.New("SignatureCertChainUrl is not a valid URL")
	}
	if !strings.EqualFold(u.Scheme, "https") {
		return errors.New("SignatureCertChainUrl is not https")
	}
	if !strings.EqualFold(u.Hostname(), alexaCertChainHost) {
		return errors.New("SignatureCertChainUrl host is not " + alexaCertChainHost)
	}
	if port := u.Port(); port != "" && port != "443" {
		return errors.New("SignatureCertChainUrl port is not 443")
	}
	// Cleaning the path first stops /echo.api/../elsewhere passing the check
	if !strings.HasPrefix(path.Clean(u.Path), alexaCertChainPath) {
		return errors.New("SignatureCertChainUrl path is not under " + alexaCertChainPath)
	}
	return nil
}

// Reject stale (replayed) and implausibly future requests
func checkAlexaTimestamp(stamp time.Time, now time.Time) error {

	if stamp.IsZero() {
		return errors.New("request has no timestamp")
	}
	skew := now.Sub(stamp)
	if skew > alexaTimestampWindow || skew < -alexaTimestampWindow {
		return errors.New("request timestamp is outside the allowed window")
	}
	return nil
}
