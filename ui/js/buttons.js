/* global require */
// Paths are relative to buttons.html
require([
        "./require/domReady!",
        "../js/ractive/ractive.js",
        "../js/ractive/ractive-events-tap.js",
        "../js/lib/lights.js",
        "../js/lib/scroll",
        "../js/lib/nav"],
    function(doc,R,tap,lights,scroll,nav) {
        'use strict';

        // Ractive data binding object
        var dto = {
            buttons: undefined,
            editMode: false
        };

        // Page initialisation
        var init = function () {
            // json description of button mappings
            dto.buttons = {
                leftColumn: [
                    {id: "l1", name: "OFF", action: "allLights", params: "000000"},
                    {id: "l2", name: "6f16d4", action: "allLights", params: "6f16d4"},
                    {id: "l3", name: "c8721f", action: "allLights", params: "c8721f"},
                    {id: "l4", name: "d71e1e", action: "allLights", params: "d71e1e"}
                ],
                midColumn: [
                    {id: "m1", name: "Full", action: "allLights", params: "ffffff"},
                    {id: "m2", name: "High", action: "allLights", params: "e0e0e0"},
                    {id: "m3", name: "Mid", action: "allLights", params: "808080"},
                    {id: "m4", name: "Low", action: "allLights", params: "3f3f3f"},
                    {id: "m5", name: "Very Low", action: "allLights", params: "101010"},
                    {id: "m6", name: "Nearly Off", action: "allLights", params: "020202"}
                ],
                rightColumn: [
                    {id: "r1", name: "Sweet Shop", action: "sweetshop"},
                    {id: "r2", name: "Runner", action: "runner"},
                    {id: "r3", name: "Rainbow", action: "rainbow"},
                    {id: "r4", name: "Cylon", action: "cylon"}
                ]
            };

            // Helper to find a button by id
            var findButtonById = function (id) {
                var columns = ["leftColumn", "midColumn", "rightColumn"];
                for (var i = 0; i < columns.length; i++) {
                    var buttons = dto.buttons[columns[i]];
                    for (var b = 0; b < buttons.length; b++) {
                        if (buttons[b].id === id) {
                            return buttons[b];
                        }
                    }
                }
                return undefined;
            };

            // Read write mode gets OK and Edit buttons
            var isReadWrite = location.search.split('rw=')[1];
            if (isReadWrite) {
                dto.buttons.leftColumn.push({id: "edit", name: "EDIT", action: "action-edit"});
                dto.buttons.rightColumn.push({id: "ok", name: "OK", action: "action-back"});
            }

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

            // Set up an on tap handler for all the buttons
            ractive.on('buttonHandler', function (event) {
                var id = event.node.id;
                var button = findButtonById(id);
                if (button === undefined) {
                    return;
                }
                if (button.action === "action-edit") {
                    dto.editMode = !dto.editMode;
                    ractive.set(dto);
                }
                else if (button.action === "action-back") {
                    window.location.href = "./index.html";
                }
                else if (button.action === "allLights") {
                    if (dto.editMode) {
                        nav.call("./colourpicker.html", "./buttons.html?rw=true", {buttonId: id, colour: button.params});
                    }
                    else {
                        lights.allLights(button.params);
                    }
                }
                else {
                    if (dto.editMode) {
                        // TODO edit mode for animations
                    }
                    else {
                        lights.animation(button.action);
                    }
                }
            });

            // TODO: Update UI when Teensy connection status changes
            //lights.cbConnectionStatusChanged = function() {
            //    document.getElementById("status").innerHTML =  lights.lastCallStatus ? "OK" : "Not Connected";
            //}

            // Allow button column scrolling
            scroll.enable(document.getElementsByClassName("left-column")[0]);
            scroll.enable(document.getElementsByClassName("mid-column")[0]);
            scroll.enable(document.getElementsByClassName("right-column")[0]);

            // Disable scrolling on everything by default
            scroll.disable(document.body);

            // If we are returning from colour selection set new colour here
            var p = nav.getParam();
            if (p !== undefined) {
                var button = findButtonById(p.buttonId);
                if (button !== undefined) {
                    button.params = p.colour;
                }
                // TODO update button id with name
            }
        };
        // Initialise the page
        init();
    }
);
