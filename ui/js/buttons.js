/* global require */
// Paths are relative to buttons.html
require([
        "./require/domReady!",
        "../js/ractive/ractive",
        "../js/ractive/ractive-touch",
        "../js/lib/lights",
        "../js/lib/scroll",
        "../js/lib/nav",
        "../js/lib/util"],
    function(doc,R,touch,lights,scroll,nav,util) {
        'use strict';

        // Ractive data binding object
        var dto = {
            buttons: undefined,
            editMode: false,
            error: undefined
        };

        // Page initialisation
        var init = function () {

            // Helper to find a button by id
            var findButtonById = function (id) {
                if (dto.buttons === undefined) {
                    return;
                }
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

            // Helper to save the updated button layout back to server
            var saveButtonLayout = function() {

                // Copy buttons WITHOUT ok and edit
                var clone = function(obj) {
                    var copy = {};
                    var columns = ["leftColumn", "midColumn", "rightColumn"];
                    for (var i = 0; i < columns.length; i++) {
                        copy[columns[i]] = [];
                        var buttons = dto.buttons[columns[i]];
                        for (var b = 0; b < buttons.length; b++) {
                            if (buttons[b].id !== 'edit' && buttons[b].id !== 'ok') {
                                copy[columns[i]].push(buttons[b]);
                            }
                        }
                    }
                    return copy;
                };

                var test2 = util.putJson('../../config/user.json', clone(dto.buttons),
                    function () {
                        // TODO LITTLE SAVE POPUP
                    },
                    function (err) {
                        dto.error = err.responseURL + ' : ' + err.responseText;
                        ractive.set(dto);
                    });
            };

            // Read write mode gets OK and Edit buttons
            var isReadWrite = location.search.split('rw=')[1];
            if (isReadWrite && dto.buttons !== undefined) {
                dto.buttons.leftColumn.push({id: 'edit', name: 'EDIT', action: 'action-edit'});
                dto.buttons.rightColumn.push({id: 'ok', name: 'OK', action: 'action-back'});
            }

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

            // TODO
            document.body.classList.add("disable-user-select");

            // If we are returning from colour selection set new colour here
            var p = nav.getParam();
            if (p !== undefined) {
                var button = findButtonById(p.buttonId);
                if (button !== undefined) {
                    button.params = p.colour;
                    // Save changes to server
                    saveButtonLayout();
                }
                // TODO update button id with name
            }
        };

        // Initialise the page after loading button layout
        var test = util.getJson('../../config/user.json',
            function (buttons) {
                dto.buttons = buttons;
                init();
            },
            function () {
                // Couldn't find user.json, try default.json
                var test = util.getJson('../../config/default.json',
                    function (buttons) {
                        dto.buttons = buttons;
                        init();
                    },
                    function (err) {
                        dto.error = err.responseURL + ' : ' + err.responseText;
                        init();
                    }
                );
            }
        );
    }
);
