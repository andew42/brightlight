// Retrieve a json file stored on the server at path
// Calls back with success(jsonObject) or error(xhr)
export function getJson(path, success, error) {
    const xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                if (success) {
                    success(JSON.parse(xhr.responseText));
                }
            } else {
                if (error) {
                    error(xhr);
                }
            }
        }
    };
    xhr.open('GET', path, true);
    xhr.send();
}

// Save object as json on server, if obj is a string
// it is assumed to be json and is sent as is. Calls
// back with success() or error(xhr)
export function putJson(path, obj, success, error) {
    const xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                if (success) {
                    success();
                }
            } else {
                if (error) {
                    error(xhr);
                }
            }
        }
    };
    xhr.open('PUT', path, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(typeof obj === "string" ? obj : JSON.stringify(obj));
}
