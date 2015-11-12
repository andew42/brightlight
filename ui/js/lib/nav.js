/* global define */
define (function() {
    'use strict';
    var paramName = "paramName";
    var returnUriName = "navCallReturnUri";
    return {
        call : function (targetUri, returnUri, param) {
            window.sessionStorage.setItem(returnUriName, returnUri);
            window.sessionStorage.setItem(paramName, JSON.stringify(param));
            window.location.href = targetUri;
        },

        getParam : function() {
            var p = window.sessionStorage.getItem(paramName);
            window.sessionStorage.removeItem(paramName);
            return (p === undefined || p === null) ? undefined : JSON.parse(p);
        },

        ret : function (param) {
            var returnUri = window.sessionStorage.getItem(returnUriName);
            if (param !== undefined) {
                window.sessionStorage.setItem(paramName, JSON.stringify(param));
            }
            else {
                window.sessionStorage.removeItem(paramName);
            }
            window.location.href = returnUri;
        }
    };
});
