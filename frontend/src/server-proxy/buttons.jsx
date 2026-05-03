import {getJson, putJson} from "./jsonHttpRequest";

export function getButtons(success, error) {

    getJson('/api/ui-config/user-buttons.json', success,
        function () {
            // Couldn't find user-buttons.json, try default-buttons.json
            getJson('/api/ui-config/default-buttons.json', success, error);
        }
    )
}

export function saveButtons(cfg, success, error) {
    putJson('/api/ui-config/user-buttons.json', cfg, success, error);
}
