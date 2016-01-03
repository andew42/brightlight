// Static config data
/* global define */
define(function () {
    'use strict';
    return {
        // List of possible animation with display name and action name
        animations: [
            {name: "Off", action: "off"},
            {name: "Static", action: "static"},
            {name: "Cylon", action: "cylon"},
            {name: "Rainbow", action: "rainbow"},
            {name: "Runner", action: "runner"},
            {name: "Sweet Shop", action: "sweetshop"},
            {name: "Candle", action: "candle"},
            {name: "Christmas", action: "christmas"}
        ],
        // Logical segments for constructing scenes TODO GENERATE FROM GO NamedSegments
        segments: [
            {name: "All", id: "s0"},
            {name: "Bedroom", id: "s1"},
            {name: "Bathroom", id: "s2"},
            {name: "Curtains", id: "s3"},
            {name: "Test 4", id: "s4"},
            {name: "Test 5", id: "s5"},
        ]
    };
});
