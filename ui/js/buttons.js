/* global require, console */
// Paths are relative to buttons.html
require([
        "./require/domReady!",
        "../js/ractive/ractive",
        "../js/ractive/ractive-touch",
        "../js/lib/lights",
        "../js/lib/scroll",
        "../js/lib/nav",
        "../js/lib/util",
        "../config/static"],
    function(doc,R,touch,lights,scroll,nav,util,sc) {
        'use strict';

        // Ractive data binding object
        var dto = {
            buttons: undefined,
            error: undefined,
            editButton: undefined,
            animations: sc.animations
        };

        // The button that was last long pressed
        var buttonBeingEdited;

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

                // Copy buttons WITHOUT OK
                var cloneEditableButtons = function() {
                    var copy = {};
                    var columns = ["leftColumn", "midColumn", "rightColumn"];
                    for (var i = 0; i < columns.length; i++) {
                        copy[columns[i]] = [];
                        var buttons = dto.buttons[columns[i]];
                        for (var b = 0; b < buttons.length; b++) {
                            if (buttons[b].id !== 'ok') {
                                copy[columns[i]].push(buttons[b]);
                            }
                        }
                    }
                    return copy;
                };

                util.putJson('../../config/user.json', cloneEditableButtons(),
                    function () {
                        // TODO LITTLE SAVE POPUP
                    },
                    function (err) {
                        /** @namespace err.responseURL */
                        dto.error = err.responseURL + ' : ' + err.responseText;
                        ractive.set(dto);
                    });
            };

            // Helper to position menu div relative to mouse
            var positionMenu = function (menu, x, y) {
                // Ensure menu is visible so we can read width and height
                menu.style.display='block';
                var menuWidth = menu.clientWidth;
                var menuHeight = menu.clientHeight;

                // Set initial position slightly offset from x,y
                var menuX = Math.max(0, x - 10);
                var menuY = Math.max(0, y - 10);

                // Now ensure menu fits entirely within the viewport
                menuX = Math.min(menuX, window.innerWidth - menuWidth);
                menuY = Math.min(menuY, window.innerHeight - menuHeight);

                // Set menu position
                menu.style.top = menuY + 'px';
                menu.style.left = menuX + 'px';
            };

            // Read write mode gets OK button
            var isReadWrite = location.search.split('rw=')[1];
            if (isReadWrite && dto.buttons !== undefined) {
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
            ractive.on('tapHandler', function (event) {
                // Close edit menu
                var editMenu = document.getElementById("edit-menu");
                editMenu.style.display='none';

                // Check for edit context menu buttons
                if (event.node.id === 'menu-name') {
                    console.info('NAME:' + buttonBeingEdited.name);
                }
                else if (event.node.id === 'menu-white') {
                    console.info('WHITE');
                }
                else if (event.node.id === 'menu-colour') {
                    console.info('COLOUR');
                    // Colour picker cares about colour, other parameters are just passed back
                    nav.call("./colourpicker.html", "./buttons.html?rw=true", {
                        colour: dto.editButton.params,
                        editButton: dto.editButton,
                        menuPos: {x: editMenu.style.left, y: editMenu.style.top}
                    });
                }
                else if (event.node.id === 'menu-ok') {
                    console.info('OK');
                    util.hideKeyboard();
                    // Update button info
                    var buttonToUpdate = findButtonById(dto.editButton.id);
                    buttonToUpdate.name = dto.editButton.name;
                    buttonToUpdate.params = dto.editButton.params;
                    buttonToUpdate.action = dto.editButton.action;
                    ractive.set(dto);
                    // Save changes to server
                    saveButtonLayout();
                }
                else if (event.node.id === 'menu-cancel') {
                    console.info('CANCEL');
                    util.hideKeyboard();
                }
                
                // Lookup dynamic lighting buttons
                var id = event.node.id;
                var button = findButtonById(id);
                if (button === undefined) {
                    return;
                }
                // OK Button?
                if (button.action === "action-back") {
                    window.location.href = "./index.html";
                }
                // Single colour light button?
                else if (button.action === "allLights") {
                    lights.allLights(button.params);
                }
                // Animation button
                else {
                    lights.animation(button.action);
                }
            });

            // Set up an on long press handler for all the buttons
            ractive.on('pressHandler', function (event) {
                var id = event.node.id;
                buttonBeingEdited = findButtonById(id);
                if (buttonBeingEdited === undefined) {
                    return true;
                }
                // Make a copy so user can cancel edit
                dto.editButton = {
                    "name": buttonBeingEdited.name,
                    "id": buttonBeingEdited.id,
                    "action": buttonBeingEdited.action,
                    "params": buttonBeingEdited.params
                };
                ractive.set(dto);
                console.info("Long Press " + buttonBeingEdited.name);
                var editMenu = document.getElementById("edit-menu");
                positionMenu(editMenu, event.original.srcEvent.pageX, event.original.srcEvent.pageY);
                return false;
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
            // p will be undefined if user canceled dialog
            if (p !== undefined) {
                // Retrieve the edit button info and update menu bindings
                dto.editButton = p.editButton;
                dto.editButton.params = p.colour;
                dto.editButton.action = "allLights";
                ractive.set(dto);
                // Show context menu at correct position
                var editMenu = document.getElementById("edit-menu");
                editMenu.style.display='block';
                editMenu.style.left = p.menuPos.x;
                editMenu.style.top = p.menuPos.y;
            }
        };

        // Initialise the page after loading button layout
        util.getJson('../../config/user.json',
            function (buttons) {
                dto.buttons = buttons;
                init();
            },
            function () {
                // Couldn't find user.json, try default.json
                util.getJson('../../config/default.json',
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
