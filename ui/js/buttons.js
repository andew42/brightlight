// Ractive data binding object
var dto = {
    buttons : undefined,
    editMode : false
};

// Page initialisation
var init = function() {
    // json description of button mappings
    dto.buttons = {
        leftColumn : [
            {name:"OFF", action:"allLights", params:"000000"},
            {name:"6f16d4", action:"allLights", params:"6f16d4"},
            {name:"c8721f", action:"allLights", params:"c8721f"},
            {name:"d71e1e", action:"allLights", params:"d71e1e"}
        ],
        midColumn : [
            {name:"Full", action:"allLights", params:"ffffff"},
            {name:"High", action:"allLights", params:"e0e0e0"},
            {name:"Mid", action:"allLights", params:"808080"},
            {name:"Low", action:"allLights", params:"3f3f3f"},
            {name:"Very Low", action:"allLights", params:"101010"},
            {name:"Nearly Off", action:"allLights", params:"020202"}
        ],
        rightColumn : [
            {name:"Sweet Shop", action:"sweetshop", params:""},
            {name:"Runner", action:"runner", params:""},
            {name:"Rainbow", action:"rainbow", params:""},
            {name:"Cylon", action:"cylon", params:""}
        ]
    };

    // Read write mode gets OK and Edit buttons
    var isReadWrite = location.search.split('rw=')[1];
    if (isReadWrite) {
        dto.buttons.leftColumn.push({name:"EDIT", action:"action-edit", params:0, editMode:false});
        dto.buttons.rightColumn.push({name:"OK", action:"action-back", params:0, editMode:false});
    }

    var ractive = new Ractive({
        // The `el` option can be a node, an ID, or a CSS selector.
        el: 'container',

        // We could pass in a string, but for the sake of convenience
        // we're passing the ID of the <script> tag above.
        template: '#template',

        // Here, we're passing in some initial data
        data: dto
    });

    // Set up an on tap handler for all the buttons
    ractive.on( 'buttonHandler', function ( event ) {
        var action = event.node.attributes["button-action"].value;
        var params = event.node.attributes["button-params"].value;
        if (action == "action-edit") {
            dto.editMode = !dto.editMode;
            ractive.set(dto);
        }
        else if (action == "action-back") {
            window.location.href = "./index.html";
        }
        else if (action == "allLights") {
            lights.allLights(params);
        }
        else {
            lights.animation(action);
        }
    });

    // TODO: Update UI when connection status changes
    //lights.cbConnectionStatusChanged = function() {
    //    document.getElementById("status").innerHTML =  lights.lastCallStatus ? "OK" : "Not Connected";
    //}

    // Allow button column scrolling
    scroll.enable(document.getElementsByClassName("left-column")[0]);
    scroll.enable(document.getElementsByClassName("mid-column")[0]);
    scroll.enable(document.getElementsByClassName("right-column")[0]);

    // Disable scrolling on everything by default
    scroll.disable(document.body);
};

// Initialise the page
window.onload = function () { init(); };
