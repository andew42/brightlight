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

        // Segments to animate to show selected colour in the room
        var roomSegments;
        // Reference to segment to edit within roomSegments
        var editSegment;

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

            controls.verticalTouchSlider(document.getElementsByClassName("cp-slider-hue")[0], function (pos) {
                dto.hue = pos;
                updateUi();
            }, 0, 360);

            controls.verticalTouchSlider(document.getElementsByClassName("cp-slider-saturation")[0], function (pos) {
                dto.saturation = pos;
                updateUi();
            }, 0, 100);

            controls.verticalTouchSlider(document.getElementsByClassName("cp-slider-value")[0], function (pos) {
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
                editSegment.params = dto.colour;
                lights.runAnimations(roomSegments);
            };

            // We expect to be passed roomSegments (all the animations active in the room)
            // plus editSegmentIndex which is a relative to roomSegments
            var initialColour;
            var p = nav.getParam();
            if (p !== undefined) {
                roomSegments = p.roomSegments;
                editSegment = p.roomSegments[p.editSegmentIndex];
                initialColour = editSegment.params;
                var c = tinycolor(initialColour).toHsv();
                dto.hue = c.h;
                dto.saturation = c.s * 100;
                dto.value = c.v * 100;
                updateUi();
            }

            // OK button hit, return new colour
            ractive.on('okButtonHandler', function () {
                p.newColour = editSegment.params;
                nav.ret(p);
            });

            // Cancel button hit
            ractive.on('cancelButtonHandler', function () {

                // Restore original room light colour before returning
                editSegment.params = initialColour;
                lights.runAnimations(roomSegments);
                nav.ret(p);
            });

            updateUi();
        }

        // Initialise the page
        init();
    });
