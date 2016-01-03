/* global require */
// Paths are relative to config.html
require([
        "./require/domReady!",
        "../js/ractive/ractive",
        "../js/ractive/ractive-touch",
        "../js/lib/controls",
        "../js/lib/scroll"],
    function (doc, R, touch, controls, scroll) {
        'use strict';

        // Ractive binding object
        var dto = {
            stripIndex: 0,
            stripLength: 100
        };

        function init() {

            var updateConfig = function () {
                var req = new XMLHttpRequest();
                req.open('PUT', '/StripLength/' + dto.stripIndex + ',' + dto.stripLength, true);
                req.send();
            };

            var ractive = new R({
                // The `el` option can be a node, an ID, or a CSS selector.
                el: 'container',

                // We could pass in a string, but for the sake of convenience
                // we're passing the ID of the <script> tag above.
                template: '#template',

                // Here, we're passing in some initial data
                data: dto
            });

            // Track user changes and update lights
            ractive.observe('stripIndex', updateConfig);
            ractive.observe('stripLength', updateConfig);

            // Wire up touch sliders to adjust index and length
            controls.touchSlider(document.getElementById("si"),
                function (v) {
                    dto.stripIndex = v;
                    ractive.set(dto);
                },
                dto.stripIndex, 0, 15);
            controls.touchSlider(document.getElementById("sl"),
                function (v) {
                    dto.stripLength = v;
                    ractive.set(dto);
                },
                dto.stripLength, 0, 1000);

            // Wire up the OK button
            controls.setOnTapNavigationMapping(ractive, 'btnGoBack', 'OK', './index.html');

            // Disable scrolling on everything
            scroll.disable(document.body);
        }

        init();
    });
