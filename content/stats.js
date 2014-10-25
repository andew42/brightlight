window.onload = function () {main = init();}

var main;

function init() {

    var stats = document.getElementById("stats");

    // e.g. "ws://localhost:8080/Stats"
    var wsUri = function () {
        // http:, , 192.168.0.X:8080, stats.html
        var parts = document.URL.split('/', 4);
        return "ws://" + parts[2] + "/Stats";
    }();

    // Local high water marks
    var latchedMinFrameTime;
    var latchedMaxFrameTime;
    var latchedMinJitter;
    var latchedMaxJitter;
    var latchedMinSendTime;
    var latchedMaxSendTime;

    // Reset local high water marks
    var resetLatchedValues = function () {
        latchedMinFrameTime = 999.99;
        latchedMaxFrameTime = 0;
        latchedMinJitter = 999.99;
        latchedMaxJitter = 0;
        latchedMinSendTime = 999.99;
        latchedMaxSendTime = 0;

        var resetables = document.getElementsByName('resetMe');
        for (var i = 0; i < resetables.length; ++i) {
            resetables[i].innerHTML = '';
        }
    }

    var ws = new WebSocket(wsUri, "P1");

    ws.onopen = function (evt) {
        stats.innerHTML = "CONNECTED";
    };

    ws.onclose = function (evt) {
        stats.innerHTML = "DISCONNECTED";
    };

    ws.onmessage = function (evt) {
        var s = JSON.parse(evt.data);

        // Convert stats to ms values
        var averageFrameTime = (s.AverageFrameTime / 1000000);
        var minFrameTime = (s.MinFrameTime / 1000000);
        var maxFrameTime = (s.MaxFrameTime / 1000000);
        var averageJitter = (s.AverageJitter / 1000000);
        var minJitter = (s.MinJitter / 1000000);
        var maxJitter = (s.MaxJitter / 1000000);
        var averageSendTime = (s.AverageSendTime / 1000000);
        var minSendTime = (s.MinSendTime / 1000000);
        var maxSendTime = (s.MaxSendTime / 1000000);

        // Update high water marks
        if (minFrameTime < latchedMinFrameTime)
            latchedMinFrameTime = minFrameTime;
        if (maxFrameTime > latchedMaxFrameTime)
            latchedMaxFrameTime = maxFrameTime;
        if (minJitter < latchedMinJitter)
            latchedMinJitter = minJitter;
        if (maxJitter > latchedMaxJitter)
            latchedMaxJitter = maxJitter;
        if (minSendTime < latchedMinSendTime)
            latchedMinSendTime = minSendTime;
        if (maxSendTime > latchedMaxSendTime)
            latchedMaxSendTime = maxSendTime;

        var html = '<table>'
        html += '<tr><td style="color: gray">Avg Frame Time</td><td style="text-align: right">' + averageFrameTime.toFixed(2) + 'ms</td></tr>'
        html += '<tr><td style="color: gray">Min Frame Time</td><td style="text-align: right">' + minFrameTime.toFixed(2) + 'ms</td><td name="resetMe" style="text-align: right">' + latchedMinFrameTime.toFixed(2) + 'ms</td></tr>'
        html += '<tr><td style="color: gray">Max Frame Time</td><td style="text-align: right">' + maxFrameTime.toFixed(2) + 'ms</td><td name="resetMe" style="text-align: right">' + latchedMaxFrameTime.toFixed(2) + 'ms</td></tr>'
        html += '<tr><td style="color: gray">Avg Jitter</td><td style="text-align: right">' + averageJitter.toFixed(2) + 'ms</td></tr>'
        html += '<tr><td style="color: gray">Min Jitter</td><td style="text-align: right">' + minJitter.toFixed(2) + 'ms</td><td name="resetMe" style="text-align: right">' + latchedMinJitter.toFixed(2) + 'ms</td></tr>'
        html += '<tr><td style="color: gray">Max Jitter</td><td style="text-align: right">' + maxJitter.toFixed(2) + 'ms</td><td name="resetMe" style="text-align: right">' + latchedMaxJitter.toFixed(2) + 'ms</td></tr>'
        if (s.SendCount > 0) {
            html += '<tr><td style="color: gray">Avg Send Time</td><td style="text-align: right">' + averageSendTime.toFixed(2) + 'ms</td></tr>'
            html += '<tr><td style="color: gray">Min Send Time</td><td style="text-align: right">' + minSendTime.toFixed(2) + 'ms</td><td name="resetMe" style="text-align: right">' + latchedMinSendTime.toFixed(2) + 'ms</td></tr>'
            html += '<tr><td style="color: gray">Max Send Time</td><td style="text-align: right">' + maxSendTime.toFixed(2) + 'ms</td><td name="resetMe" style="text-align: right">' + latchedMaxSendTime.toFixed(2) + 'ms</td></tr>'
        }
        html += '<tr><td/><td/><td style="text-align: center"><button onclick="main.resetLatchedValues()">RESET</button></td></tr>'
        html += '</table>'
        stats.innerHTML = html;
    };

    ws.onerror = function (evt) {
        stats.innerHTML = '<span style="color: red;">ERROR:</span> ' + evt.data;
    };

    resetLatchedValues();

    // Return an object containing public properties
    return {resetLatchedValues:resetLatchedValues};
};
