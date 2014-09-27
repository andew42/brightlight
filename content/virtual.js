window.onload = init;

function init() {
    var output = document.getElementById("output");
    var frameBuffer = document.getElementById("frameBuffer");

    // e.g. "ws://localhost:8080/FrameBuffer"
    var wsUri = function () {
        // http:, , 192.168.0.X:8080, virtual.html
        var parts = document.URL.split('/', 4);
        return "ws://" + parts[2] + "/FrameBuffer";
    }();

    var writeToScreen = function (message) {
        var pre = document.createElement("p");
        pre.style.wordWrap = "break-word";
        pre.innerHTML = message;
        output.appendChild(pre);
    }

    // REMOVE
    var writeHtmlToScreen = function (html) {
        var pre = document.createElement("div");
        pre.innerHTML = html;
        output.appendChild(pre);
    }

    var numberToColourCode = function (n) {
        var c = n.toString(16);
        while (c.length < 6)
            c = '0' + c;
        return c;
    }

    // Test web socket
    var ws = new WebSocket(wsUri, "P1");
    ws.onopen = function (evt) {
        var message = "WebSocket rocks";
        writeToScreen("CONNECTED");
        writeToScreen("SENT: " + message);
        ws.send(message);
    };
    ws.onclose = function (evt) {
        writeToScreen("DISCONNECTED");
    };
    ws.onmessage = function (evt) {
        // writeToScreen('<span style="color: blue;">RESPONSE: ' + evt.data + '</span>');

        var html = '<table>'
        var strips = JSON.parse(evt.data).Strips;
        for (s in strips) {
            html += '<tr>';
            var leds = strips[s].Leds;
            for (l in leds) {
                html += ('<td><div style="background-color:#' + numberToColourCode(leds[l]) + ';width:10px;height:10px"></div></td>');
            }
            html += '</tr>';
        }
        html += '</table>'

        // writeHtmlToScreen(evt.data);
        frameBuffer.innerHTML = html;

        // ws.close();
    };
    ws.onerror = function (evt) {
        writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data);
    };
}
