package segment

import (
	"testing"

	"github.com/andew42/brightlight/framebuffer"
)

// Physical pXX:YY segment names arrive from the HTTP API so out of range
// values must return an error rather than panicking the animation driver
func TestGetNamedSegmentPhysicalBounds(t *testing.T) {

	// A valid strip index and length
	if _, err := GetNamedSegment("p0:10"); err != nil {
		t.Errorf("GetNamedSegment(p0:10) unexpected error: %v", err)
	}

	// Out of range indexes and lengths must be rejected
	for _, name := range []string{"p99:10", "p-1:5", "p0:99999", "p0:-1"} {
		if _, err := GetNamedSegment(name); err == nil {
			t.Errorf("GetNamedSegment(%q) expected an error", name)
		}
	}
}

// A zero length request (including on an unused zero length strip) is valid
// and must resolve to an empty segment without panicking
func TestGetNamedSegmentZeroLength(t *testing.T) {

	fb := framebuffer.NewFrameBuffer()
	for i := range fb.Strips {
		name := "p" + string(rune('0'+i%10)) + ":0"
		if i > 9 {
			continue // single digit indexes are enough for this test
		}
		ns, err := GetNamedSegment(name)
		if err != nil {
			t.Errorf("GetNamedSegment(%q) unexpected error: %v", name, err)
			continue
		}
		if got := ns.GetSegment(fb).Len(); got != 0 {
			t.Errorf("GetNamedSegment(%q).Len() = %d, want 0", name, got)
		}
	}
}
