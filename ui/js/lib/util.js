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
        }
    };
});
