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
    function (doc, R, touch, lights, scroll, nav, util, sc) {
        'use strict';

        // Ractive data binding object
        var dto = {
            buttons: undefined,
            error: undefined,
            editButton: undefined,
            animations: sc.animations,
            segments: undefined,
            selectedSegment: undefined
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
            var saveButtonLayout = function () {
                // Copy buttons WITHOUT OK
                var cloneEditableButtons = function () {
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
                menu.style.display = 'block';
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

            // Helper returns:
            // { "name": segment display name (for segmentId) with active animation name appended,
            //   "index": within the active segments iff the segment is active (so we can put active segments first)
            //   "animation": the animation iff the segment has and active animation}
            var getSegmentWithAnimationName = function (segmentId, activeSegments) {
                // Look up the segment display name, given it's id
                var segmentName = util.findFirst(sc.segments,
                    function (seg) {return seg.id === segmentId;},
                    function (seg) {return seg.name;});

                // Give up if the segment id is unknown
                if (segmentName === undefined) {
                    return segmentName;
                }

                // Determine if the segment has an animation running
                var activeSegmentIndex = util.findFirstIndex(activeSegments,
                    function (seg) {return seg.segmentId === segmentId;});

                // Not active so just return segment name
                if (activeSegmentIndex === undefined) {
                    return {"name": segmentName};
                }

                // Otherwise look up the actions display name
                var action = activeSegments[activeSegmentIndex].action;
                var animationIndex = util.findFirstIndex(sc.animations,
                    function (animation) {return animation.action === action;});

                if (animationIndex === undefined) {
                    return {"name": segmentName};
                }

                // Finally return the composite "segment name (animation name)" index and animation
                return {"name": segmentName + " (" + sc.animations[animationIndex].name + ")",
                        "index": activeSegmentIndex,
                        "animation": sc.animations[animationIndex]};
            };

            // Takes a button's active segment list and builds the full segment
            // list which will include animation names for the active segments
            var buildSegmentList = function(buttonsActiveSegments) {
                // Build a custom list of all segment with appended actions specific for this button
                var segments = [];
                for (var s = 0; s < sc.segments.length; s++) {
                    var id = sc.segments[s].id;

                    // Retrieve entry and insert into list respecting index (or append if no index)
                    var si = getSegmentWithAnimationName(id, buttonsActiveSegments);
                    if (si !== undefined) {
                        if (si.index !== undefined) {
                            // Here to insert at a specific index
                            for (var i = 0; i < segments.length; i++) {
                                if (segments[i].index === undefined || segments[i].index > si.index) {
                                    break;
                                }
                            }
                            // Must be a copy as we bind this to a combo which will change action
                            var animation = {"name":si.animation.name, "action":si.animation.action}
                            segments.splice(i, 0, {"name":si.name, "id":id, "animation":animation});
                        }
                        else {
                            // Here to just append
                            segments.push({"name": si.name, "id": id});
                        }
                    }
                }
                return segments;
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

            // Observe selection changes to context menu - segment
            ractive.observe( 'selectedSegment.id', function ( newValue ) {
                console.info(newValue);

                // Ignore initial update
                if (newValue === undefined) {
                    return;
                }

                // Update selected segment
                dto.selectedSegment = util.findFirst(dto.segments,
                    function (seg) {return seg.id === newValue;},
                    function (seg) {return seg;});

                // Force the selected animation to update to off TODO THIS PROBABLY DOESN'T WORK
                if (dto.selectedSegment.animation === undefined) {
                    dto.selectedSegment.animation = {"action": "off"};
                }

                // Update animation drop down
                ractive.set(dto);
            });

            // Observe selection changes to context menu - animation
            ractive.observe( 'selectedSegment.animation.action', function ( newValue ) {
                console.info(newValue);

                var selectedSegmentId = dto.selectedSegment.id;
                // Ignore initial update
                if (selectedSegmentId === undefined) {
                    return;
                }

                // Find the selected segment in the buttons segment list
                var selectedButtonSegmentIndex = util.findFirstIndex(dto.editButton.segments,
                    function (s) {return s.segmentId === selectedSegmentId;});

                // If the action if off, remove the segment for the button
                if (newValue === 'off') {
                    dto.editButton.segments.splice(selectedButtonSegmentIndex, 1);
                }
                else {
                    if (selectedButtonSegmentIndex === undefined) {
                        // Here if segment was off, append a new button action with TODO: set default params
                        if (dto.editButton.segments === undefined) {
                            dto.editButton.segments = [];
                        }
                        dto.editButton.segments.push(
                            {"segmentId": selectedSegmentId, "action":newValue, "params":808080});
                    } else {
                        // otherwise update the (existing) segments action TODO: set default params
                        dto.editButton.segments[selectedButtonSegmentIndex].action = newValue;
                    }
                }

                // Rebuild a custom list of all segment with appended actions specific for this button
                dto.segments = buildSegmentList(dto.editButton.segments);
                ractive.set(dto);
            });

            // Set up an on tap handler for all the buttons
            ractive.on('tapHandler', function (event) {
                // Close edit menu
                var editMenu = document.getElementById("edit-menu");
                editMenu.style.display = 'none';

                // Check for edit context menu button taps
                if (event.node.id === 'menu-name') {
                    // TODO
                    console.info('NAME:' + buttonBeingEdited.name);
                }
                else if (event.node.id === 'menu-white') {
                    // TODO
                    console.info('WHITE');
                }
                else if (event.node.id === 'menu-colour') {
                    console.info('COLOUR');
                    // Colour picker cares about colour, other parameters are just passed back
                    nav.call("./colourpicker.html", "./buttons.html?rw=true", {
                        colour: dto.editButton.params,
                        // TODO
                        editButton: dto.editButton,
                        menuPos: {x: editMenu.style.left, y: editMenu.style.top}
                    });
                }
                else if (event.node.id === 'menu-ok') {
                    console.info('OK');
                    util.hideKeyboard();
                    // Update button info TODO now for segments!!!
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

                // Lookup dynamic (user programmable) lighting buttons
                var id = event.node.id;
                var button = findButtonById(id);
                if (button === undefined) {
                    return;
                }
                // OK Button?
                if (button.action === "action-back") {
                    window.location.href = "./index.html";
                }
                else {
                    lights.runAnimations(button.segments);
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
                    "segments": []
                };
                // Copy across the zero or more active segments for this button
                var buttonsActiveSegments = buttonBeingEdited.segments;
                for (var s = 0; s < buttonsActiveSegments.length; s++) {
                    var copy = {"segmentId": buttonsActiveSegments[s].segmentId, "action": buttonsActiveSegments[s].action};
                    if (buttonsActiveSegments[s].params !== undefined) {
                        copy.params = buttonsActiveSegments[s].params;
                    }
                    dto.editButton.segments.push(copy);
                }
                // Build a custom list of all segment with appended actions specific for this button
                dto.segments = buildSegmentList(buttonsActiveSegments);
                dto.selectedSegment = dto.segments[0];
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
                editMenu.style.display = 'block';
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
