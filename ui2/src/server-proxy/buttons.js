import {getJson, putJson} from "./jsonHttpRequest";

// Get the list of user defined buttons (or the default set)
export function getButtons(success, error) {

    getJson('/ui-config/user-buttons.json', success,
        function () {
            // Couldn't find user-buttons.json, try default-buttons.json
            getJson('/ui-config/default-buttons.json', success, error);
        }
    )
}

// Save the supplied user defined buttons
export function saveButtons(buttons, success, error) {
    putJson('/ui-config/user-buttons.json', buttons, success, error);
}
