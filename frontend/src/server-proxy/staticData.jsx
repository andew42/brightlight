// Get the static data (animation and segment names)
import {getJson} from "./jsonHttpRequest";

export function getStaticData(success, error) {

    getJson('/api/ui-config/static-data.json', success, error);
}
