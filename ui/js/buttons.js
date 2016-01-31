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

        // Page initialisation
        var init = function () {
        // -----
            // Helper to find a button by id
            var findButtonById = function (id) {
            // ---------------
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
            // -----------------
                // Copy buttons WITHOUT OK
                var cloneEditableButtons = function () {
                // ---------------------
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

                // Save to server
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
            // -------------
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

            // return {
            //   "segment": segment name,
            //   "index": within the active segments iff the segment is active (so we can put active segments first)
            //   "animation": the animation iff the segment has and active animation,
            //   "params": the animation parameters
            // }
            var getSegmentWithAnimation = function (segment, activeSegments) {
            //-------------------------
                // Determine if the segment has an animation running on this button
                var activeSegmentIndex = util.findFirstIndex(activeSegments,
                    function (seg) {return seg.segment === segment;});

                // Not active so just return the segment name
                if (activeSegmentIndex === undefined) {
                    return {"segment": segment};
                }

                // Active, return segment with animation
                return {
                    "segment": segment,
                    "index": activeSegmentIndex,
                    "animation": activeSegments[activeSegmentIndex].animation,
                    "params": activeSegments[activeSegmentIndex].params
                };
            };

            // Takes a button's active segment list and builds the full
            // segment list which will include any active button animation
            var buildSegmentAnimations = function(buttonsActiveSegments) {
            // -----------------------
                // Build a custom list of all segments with animations specific for this button
                var segments = [];
                for (var s = 0; s < sc.segments.length; s++) {

                    // Retrieve entry and insert into list respecting index (or append if no index)
                    var sa = getSegmentWithAnimation(sc.segments[s], buttonsActiveSegments);
                    if (sa.index !== undefined) {

                        // Here to insert at a specific index
                        for (var i = 0; i < segments.length; i++) {
                            if (segments[i].index === undefined || segments[i].index > sa.index) {
                                break;
                            }
                        }

                        // Must be a copy as we bind this to a combo which will change action
                        segments.splice(i, 0, sa);
                    }
                    else {
                        // Here to just append
                        segments.push(sa);
                    }
                }
                return segments;
            };

            // Create the segment list for the button from segmentAnimations binding
            var buildButtonSegments = function() {
            // --------------------
                var segments = [];
                for (var i = 0; i < dto.segmentAnimations.length; i++) {
                    var sa = dto.segmentAnimations[i];
                    if (sa.animation !== undefined) {
                        segments.push({
                            "segment": sa.segment,
                            "animation": sa.animation,
                            "params": sa.params
                        });
                    }
                }
                return segments;
            };

            // Observe selection changes to context menu - segment
            var observeSelectedSegment = function ( newValue ) {
            // -----------------------
                console.info(newValue);

                // Ignore initial update
                if (newValue === undefined) {
                    return;
                }

                // Select animation attached to this segment (if any)
                dto.selectedAnimation = util.findFirst(dto.segmentAnimations,
                    function (seg) {return seg.segment === newValue;},
                    function (seg) {return seg.animation;});

                // Force the selected animation to update to off if undefined
                if (dto.selectedAnimation === undefined) {
                    dto.selectedAnimation = "Off";
                }

                // Update UI
                ractive.set(dto);
            };

            // Observe selection changes to context menu - animation
            var observeSelectedAnimation = function ( newValue ) {
            // -------------------------
                console.info(newValue);

                // Ignore initial update
                if (dto.segmentAnimations === undefined) {
                    return;
                }

                // Find selected segment animations
                var sa = util.findFirst(dto.segmentAnimations,
                    function (s) {return s.segment === dto.selectedSegment;},
                    function (s) {return s;});

                // Update animation removing if Off (TODO: Set animation default params)
                sa.animation = newValue === "Off" ? undefined : newValue;

                // The segmentAnimation list may now be out of order (i.e. all
                // animated segments may not be at the top) so we rebuild list
                var buttonSegments = buildButtonSegments();
                dto.segmentAnimations = buildSegmentAnimations(buttonSegments);

                // Make sure room lights reflect new state (especially the all off state)
                lights.runAnimations(buttonSegments);

                // Update UI
                ractive.set(dto);
            };

            // Set up an on tap handler for all the buttons
            var tapHandler = function (event) {
            // -----------

                var editMenu = document.getElementById("edit-menu");

                // Check for edit context menu button taps
                if (event.node.id === 'menu-colour') {

                    console.info('COLOUR');

                    // Get the segment we are configuring, may be undefined
                    var selectedSegment = util.findFirst(dto.segmentAnimations,
                        function(s) {return s.segment === dto.selectedSegment;},
                        function(s) {return s;});

                    // Currently we only support colour for static segments
                    if (selectedSegment === undefined || selectedSegment.animation !== "Static") {
                        return;
                    }

                    // Colour picker cares about colour and segment,
                    // other parameters are just passed back
                    nav.call("./colourpicker.html", "./buttons.html?rw=true", {
                        colour: selectedSegment.params,
                        segment: selectedSegment.segment,
                        // Used to restore the edit menu so far
                        editButton: dto.editButton,
                        segmentAnimations: dto.segmentAnimations,
                        selectedSegment: dto.selectedSegment,
                        selectedAnimation: dto.selectedAnimation,
                        // Used to restore the edit menu position
                        menuPos: {x: editMenu.style.left, y: editMenu.style.top}
                    });
                }
                else if (event.node.id === 'menu-ok') {

                    console.info('OK');
                    util.hideKeyboard();

                    // Close edit menu
                    editMenu.style.display = 'none';

                    // Find the button we are editing and update it
                    var buttonToUpdate = findButtonById(dto.editButton.id);

                    // Update the button segment list from segmentAnimations menu
                    buttonToUpdate.segments = buildButtonSegments();

                    // If all segments are off, replace with all off animation
                    if (buttonToUpdate.segments === undefined || buttonToUpdate.segments.length === 0) {
                        buttonToUpdate.segments = [{"segment": "All", "animation": "Static", "params": "000000"}];
                    }

                    // Update the button name if it changed
                    if (buttonToUpdate.name !== dto.editButton.name) {
                        buttonToUpdate.name = dto.editButton.name;
                        ractive.set(dto);
                    }

                    // Make sure room lights reflect new state (especially the all off state)
                    lights.runAnimations(buttonToUpdate.segments);

                    // Save changes to server
                    saveButtonLayout();
                }
                else if (event.node.id === 'menu-cancel') {

                    // Restore room lights to pre edit state
                    var b = findButtonById(dto.editButton.id);
                    if (b !== undefined && b.segments !== undefined) {
                        lights.runAnimations(b.segments);
                    }

                    console.info('CANCEL');
                    util.hideKeyboard();

                    // Close edit menu
                    editMenu.style.display = 'none';
                }

                // If the edit menu is visible, other buttons are disabled
                if (editMenu.style.display !== 'none') {
                    return;
                }

                // Lookup dynamic (user programmable) lighting buttons
                var id = event.node.id;
                var button = findButtonById(id);
                if (button === undefined) {
                    return;
                }
                // OK Button or lighting button?
                if (button.action === "action-back") {
                    window.location.href = "./index.html";
                }
                else {
                    lights.runAnimations(button.segments);
                }
            };

            // Set up an on long press handler for all the buttons
            var pressHandler = function (event) {
            // -------------
                // Find the button being pressed
                var id = event.node.id;
                var buttonBeingEdited = findButtonById(id);
                if (buttonBeingEdited === undefined) {
                    return true;
                }

                console.info("Long Press " + buttonBeingEdited.name);

                // Make sure the room lights reflect edit button
                lights.runAnimations(buttonBeingEdited.segments);

                // Make a copy of id and name, so user can cancel edit
                dto.editButton = {
                    "name": buttonBeingEdited.name,
                    "id": buttonBeingEdited.id
                };

                // Build a list of all segments with selected buttons animations attached
                dto.segmentAnimations = buildSegmentAnimations(buttonBeingEdited.segments);

                // Select the first segment segment
                dto.selectedSegment = dto.segmentAnimations[0].segment;

                // Select first segment animation
                dto.selectedAnimation = dto.segmentAnimations[0].animation;

                // Update the UI
                ractive.set(dto);

                // Show and position the edit menu near pressed button
                var editMenu = document.getElementById("edit-menu");
                positionMenu(editMenu, event.original.srcEvent.pageX, event.original.srcEvent.pageY);
                return false;
            };

            // --- Init() --- Execution proper starts here

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

                // Retrieve the edit button info we sent to colour editor
                dto.editButton = p.editButton;
                dto.segmentAnimations = p.segmentAnimations;
                dto.selectedSegment = p.selectedSegment;
                dto.selectedAnimation = p.selectedAnimation;

                // Colour picker returns undefined colour for cancel
                if (p.colour !== undefined) {

                    // Get the colour of the selected button, may be undefined
                    var selected = util.findFirst(dto.segmentAnimations,
                        function (s) {return s.segment === dto.selectedSegment;},
                        function (s) {return s;});

                    if (selected !== undefined) {
                        selected.params = p.colour;
                    }
                }

                // Update the UI
                ractive.set(dto);

                // Show context menu at correct position
                var editMenu = document.getElementById("edit-menu");
                editMenu.style.display = 'block';
                editMenu.style.left = p.menuPos.x;
                editMenu.style.top = p.menuPos.y;
            }

            // Hook up ractive handlers
            ractive.observe('selectedSegment', observeSelectedSegment);
            ractive.observe('selectedAnimation', observeSelectedAnimation);
            ractive.on('tapHandler', tapHandler);
            ractive.on('pressHandler', pressHandler);
        };

        // --- Page Execution Starts Here ---

        // Ractive data binding object
        var dto = {
            // See default.json for the schema
            buttons: undefined,

            // Fatal error string returned from server
            error: undefined,

            // A copy of selected button id and name
            // "id": "l1",
            // "name": "OFF",
            editButton: undefined,

            // A list of animation names
            // ["Off", "Static", "Rainbow"]
            animations: sc.animations,

            // A list of all segments with animation for edit button attached
            // [
            //   {
            //     "segment": "All"
            //   },
            //   {
            //     "segment": "Bedroom",
            //     "animation": "Rainbow",
            //     "params": "000000"
            //   }
            // ]
            segmentAnimations: undefined,

            // A segment name e.g. All
            selectedSegment: undefined,

            // An animation name e.g. Rainbow
            selectedAnimation: undefined
        };

        // Initialise the page after loading button layout
        util.getJson('../../config/user.json',
        //----------
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
