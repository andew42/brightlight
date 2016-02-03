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

            // We expect the element to have some padding which will be
            // treated as a dead zone to allow easy selection of min max
            // values. Padding is expected to be in pixels
            var padding = parseInt(window.getComputedStyle(el,null).padding, 10);
            if (padding === undefined) {
                padding = 0;
            }

            // Control range (typically 360 or 100)
            var range = max - min;

            // Helper converts pixel positions into range position and calls onValueChange
            var updatePosition = function(mouseX, mouseY) {

                var r = el.getBoundingClientRect();

                // Ensure mouse / finger is within control
                if (mouseX < r.left || mouseX > r.right || mouseY > r.bottom || mouseY < r.top) {
                    return;
                }

                // Where is the mouse positioned within track
                var pos;
                if ((r.bottom - mouseY) < padding) {
                    // Position is in the lower dead zone
                    pos = 0;
                } else if (mouseY < (r.bottom - r.height + padding)) {
                    // Position is in the upper dead zone
                    pos = max;
                } else {
                    // Calculate position in pixels
                    pos = (r.bottom - padding) - mouseY;
                    // Convert to value between min and max
                    pos = (pos / (r.height - 2 * padding)) * range + min;
                }
                // Inform the caller
                onValueChange(pos);
            };

            // === Touch support for iOS ===

            // Respond to initial touch
            el.addEventListener('touchstart', function (evt) {
                updatePosition(evt.targetTouches[0].clientX, evt.targetTouches[0].clientY);
                evt.preventDefault();
            });

            // Track touch movement over track
            el.addEventListener('touchmove', function (evt) {
                updatePosition(evt.targetTouches[0].clientX, evt.targetTouches[0].clientY);
                evt.preventDefault();
            });

            // === Mouse support for desktop browsers ===

            var mousedown = false;

            // Track mouse downs over document
            document.addEventListener('mousedown', function (evt) {
                mousedown = true;
                updatePosition(evt.clientX, evt.clientY);
            });

            // Track mouse up over document
            document.addEventListener('mouseup', function () {
                mousedown = false;
            }, true);

            // Track mouse move over slider track
            el.addEventListener('mousemove', function (evt) {
                if (mousedown) {
                    updatePosition(evt.clientX, evt.clientY);
                }
            }, true);
        }
    };
    return controls;
});
