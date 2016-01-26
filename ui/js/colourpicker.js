/* global require */
// Paths are relative to colourpicker.html
require([
        "./require/domReady!",
        "../js/ractive/ractive",
        "../js/ractive/ractive-touch",
        "../js/tinycolor/tinycolor",
        "../js/lib/scroll",
        "../js/lib/controls",
        "../js/lib/nav",
        "../js/lib/lights"],
    function (doc, R, touch, tinycolor, scroll, controls, nav, lights) {
        'use strict';
        // Ractive data binding object
        var dto = {
            hue: 0,
            saturation: 100,
            value: 100,
            colour: undefined,
            hueColour: undefined,
            asHsv: undefined,
            asRgb: undefined
        };

        // Segment to update to show selected colour
        var segmentToUpdate;

        function init() {
            // Disable scrolling on everything
            scroll.disable(document.body);

            var ractive = new R({
                // The `el` option can be a node, an ID, or a CSS selector.
                el: 'container',

                // We could pass in a string, but for the sake of convenience
                // we're passing the ID of the <script> tag above.
                template: '#template',

                // Here, we're passing in some initial data
                data: dto
            });

            controls.verticalTouchSlider(document.getElementsByClassName("cp-slider-track-hue")[0], function (pos) {
                dto.hue = pos;
                updateUi();
            }, 0, 360);

            controls.verticalTouchSlider(document.getElementsByClassName("cp-slider-track-saturation")[0], function (pos) {
                dto.saturation = pos;
                updateUi();
            }, 0, 100);

            controls.verticalTouchSlider(document.getElementsByClassName("cp-slider-track-value")[0], function (pos) {
                dto.value = pos;
                updateUi();
            }, 0, 100);

            var updateUi = function () {

                // Compute derived fields
                var colour = tinycolor("hsv " + dto.hue + ", " + dto.saturation + "%, " + dto.value + "%");
                dto.colour = colour.toHex();
                dto.hueColour = tinycolor("hsv " + dto.hue + ", 100%, 100%").toHexString();
                dto.asHsv = colour.toHsvString();
                dto.asRgb = colour.toRgbString();

                // Update binding
                ractive.set(dto);

                // Update the room lights
                segmentToUpdate.params = dto.colour;
                lights.runAnimations([segmentToUpdate]);
            };

            // We expect to be passed a string colour parameter e.g. "FF0000" and a segment to update
            var p = nav.getParam();
            if (p !== undefined) {
                var c = tinycolor(p.colour).toHsv();
                segmentToUpdate = {
                    "segment": p.segment,
                    "animation": "Static",
                    "params": p.colour};
                dto.hue = c.h;
                dto.saturation = c.s * 100;
                dto.value = c.v * 100;
                updateUi();
            }

            // OK Cancel buttons
            ractive.on('okButtonHandler', function () {
                p.colour = dto.colour;
                nav.ret(p);
            });

            ractive.on('cancelButtonHandler', function () {

                // Restore original room light colour before returning
                segmentToUpdate.params = p.colour;
                lights.runAnimations([segmentToUpdate]);

                // Indicate we canceled
                p.colour = undefined;
                nav.ret(p);
            });

            updateUi();
        }

        // Initialise the page
        init();
    });
