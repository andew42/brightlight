// Lighting API
/* global define, console */
define(function () {
    'use strict';
    var lights = {

        // Convert a number to a 6 digit zero padded colour
        numberToColourCode: function (n) {
            var c = n.toString(16);
            while (c.length < 6) {
                c = '0' + c;
            }
            return c;
        },

        // Animate a segment list, throttle to speed of server
        runAnimations: function (segments) {

            console.log(segments);

            // Return immediatly if nothing to do
            if (segments === undefined || segments.length === 0) {
                return;
            }

            // Throttle calls to server when user is swiping around
            if (busy) {
                // Make a copy of the segment list to run when last call completes
                nextRunAnimations = [];
                for (var s = 0; s < segments.length; s++) {
                    nextRunAnimations.push({
                        "segment": segments[s].segment,
                        "animation": segments[s].animation,
                        "params": segments[s].params});
                }
            } else {
                busy = true;
                var req = new XMLHttpRequest();
                req.open("POST", "/RunAnimations/");
                req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
                req.send(JSON.stringify(segments));
                req.onreadystatechange = function () {
                    if (req.readyState === 4) {
                        setConnectionStatus(req.responseText);
                        req = null;
                    }

                    // Run any pending animations so we always have the last one
                    busy = false;
                    if (nextRunAnimations !== undefined) {
                        lights.runAnimations(nextRunAnimations);
                        nextRunAnimations = undefined;
                    }
                };
            }
        },

        // The connection status (of teensy) determined by last api call
        lastCallStatus: undefined,

        // Raised when connection status changes as the result of an api call
        cbConnectionStatusChanged: undefined
    };

    var busy = false;
    var nextRunAnimations;

    var setConnectionStatus = function (responseText) {
        var status = JSON.parse(responseText);
        if (lights.lastCallStatus !== status) {
            lights.lastCallStatus = status;
            if (lights.cbConnectionStatusChanged !== undefined) {
                lights.cbConnectionStatusChanged();
            }
        }
    };

    return lights;
});
