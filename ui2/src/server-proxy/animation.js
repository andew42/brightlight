// Run an animation
let busy;
let nextRunAnimation;

export function runAnimation(segments) {

    console.log(segments);

    // Return immediately if nothing to do
    if (segments === undefined || segments.length === 0) {
        return;
    }

    // Throttle calls to server when user is swiping around
    if (busy) {
        // Make a copy of the segment list to run when last call completes
        nextRunAnimation = [];
        for (let s = 0; s < segments.length; s++) {
            nextRunAnimation.push({
                "segment": segments[s].segment,
                "animation": segments[s].animation,
                "params": segments[s].params
            });
        }
    } else {
        busy = true;
        let req = new XMLHttpRequest();
        req.open("POST", "/RunAnimations/");
        req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
        req.send(JSON.stringify(segments));
        req.onreadystatechange = function () {
            if (req.readyState === 4) {
// TODO                setConnectionStatus(req.responseText);
                req = null;
            }

            // Run any pending animations so we always have the last one
            busy = false;
            if (nextRunAnimation !== undefined) {
                runAnimation(nextRunAnimation);
                nextRunAnimation = undefined;
            }
        };
    }
}
