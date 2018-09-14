// Returns a web socket. The caller should close the
// returned socket when done. Calls cb with JSON pay
// load each time a message is received.
export function OpenWebSocket(url, cb) {

    // If url is FrameBuffer returns e.g.
    // "ws://localhost:8080/FrameBuffer"
    function wsUri(url) {

        // http:, , 192.168.0.X:8080, virtual.html
        let parts = document.URL.split('/', 4);
        return "ws://" + parts[2] + "/" + url;
    }

    // Open socket and wire up handlers
    let ws = new WebSocket(wsUri(url), "P1");

    ws.onopen = function () {
        console.info(url + " web socket open")
    };

    ws.onclose = function () {
        console.info(url + " web socket closed")
    };

    ws.onerror = function (evt) {
        console.info(url + " web socket error: " + evt.data);
    };

    let processingMessage = false;

    ws.onmessage = function (evt) {

        // Can this ever happen? (single threaded)
        if (processingMessage) {
            console.info(url + " dropped a message");
            return;
        }

        processingMessage = true;
        try {
            cb(JSON.parse(evt.data));
        }
        finally {
            processingMessage = false;
        }
    };
    return ws;
}
