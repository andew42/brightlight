// Lighting API
/* global define, console */
define (function() {
    'use strict';
    var lights = {
        numberToColourCode: function (n) {
            var c = n.toString(16);
            while (c.length < 6) {
                c = '0' + c;
            }
            return c;
        },

        // Turn all lights on with the specified colour
        allLights: function (colour)
        {
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

        // Run the named animation
        animation : function (name) {
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

        // private
        setConnectionStatus : function(responseText) {
            var status = JSON.parse(responseText);
            if (lights.lastCallStatus !== status) {
                lights.lastCallStatus = status;
                if (lights.cbConnectionStatusChanged !== undefined) {
                    lights.cbConnectionStatusChanged();
                }
            }
        },

        // The connection status (of teensy) determined by last api call
        lastCallStatus : undefined,

        // Raised when connection status changes as the result of an api call
        cbConnectionStatusChanged : undefined
    };
    return lights;
});
