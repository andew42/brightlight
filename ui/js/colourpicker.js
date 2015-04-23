// jshint -W097
"use strict";

// Ractive data binding object
var dto = {
    hue: 0,
    saturation: 100,
    value: 100,
    colour: undefined,
    hueColour: undefined,
    asHsv: undefined,
    asRgb: undefined
};

function init() {
    // Disable scrolling on everything
    scroll.disable(document.body);

    var ractive = new Ractive({
        // The `el` option can be a node, an ID, or a CSS selector.
        el: 'container',

        // We could pass in a string, but for the sake of convenience
        // we're passing the ID of the <script> tag above.
        template: '#template',

        // Here, we're passing in some initial data
        data: dto
    });

    buttons.verticalTouchSlider(document.getElementsByClassName("cp-slider-track-hue")[0], function(pos) {
        dto.hue = pos;
        updateUi();
    }, 0, 360);

    buttons.verticalTouchSlider(document.getElementsByClassName("cp-slider-track-saturation")[0], function(pos) {
        dto.saturation = pos;
        updateUi();
    }, 0, 100);

    buttons.verticalTouchSlider(document.getElementsByClassName("cp-slider-track-value")[0], function(pos) {
        dto.value = pos;
        updateUi();
    }, 0, 100);

    var updateUi = function () {
        // Compute derived fields
        var colour = tinycolor("hsv " + dto.hue + ", " + dto.saturation + "%, " + dto.value + "%");
        dto.colour = colour.toHexString();
        dto.hueColour = tinycolor("hsv " + dto.hue + ", 100%, 100%").toHexString();
        dto.asHsv = colour.toHsvString();
        dto.asRgb = colour.toRgbString();
        // Update binding
        ractive.set(dto);
    };

    updateUi();
}

// Initialise the page
window.onload = function () { init(); };
