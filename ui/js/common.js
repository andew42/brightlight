// jshint -W097
"use strict";

var util = {
    // Set default function parameter if parameter undefined
    defaultFor : function (arg, val) {
        return typeof arg !== 'undefined' ? arg : val;
    },

    // Apply a function to a single element or array
    applyTo : function (arg, f) {
        if (arg.constructor === Array || arg.constructor === HTMLCollection) {
            for (var i = 0; i < arg.length; i++) {
                f(arg[i]);
            }
        }
        else {
            f(arg);
        }
    }
};

// Lighting API
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

// iPhone element scroll enable / disable
var scroll = {
    // disable scrolling on this element
    disable : function(e) {
        // Allow single elements or arrays with applyTo
        util.applyTo(e, function(elmt) {
            elmt.addEventListener('touchmove', function (evt) {
                // Stop everything scrolling unless it is overridden
                if (!evt._isScroller) {
                    evt.preventDefault();
                }
            });
        });
    },

    // enable scroll on this element
    enable : function(e) {
        // Allow single elements or arrays with applyTo
        util.applyTo(e, function(elmt) {
            // range controls are a special case
            if (elmt.nodeName === "INPUT" && elmt.type === "range") {
                // touchmove handler - unconditional enable
                elmt.addEventListener('touchmove', function(evt) {
                    evt._isScroller = true;
                });
            }
            else {
                // touchstart handler
                elmt.addEventListener('touchstart', function() {
                    var top = elmt.scrollTop,
                        totalScroll = elmt.scrollHeight,
                        currentScroll = top + elmt.offsetHeight;

                    // If we're at the top or the bottom of the containers scroll, push up or down
                    // one pixel. This prevents the scroll from "passing through" to the body.

                    if(top === 0) {
                        elmt.scrollTop = 1;
                    }
                    else if(currentScroll === totalScroll) {
                        elmt.scrollTop = top - 1;
                    }
                });

                // touchmove handler
                elmt.addEventListener('touchmove', function(evt) {
                    // If content is actually scrollable, i.e. the content is long enough that scrolling can occur
                    if(elmt.offsetHeight < elmt.scrollHeight) {
                        evt._isScroller = true;
                    }
                });
            }
        });
    }
};

// Ractive button handling
var buttons = {
    // Setup a ractive event navigation mapping using on-tap extension
    // Not very DRY but on-tap is stripped from the rendered template
    // <div id="btnGoBack" on-tap="btnGoBack">
    setOnTapNavigationMapping : function(r, id, text, href, colour) {
        buttons.setOnTapMapping(r, id, text, function () {
            window.location.href = href;
        }, colour);
    },

    // Setup a ractive event mapping using on-tap extension (see above)
    setOnTapMapping : function(r, id, text, fn, colour) {
        // Register on-tap event with same target name as id
        r.on( id, fn);
        // Set the button text...
        var e = document.getElementById(id);
        e.innerHTML = text;
        // ...and colour
        if (colour !== undefined) {
            e.setAttribute('style', ' color: #' + colour);
        }
    },

    // Pad touch slider control (adjusts value relative to existing value)
    touchSlider : function(el, onValueChange, initial, min, max, sensitivity) {

        sensitivity = util.defaultFor(sensitivity, 0.1);
        var origin = initial;
        var current = initial;

        // Capture the origin value at the start of adjustment
        el.addEventListener('touchstart', function() {
            origin = current;
        });

        // Track movement relative to start origin
        el.addEventListener('touchmove', function(evt) {
            /** @namespace evt.targetTouches */
            if (evt.targetTouches.length === 1) {
                var t = evt.targetTouches[0];
                var r = el.getBoundingClientRect();
                var dif = ((r.top + (r.height / 2)) - t.clientY) * sensitivity;
                current = Math.round(origin + dif);
                if (current < min) {
                    current = min;
                }
                else if (current > max) {
                    current = max;
                }
                onValueChange(current);
                evt.preventDefault();
            }
        });
    },

    // Touch slider control (adjusts value based on position on track)
    verticalTouchSlider : function(el, onValueChange, min, max) {
        var range = max - min;

        // ### Touch support for iOS ###

        // Track touch movement over track
        el.addEventListener('touchmove', function(evt) {
            // Calculate position in pixels
            var r = el.getBoundingClientRect();
            var pos = r.bottom - evt.targetTouches[0].clientY;

            // Convert to value between min and max
            pos = (pos / r.height) * range + min;

            // Ensue position falls withing range
            if (pos < min) {
                pos = min;
            }
            else if (pos > max) {
                pos = max;
            }
            onValueChange(pos);
            evt.preventDefault();
        });

        // ### Mouse support for desktop browsers ###

        var mousedown = false;

        // Track mouse downs over document
        document.addEventListener('mousedown', function() {
            mousedown = true;
        });

        // Track mouse up over document
        document.addEventListener('mouseup', function() {
            mousedown = false;
        }, true);

        // Track mouse move over slider track
        el.addEventListener('mousemove', function(evt) {

            if (!mousedown) {
                return;
            }

            // Calculate position in pixels
            var r = el.getBoundingClientRect();
            var pos = r.bottom - evt.clientY;

            // Convert to value between min and max
            pos = (pos / r.height) * range + min;

            // Ensue position falls withing range
            if (pos < min) {
                pos = min;
            }
            else if (pos > max) {
                pos = max;
            }
            onValueChange(pos);
        }, true);
    }
};
