/* global require */
// Paths are relative to stats.html
require([
        "./require/domReady!",
        "../js/ractive/ractive.js",
        "../js/ractive/ractive-events-tap.js",
        "../js/lib/controls.js",
        "../js/lib/scroll"],
    function(doc,R,tap,controls,scroll) {
        'use strict';

        // The ractive binding object
        var dto = {
            // Must start true so we can wire up the buttons (i.e. reset button must exist in DOM)
            socketConnected: true,
            socketStatus: undefined,

            averageFrameTime: undefined,
            minFrameTime: undefined,
            latchedMinFrameTime: undefined,
            maxFrameTime: undefined,
            latchedMaxFrameTime: undefined,

            averageJitter: undefined,
            minJitter: undefined,
            latchedMinJitter: undefined,
            maxJitter: undefined,
            latchedMaxJitter: undefined,

            showSendTimes: false,
            averageSendTime: undefined,
            minSendTime: undefined,
            latchedMinSendTime: undefined,
            maxSendTime: undefined,
            latchedMaxSendTime: undefined
        };

        function init() {

            var ractive = new R({
                // Attach tap handler extension
                events: {tap: tap },

                // The `el` option can be a node, an ID, or a CSS selector.
                el: 'container',

                // We could pass in a string, but for the sake of convenience
                // we're passing the ID of the <script> tag above.
                template: '#template',

                // Here, we're passing in some initial data
                data: dto
            });

            // Set OK and Reset button mappings
            controls.setOnTapNavigationMapping(ractive, 'btnGoBack', 'OK', './index.html');
            controls.setOnTapMapping(ractive, 'btnReset', 'RESET', function () {
                dto.latchedMinFrameTime = undefined;
                dto.latchedMaxFrameTime = undefined;
                dto.latchedMinJitter = undefined;
                dto.latchedMaxJitter = undefined;
                dto.latchedMinSendTime = undefined;
                dto.latchedMaxSendTime = undefined;
                ractive.set(dto);
            });

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
                // Parse socket payload into a status (s) object
                /** @namespace s.AverageFrameTime */
                /** @namespace s.MinFrameTime */
                /** @namespace s.MaxFrameTime */
                /** @namespace s.AverageJitter */
                /** @namespace s.MaxJitter */
                /** @namespace s.MinJitter */
                /** @namespace s.AverageSendTime */
                /** @namespace s.MinSendTime */
                /** @namespace s.MaxSendTime */
                /** @namespace s.SendCount */
                var s = JSON.parse(evt.data);

                // Convert stats to 2dp ms values
                dto.averageFrameTime = (s.AverageFrameTime / 1000000).toFixed(2);
                dto.minFrameTime = (s.MinFrameTime / 1000000).toFixed(2);
                dto.maxFrameTime = (s.MaxFrameTime / 1000000).toFixed(2);
                dto.averageJitter = (s.AverageJitter / 1000000).toFixed(2);
                dto.minJitter = (s.MinJitter / 1000000).toFixed(2);
                dto.maxJitter = (s.MaxJitter / 1000000).toFixed(2);
                dto.averageSendTime = (s.AverageSendTime / 1000000).toFixed(2);
                dto.minSendTime = (s.MinSendTime / 1000000).toFixed(2);
                dto.maxSendTime = (s.MaxSendTime / 1000000).toFixed(2);

                // Any send time to display (real light controller connected)?
                dto.showSendTimes = s.SendCount > 0;

                // Update high water marks
                if (dto.latchedMinFrameTime === undefined || parseFloat(dto.minFrameTime) < parseFloat(dto.latchedMinFrameTime)) {
                    dto.latchedMinFrameTime = dto.minFrameTime;
                }
                if (dto.latchedMaxFrameTime === undefined || parseFloat(dto.maxFrameTime) > parseFloat(dto.latchedMaxFrameTime)) {
                    dto.latchedMaxFrameTime = dto.maxFrameTime;
                }

                if (dto.latchedMinJitter === undefined || parseFloat(dto.minJitter) < parseFloat(dto.latchedMinJitter)) {
                    dto.latchedMinJitter = dto.minJitter;
                }
                if (dto.latchedMaxJitter === undefined || parseFloat(dto.maxJitter) > parseFloat(dto.latchedMaxJitter)) {
                    dto.latchedMaxJitter = dto.maxJitter;
                }

                if (dto.showSendTimes) {
                    if (dto.latchedMinSendTime === undefined || parseFloat(dto.minSendTime) < parseFloat(dto.latchedMinSendTime)) {
                        dto.latchedMinSendTime = dto.minSendTime;
                    }
                    if (dto.latchedMaxSendTime === undefined || parseFloat(dto.maxSendTime) > parseFloat(dto.latchedMaxSendTime)) {
                        dto.latchedMaxSendTime = dto.maxSendTime;
                    }
                }

                // We could use magic : true but that would cause an update for every property set
                // so instead we do ractive.set(dto) to update once after everything has been set
                ractive.set(dto);
            };

            // Disable scrolling on everything
            scroll.disable(document.body);
        }
        init();
    });
