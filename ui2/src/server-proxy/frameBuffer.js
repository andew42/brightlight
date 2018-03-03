// Returns a web socket to for a frame buffer connection
// The caller should close the returned socket when done
// Calls cb with each frame buffer frame every 40ms
export function Open(cb) {

    // Frame Buffer URL e.g. "ws://localhost:8080/FrameBuffer"
    function wsUri() {

        // http:, , 192.168.0.X:8080, virtual.html
        let parts = document.URL.split('/', 4);
        return "ws://" + parts[2] + "/FrameBuffer";
    }

    // Convert a number to a 6 digit zero padded colour
    function numberToColourCode(n) {
        let c = n.toString(16);
        while (c.length < 6) {
            c = '0' + c;
        }
        return c;
    }

    // Open socket and wire up handlers
    let ws = new WebSocket(wsUri(), "P1");

    ws.onopen = function () {
        console.info("Frame Buffer web socket open")
    };

    ws.onclose = function () {
        console.info("Frame Buffer web socket closed")
    };

    ws.onerror = function (evt) {
        console.info("Frame Buffer web socket error: " + evt.data);
    };

    let processingMessage = false;

    ws.onmessage = function (evt) {

        // TODO: Can this ever happen? (single threaded)
        if (processingMessage) {
            console.info("Dropped frame buffer frame");
            return;
        }

        processingMessage = true;
        try {
            let fb = JSON.parse(evt.data);
            let strips = fb.Strips;

            // Rebuild led array with first few (20) leds from each strip
            let rc = [];
            for (let s = 0; s < strips.length; s++) {
                let fullStrip = strips[s];
                let leds = fullStrip.Leds;
                let displayStrip = [];
                let limitStripCount = 20;
                for (let l = 0; l < leds.length; l++) {
                    if (limitStripCount-- <= 0) {
                        break;
                    }
                    displayStrip.push(numberToColourCode(leds[l]));
                }
                rc.push(displayStrip);
            }
            cb(rc);
        }
        finally {
            processingMessage = false;
        }
    };
    return ws;
}
