// Opens an SSE stream at /api/<url>. Returns the EventSource so the caller
// can close it. Calls cb with the parsed JSON payload of each event.
export function OpenWebSocket(url, cb) {

    const fullUrl = '/api/' + url;
    console.debug('Opening SSE stream at ' + fullUrl);

    const es = new EventSource(fullUrl);

    es.onopen = () => console.debug(url + ' SSE open');
    es.onerror = () => console.error(url + ' SSE error');
    es.onmessage = evt => {
        console.debug(url + ' SSE message');
        cb(JSON.parse(evt.data));
    };

    // Mirror the WebSocket close() interface so callers need no changes
    es.close = es.close.bind(es);

    return es;
}
