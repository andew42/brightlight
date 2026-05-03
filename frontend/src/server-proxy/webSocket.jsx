// Opens an SSE stream at /api/<url>. Returns a handle with a close() method.
// Automatically reconnects when the page becomes visible after being backgrounded.
export function OpenWebSocket(url, cb) {

    const fullUrl = '/api/' + url;
    let es = null;
    let closed = false;

    function connect() {
        if (closed) return;
        console.debug('Opening SSE stream at ' + fullUrl);
        es = new EventSource(fullUrl);
        es.onopen = () => console.debug(url + ' SSE open');
        es.onerror = () => console.error(url + ' SSE error');
        es.onmessage = evt => cb(JSON.parse(evt.data));
    }

    function onVisibilityChange() {
        if (document.visibilityState === 'visible') {
            if (es) es.close();
            connect();
        }
    }

    document.addEventListener('visibilitychange', onVisibilityChange);
    connect();

    return {
        close() {
            closed = true;
            document.removeEventListener('visibilitychange', onVisibilityChange);
            if (es) es.close();
        }
    };
}
