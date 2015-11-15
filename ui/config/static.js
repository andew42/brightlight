// Static config data
/* global define */
define (function() {
    'use strict';
    return {
        // List of possible animation with display name and action name
        animations: [
            {name: "Cylon", action: "cylon"},
            {name: "Rainbow", action: "rainbow"},
            {name: "Runner", action: "runner"},
            {name: "Sweet Shop", action: "sweetshop"},
            {name: "Candle", action: "candle"},
            {name: "Christmas", action: "christmas"}
        ],
        // Logical segments for constructing scenes
        segments: [
            "Bedroom",
            "Bathroom"
        ]
    };
});
