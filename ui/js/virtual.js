// The ractive binding object
var dto = {
    socketConnected : true,
    socketStatus : undefined,
    leds : [[]]
};

function init() {
    var ractive = new Ractive({
        // The `el` option can be a node, an ID, or a CSS selector.
        el: 'container',

        // We could pass in a string, but for the sake of convenience
        // we're passing the ID of the <script> tag above.
        template: '#template',

        // Here, we're passing in some initial data
        data: dto
    });

    // Wire up back button
    buttons.setOnTapNavigationMapping(ractive, 'btnGoBack', 'OK', './index.html');

    // e.g. "ws://localhost:8080/FrameBuffer"
    var wsUri = function () {
        // http:, , 192.168.0.X:8080, virtual.html
        var parts = document.URL.split('/', 4);
        return "ws://" + parts[2] + "/FrameBuffer";
    }();

    // Open socket and wire up handlers
    var ws = new WebSocket(wsUri, "P1");

    ws.onopen = function () {
        dto.socketConnected = true;
        dto.socketStatus = "";
        ractive.set(dto);
    };

    ws.onclose = function () {
        dto.socketConnected = false;
        dto.socketStatus = "SOCKET DISCONNECTED";
        ractive.set(dto);
    };

    ws.onerror = function (evt) {
        dto.socketConnected = false;
        dto.socketStatus = 'SOCKET ERROR:' + evt.data;
        ractive.set(dto);
    };

    ws.onmessage = function (evt) {

        /** @namespace fb.Strips */
        var fb = JSON.parse(evt.data);
        var strips = fb.Strips;

        // Rebuild led array with first few leds from each strip
        // Ractive approach uses about half the cpu compared to
        // building HTML here (in chrome)
        dto.leds = [];
        for (var s = 0; s < strips.length; s++) {
            var fullStrip = strips[s];
            /** @namespace fullStrip.Leds */
            var leds = fullStrip.Leds;
            var displayStrip = [];
            var limitStripCount = 20;
            for (var l = 0; l < leds.length; l++) {
                if (limitStripCount-- <= 0)
                    break;
                displayStrip.push(lights.numberToColourCode(leds[l]));
            }
            dto.leds.push(displayStrip);
        }
        ractive.set(dto);
    };
}

window.onload = init;
