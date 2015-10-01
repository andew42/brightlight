/* global define, HTMLCollection */
define (function() {
    'use strict';
    return {
        // Set default function parameter if parameter undefined
        defaultFor : function (arg, val) {
            return typeof arg !== 'undefined' ? arg : val;
        },

        // Apply a function to a single element or array
        applyTo : function (arg, f) {
            if (arg.constructor === Array || arg.constructor === HTMLCollection) {
                for (var i = 0; i < arg.length; i++) {
                    f(arg[i]);
                }
            }
            else {
                f(arg);
            }
        },

        // Load json
        getJson : function (path, success, error)
        {
            var xhr = new XMLHttpRequest();
            xhr.onreadystatechange = function()
            {
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
        },

        // Save object as json on server, if obj is a string
        // it is assumed to be json and is sent as is
        putJson : function (path, obj, success, error)
        {
            var xhr = new XMLHttpRequest();
            xhr.onreadystatechange = function()
            {
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
    };
});
