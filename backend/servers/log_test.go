package servers

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestLogRingKeepsNewestLines(t *testing.T) {

	ring := &logRing{}
	total := logBufferLines + 200
	for i := 0; i < total; i++ {
		if _, err := fmt.Fprintf(ring, "line %d\n", i); err != nil {
			t.Fatal(err)
		}
	}

	lines := ring.snapshot()
	if len(lines) != logBufferLines {
		t.Fatalf("buffer holds %d lines, want %d", len(lines), logBufferLines)
	}
	// The oldest 200 must have been dropped, not the newest
	if want := fmt.Sprintf("line %d", total-logBufferLines); lines[0] != want {
		t.Errorf("first line is %q, want %q", lines[0], want)
	}
	if want := fmt.Sprintf("line %d", total-1); lines[len(lines)-1] != want {
		t.Errorf("last line is %q, want %q", lines[len(lines)-1], want)
	}
}

func TestLogRingSplitsMultiLineWrites(t *testing.T) {

	ring := &logRing{}
	// A panic or stack trace arrives as one write of several lines; kept apart
	// so the page can colour and filter each, and with no blank entry from the
	// trailing newline slog always writes
	if _, err := ring.Write([]byte("first\nsecond\nthird\n")); err != nil {
		t.Fatal(err)
	}

	lines := ring.snapshot()
	if got := strings.Join(lines, "|"); got != "first|second|third" {
		t.Errorf("got %q, want \"first|second|third\"", got)
	}
}

func TestLogRingIsSafeForConcurrentWriters(t *testing.T) {

	// slog is called from every driver goroutine, and the page reads while they
	// write; run with -race to make this mean something
	ring := &logRing{}
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				fmt.Fprintf(ring, "writer %d line %d\n", w, i)
			}
		}(w)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			ring.snapshot()
		}
	}()
	wg.Wait()

	if len(ring.snapshot()) != logBufferLines {
		t.Errorf("expected a full buffer, got %d lines", len(ring.snapshot()))
	}
}
