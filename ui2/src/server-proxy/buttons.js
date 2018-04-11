import {getJson, putJson} from "./jsonHttpRequest";

export function getButtons(success, error) {

    getJson('/ui-config/user-buttons.json', success,
        function () {
            // Couldn't find user-buttons.json, try default-buttons.json
            getJson('/ui-config/default-buttons.json', success, error);
        }
    )
}

export function saveButtons(buttons, success, error) {
    putJson('/ui-config/user-buttons.json', buttons, success, error);
}
