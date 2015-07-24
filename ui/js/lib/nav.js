/* global define */
define (function() {
    'use strict';
    return {
        call : function (targetUri, returnUri, param) {
            window.sessionStorage.setItem("navCallReturnUri", returnUri);
            window.sessionStorage.setItem("navCallParam", JSON.stringify(param));
            window.location.href = targetUri;
        },

        getParam : function() {
            var p = window.sessionStorage.getItem("navCallParam");
            return (p === undefined || p === null) ? undefined : JSON.parse(p);
        },

        ret : function (param) {
            var returnUri = window.sessionStorage.getItem("navCallReturnUri");
            if (param !== undefined) {
                window.sessionStorage.setItem("navCallParam", JSON.stringify(param));
            }
            else {
                window.sessionStorage.removeItem("navCallParam");
            }
            window.location.href = returnUri;
        }
    };
});
