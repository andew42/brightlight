/* global require */
// Paths are relative to index.html
require([
        "./require/domReady!",
        "../js/ractive/ractive",
        "../js/ractive/ractive-touch",
        "../js/lib/lights",
        "../js/lib/scroll"],
    function (doc, R, touch, lights, scroll) {
        'use strict';

        // Ractive data binding object
        var dto = {
            buttons: undefined,
            outputMappings: ["Linear", "Linear 128", "Cie 256", "Cie 128"],
            selectedOutputMapping: "Linear 128"
        };

        // Page initialisation
        var init = function () {
            // json description of button mappings
            dto.buttons = {
                leftColumn: [
                    {name: "OFF", action: "allLights", params: "000000"},
                    {name: "Buttons RO", action: "action-navigate", params: "./buttons.html"},
                    {name: "Virtual", action: "action-navigate", params: "./virtual.html"}
                ],
                midColumn: [
                    {name: "High", action: "allLights", params: "e0e0e0"},
                    {name: "Buttons RW", action: "action-navigate", params: "./buttons.html?rw=true"},
                    {name: "Strip Length", action: "action-navigate", params: "./striplength.html"}
                ],
                rightColumn: [
                    {name: "Purple", action: "allLights", params: "6f16d4"},
                    {name: "Rainbow", action: "rainbow", params: ""},
                    {name: "Stats", action: "action-navigate", params: "./stats.html"}
                ]
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

            // Set up an on tap handler for all the buttons
            ractive.on('buttonHandler', function (event) {
                var action = event.node.attributes["button-action"].value;
                var params = event.node.attributes["button-params"].value;
                if (action === "action-navigate") {
                    window.location.href = params;
                } else if (action === "allLights") {
                    lights.allLights(params);
                } else {
                    lights.animation(action);
                }
            });


            ractive.observe('selectedOutputMapping', function(newMapping) {
                console.info(newMapping);
                var req = new XMLHttpRequest();
                req.open("POST", "/option/");
                req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
                req.send(JSON.stringify({"cmd": "outputMapping", "param": newMapping}));
            });

            // TODO: Update UI when connection status changes
            //lights.cbConnectionStatusChanged = function() {
            //    document.getElementById("status").innerHTML =  lights.lastCallStatus ? "OK" : "Not Connected";
            //}

            // Allow button column scrolling
            scroll.enable(document.getElementsByClassName("left-column")[0]);
            scroll.enable(document.getElementsByClassName("mid-column")[0]);
            scroll.enable(document.getElementsByClassName("right-column")[0]);

            // Disable scrolling on everything by default
            scroll.disable(document.body);
        };

        // Initialise the page (domReady! ensures the dom is ready at this point)
        init();
    });
