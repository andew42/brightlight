package servers

import (
	"sync"
	"testing"
	"time"
)

// Broadcasting a button state update must never block, even when a listener
// is slow, has already exited, or is removed concurrently. The old
// implementation sent on unbuffered channels while holding listenersMux,
// which deadlocked against removeButtonListener.
func TestButtonStateBroadcastDoesNotBlock(t *testing.T) {

	done := make(chan struct{})
	go func() {
		defer close(done)

		var wg sync.WaitGroup
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				// Register a listener that never receives
				c := make(chan buttonState, 1)
				listenersMux.Lock()
				listeners = append(listeners, c)
				listenersMux.Unlock()

				updateActiveButtonKey(i)
				updateButtonPadVersion(i)

				// Remove without draining, racing with other broadcasts
				removeButtonListener(c)
			}(i)
		}
		wg.Wait()
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("button state broadcast deadlocked")
	}
}

// A listener must always end up observing the most recent state.
func TestButtonStateListenerReceivesLatestState(t *testing.T) {

	c := make(chan buttonState, 1)
	listenersMux.Lock()
	listeners = append(listeners, c)
	listenersMux.Unlock()
	defer removeButtonListener(c)

	// Two updates without the listener receiving in between: the first
	// queued state is dropped in favour of the latest
	updateActiveButtonKey(41)
	updateActiveButtonKey(42)

	select {
	case bs := <-c:
		if bs.ActiveButtonKey != 42 {
			t.Fatalf("expected latest key 42, got %d", bs.ActiveButtonKey)
		}
	case <-time.After(time.Second):
		t.Fatal("listener never received the update")
	}
}
