// Ractive binding object
var dto = {
    stripIndex: 0,
    stripLength: 100
};

function init() {

    var updateConfig = function() {
        var req = new XMLHttpRequest();
        req.open('PUT', '/Config/' + dto.stripIndex + ',' + dto.stripLength, true);
        req.send();
    };

    var ractive = new Ractive({
        // The `el` option can be a node, an ID, or a CSS selector.
        el: 'container',

        // We could pass in a string, but for the sake of convenience
        // we're passing the ID of the <script> tag above.
        template: '#template',

        // Here, we're passing in some initial data
        data: dto
    });

    // Track user changes and update lights
    ractive.observe('stripIndex', updateConfig);
    ractive.observe('stripLength', updateConfig);

    // Wire up touch sliders to adjust index and length
    buttons.touchSlider(document.getElementById("si"),
        function(v) {dto.stripIndex = v; ractive.set(dto);},
        dto.stripIndex, 0, 15);
    buttons.touchSlider(document.getElementById("sl"),
        function(v) {dto.stripLength = v; ractive.set(dto);},
        dto.stripLength, 0, 1000);

    // Wire up the OK button
    buttons.setOnTapNavigationMapping(ractive, 'btnGoBack', 'OK', './index.html');

    // Disable scrolling on everything
    //scroll.disable(document.body);
}

// Wait to page load to initialise
window.onload = init();
