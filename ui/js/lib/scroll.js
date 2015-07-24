// iPhone element scroll enable / disable
/* global define */
define (['./util'], function(util) {
    'use strict';
    return {
        // disable scrolling on this element
        disable : function(e) {
            // Allow single elements or arrays with applyTo
            util.applyTo(e, function(elmt) {
                elmt.addEventListener('touchmove', function (evt) {
                    // Stop everything scrolling unless it is overridden
                    if (!evt._isScroller) {
                        evt.preventDefault();
                    }
                });
            });
        },

        // enable scroll on this element
        enable : function(e) {
            // Allow single elements or arrays with applyTo
            util.applyTo(e, function(elmt) {
                // range controls are a special case
                if (elmt.nodeName === "INPUT" && elmt.type === "range") {
                    // touchmove handler - unconditional enable
                    elmt.addEventListener('touchmove', function(evt) {
                        evt._isScroller = true;
                    });
                }
                else {
                    // touchstart handler
                    elmt.addEventListener('touchstart', function() {
                        var top = elmt.scrollTop,
                            totalScroll = elmt.scrollHeight,
                            currentScroll = top + elmt.offsetHeight;

                        // If we're at the top or the bottom of the containers scroll, push up or down
                        // one pixel. This prevents the scroll from "passing through" to the body.

                        if(top === 0) {
                            elmt.scrollTop = 1;
                        }
                        else if(currentScroll === totalScroll) {
                            elmt.scrollTop = top - 1;
                        }
                    });

                    // touchmove handler
                    elmt.addEventListener('touchmove', function(evt) {
                        // If content is actually scrollable, i.e. the content is long enough that scrolling can occur
                        if(elmt.offsetHeight < elmt.scrollHeight) {
                            evt._isScroller = true;
                        }
                    });
                }
            });
        }
    };
});
