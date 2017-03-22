// Static config data
/* global define */
define(function () {
    'use strict';
    return {
        // List of possible animation
        animations: [
            "Off",
            "Static",
            "Cylon",
            "Rainbow",
            "Runner",
            "Sweet Shop",
            "Candle",
            "Christmas",
            "Pulse",
            "Twinkle",
            "BabyBows",
            "Test"
        ],
        // Logical segments for constructing scenes TODO GENERATE FROM GO NamedSegments
        segments: [
            "All",
            "All Ceiling",
            "All Wall",

            "Bedroom",
            "Bedroom Ceiling",
            "Bedroom Wall",

            "Bathroom",
            "Bathroom Ceiling",
            "Bathroom Wall",

            "Curtains",

            "Strip Three",
            "Strip Five"
        ]
    };
});
