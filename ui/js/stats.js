/* global require */
// Paths are relative to stats.html
require([
        "./require/domReady!",
        "../js/ractive/ractive",
        "../js/ractive/ractive-touch",
        "../js/lib/controls",
        "../js/lib/scroll"],
    function (doc, R, touch, controls, scroll) {
        'use strict';

        // The ractive binding object
        var dto = {
            stats: undefined,
            // Must start true so we can wire up the buttons (i.e. button must exist in DOM)
            socketConnected: true,
            socketStatus: undefined
        };

        function init() {

            var ractive = new R({
                // The `el` option can be a node, an ID, or a CSS selector.
                el: 'container',

                // We could pass in a string, but for the sake of convenience
                // we're passing the ID of the <script> tag above.
                template: '#template',

                // Here, we're passing in some initial data
                data: dto
            });

            // Set OK button mapping
            controls.setOnTapNavigationMapping(ractive, 'btnGoBack', 'OK', './index.html');

            // e.g. "ws://localhost:8080/Stats"
            var wsUri = function () {
                // http:, , 192.168.0.X:8080, stats.html
                var parts = document.URL.split('/', 4);
                return "ws://" + parts[2] + "/Stats";
            }();

            // Open socket and set up handlers
            var ws = new WebSocket(wsUri, "P1");

            ws.onopen = function () {
                dto.socketStatus = "";
                dto.socketConnected = true;
                ractive.set(dto);
            };

            ws.onclose = function () {
                dto.socketStatus = "SOCKET DISCONNECTED";
                dto.socketConnected = false;
                ractive.set(dto);
            };

            ws.onerror = function (evt) {
                dto.socketStatus = 'SOCKET ERROR:' + evt.data;
                dto.socketConnected = false;
                ractive.set(dto);
            };

            ws.onmessage = function (evt) {
                // Parse socket payload into a stats object
                dto.stats = JSON.parse(evt.data);
                // We could use magic : true but that would cause an update for every property set
                // so instead we do ractive.set(dto) to update once after everything has been set
                ractive.set(dto);
            };

            // Disable scrolling on everything
            scroll.disable(document.body);
        }

        init();
    });
