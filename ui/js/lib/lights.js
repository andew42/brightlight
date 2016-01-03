// Lighting API
/* global define, console */
define(function () {
    'use strict';
    var lights = {
        numberToColourCode: function (n) {
            var c = n.toString(16);
            while (c.length < 6) {
                c = '0' + c;
            }
            return c;
        },

        // Animate a segment list
        runAnimations: function (segments) {
            console.log(segments);
            var req = new XMLHttpRequest();
            req.open("POST", "/RunAnimations/");
            req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
            req.send(JSON.stringify(segments));
            req.onreadystatechange = function () {
                if (req.readyState === 4) {
                    lights.setConnectionStatus(req.responseText);
                    req = null;
                }
            };
        },

        // Turn all lights on with the specified colour TODO remove
        allLights: function (colour) {
            console.log(colour);
            var req = new XMLHttpRequest();
            req.open('PUT', '/AllLights/' + colour, true);
            req.send();
            req.onreadystatechange = function () {
                if (req.readyState === 4) {
                    lights.setConnectionStatus(req.responseText);
                    req = null;
                }
            };
        },

        // Run the named animation TODO remove
        animation: function (name) {
            console.log(name);
            var req = new XMLHttpRequest();
            req.open('PUT', '/Animation/' + name, true);
            req.send();
            req.onreadystatechange = function () {
                if (req.readyState === 4) {
                    lights.setConnectionStatus(req.responseText);
                    req = null;
                }
            };
        },

        // TODO: private
        setConnectionStatus: function (responseText) {
            var status = JSON.parse(responseText);
            if (lights.lastCallStatus !== status) {
                lights.lastCallStatus = status;
                if (lights.cbConnectionStatusChanged !== undefined) {
                    lights.cbConnectionStatusChanged();
                }
            }
        },

        // The connection status (of teensy) determined by last api call
        lastCallStatus: undefined,

        // Raised when connection status changes as the result of an api call
        cbConnectionStatusChanged: undefined
    };
    return lights;
});
