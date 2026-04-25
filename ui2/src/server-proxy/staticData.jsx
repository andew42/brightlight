// Get the static data (animation and segment names)
import {getJson} from "./jsonHttpRequest";

export function getStaticData(success, error) {

    getJson('/ui-config/static-data.json', success, error);
}
