// Returns a web socket. The caller should close the
// returned socket when done. Calls cb with JSON pay
// load each time a message is received.
export function OpenWebSocket(url, cb) {

    // If url is FrameBuffer returns e.g.
    // "ws://localhost:8080/FrameBuffer"
    function wsUri(url) {

        // http:, , 192.168.0.X:3000, virtual.html
        let parts = document.URL.split('/', 4);

        // If we are running on the development server (port 3000)
        // Force web sockets to port 8080 because the development
        // server proxy doesn't work for web sockets
        let ip = parts[2].split(":");
        if (ip.length === 2 && ip[1] === "3000")
            parts[2] = ip[0] + ":8080";

        return "ws://" + parts[2] + "/" + url;
    }

    // Open socket and wire up handlers
    let fullUrl = wsUri(url);
    console.debug("Open " + url + " at " + fullUrl);
    let ws = new WebSocket(wsUri(url), "P1");

    ws.onopen = function () {
        console.debug(url + " web socket open")
    };

    ws.onclose = function () {
        console.debug(url + " web socket closed")
    };

    ws.onerror = function (evt) {
        console.error(url + " web socket error: " + evt.data);
    };

    let processingMessage = false;

    ws.onmessage = function (evt) {

        // Can this ever happen? (single threaded)
        if (processingMessage) {
            console.debug(url + " dropped a message");
            return;
        }

        console.debug(url + " web socket message");

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
