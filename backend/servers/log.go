package servers

import (
	"encoding/json"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

// A web view of the server's own log, so diagnosing something (a voice command
// that found no button, a dropped USB port) does not mean logging into the Pi.
//
//	GET /api/Log        <- recent log lines as JSON
//	GET /api/Log.html   <- the same lines as an auto-refreshing page
//
// The lines are kept in memory rather than read back from journald: slog is
// already writing them, a ring buffer costs nothing, needs no journal
// permissions and behaves the same in development where there is no systemd.
// The cost is that this holds only the most recent lines and starts empty after
// a restart — journalctl is still the place to look for older history.
//
// Both routes are registered on the LAN server, never on the Alexa listener,
// for the same reason /api/AlexaButtons is: only /alexa/skill is reachable from
// the internet, and a log carries addresses and the configured skill id.

// Enough to cover a startup plus a good run of activity, and still trivial
// memory on a Pi 2B
const logBufferLines = 1000

type logRing struct {
	mu    sync.Mutex
	lines []string
}

var serverLog = &logRing{}

// Write accepts one formatted record per call — slog's TextHandler builds a
// whole line and writes it once — but splits anyway so a multi-line message
// cannot turn into one unreadable entry.
func (r *logRing) Write(p []byte) (int, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, line := range strings.Split(strings.TrimRight(string(p), "\n"), "\n") {
		r.lines = append(r.lines, line)
	}
	if excess := len(r.lines) - logBufferLines; excess > 0 {
		r.lines = append(r.lines[:0], r.lines[excess:]...)
	}
	return len(p), nil
}

func (r *logRing) snapshot() []string {

	r.mu.Lock()
	defer r.mu.Unlock()
	return append(make([]string, 0, len(r.lines)), r.lines...)
}

// LogWriter returns the writer slog output must be teed into for the log views
// to have anything to show. Wire it up in main alongside stdout.
func LogWriter() io.Writer {

	return serverLog
}

// LogHandler Handle GET /api/Log returning recent log lines, oldest first
func LogHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(serverLog.snapshot())
	if _, err := w.Write(response); err != nil {
		slog.Warn("Log failed to write response", "err", err.Error())
	}
}

// Kept plain and self-contained: no build step touches this, and it is read on
// a phone on the LAN while something is misbehaving.
var logPage = template.Must(template.New("log").Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Brightlight log</title>
<style>
 :root { color-scheme: light dark }
 body { font: 16px system-ui, sans-serif; margin: 0; padding: 14px }
 h1 { font-size: 1.15rem; margin: 0 0 .5rem }
 .bar { display: flex; flex-wrap: wrap; gap: .5rem; align-items: center; margin-bottom: .6rem }
 input[type=search] { font: inherit; padding: .35rem .5rem; flex: 1 1 12rem; min-width: 0 }
 label, .muted { font-size: .9rem; opacity: .75 }
 button { font: inherit; padding: .35rem .8rem; cursor: pointer }
 pre { margin: 0; padding: .6rem; border: 1px solid rgba(128,128,128,.4); border-radius: 4px;
       font: 12.5px ui-monospace, monospace; white-space: pre-wrap; word-break: break-word;
       max-height: 74vh; overflow: auto }
 .warn { color: #b40 }
 .err { color: #c00; font-weight: 600 }
</style></head><body>
<h1>Brightlight log</h1>
<div class="bar">
 <input id="filter" type="search" placeholder="filter, e.g. alexa" autocomplete="off">
 <label><input type="checkbox" id="auto" checked> auto-refresh</label>
 <button id="now">Refresh</button>
 <span id="count" class="muted"></span>
</div>
<pre id="log">loading…</pre>
<p class="muted">Last {{.Lines}} lines, held in memory — empty after a restart.
Older history is on the Pi: <code>journalctl -u brightlight</code>.</p>

<script>
 const pre = document.getElementById('log'), filter = document.getElementById('filter'),
       auto = document.getElementById('auto'), count = document.getElementById('count');
 let lines = [];

 function render() {
   const q = filter.value.toLowerCase();
   const shown = q ? lines.filter(l => l.toLowerCase().includes(q)) : lines;
   const atBottom = pre.scrollHeight - pre.scrollTop - pre.clientHeight < 40;
   pre.textContent = '';
   for (const l of shown) {
     const div = document.createElement('div');
     if (l.includes('level=ERROR')) div.className = 'err';
     else if (l.includes('level=WARN')) div.className = 'warn';
     div.textContent = l;
     pre.appendChild(div);
   }
   count.textContent = shown.length + (q ? ' of ' + lines.length : '') + ' lines';
   if (atBottom) pre.scrollTop = pre.scrollHeight;
 }

 async function load() {
   try {
     const r = await fetch('/api/Log', { cache: 'no-store' });
     lines = await r.json();
     render();
   } catch { count.textContent = 'server unreachable'; }
 }

 filter.addEventListener('input', render);
 document.getElementById('now').addEventListener('click', load);
 load();
 setInterval(() => { if (auto.checked) load(); }, 2000);
</script>
</body></html>`))

// LogPageHandler Handle GET /api/Log.html, the log as a filterable page
func LogPageHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := logPage.Execute(w, struct{ Lines int }{logBufferLines}); err != nil {
		slog.Warn("Log page failed to write response", "err", err.Error())
	}
}
