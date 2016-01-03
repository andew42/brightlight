// Controls helper API
/* global define */
define(['./util'], function (util) {
    'use strict';
    var controls = {
        // Setup a ractive event navigation mapping using on-tap extension
        // Not very DRY but on-tap is stripped from the rendered template
        // <div id="btnGoBack" on-tap="btnGoBack">
        setOnTapNavigationMapping: function (r, id, text, href, colour) {
            controls.setOnTapMapping(r, id, text, function () {
                window.location.href = href;
            }, colour);
        },

        // Setup a ractive event mapping using on-tap extension (see above)
        setOnTapMapping: function (r, id, text, fn, colour) {
            // Register on-tap event with same target name as id
            r.on(id, fn);
            // Set the button text...
            var e = document.getElementById(id);
            e.innerHTML = text;
            // ...and colour
            if (colour !== undefined) {
                e.setAttribute('style', ' color: #' + colour);
            }
        },

        // Pad touch slider control (adjusts value relative to existing value)
        touchSlider: function (el, onValueChange, initial, min, max, sensitivity) {

            sensitivity = util.defaultFor(sensitivity, 0.1);
            var origin = initial;
            var current = initial;

            // Capture the origin value at the start of adjustment
            el.addEventListener('touchstart', function () {
                origin = current;
            });

            // Track movement relative to start origin
            el.addEventListener('touchmove', function (evt) {
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
        verticalTouchSlider: function (el, onValueChange, min, max) {
            var range = max - min;

            // ### Touch support for iOS ###

            // Track touch movement over track
            el.addEventListener('touchmove', function (evt) {
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
            document.addEventListener('mousedown', function () {
                mousedown = true;
            });

            // Track mouse up over document
            document.addEventListener('mouseup', function () {
                mousedown = false;
            }, true);

            // Track mouse move over slider track
            el.addEventListener('mousemove', function (evt) {

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
    return controls;
});
