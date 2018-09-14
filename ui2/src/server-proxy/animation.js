let currentAnimationRequest;
let nextAnimationRequest;

export function runAnimation(button) {

    // Return immediately if nothing to do
    if (button === undefined || button.segments === undefined || button.segments.length === 0)
        return;

    // Already requesting current animation? This can happen on a desktop
    // browser where we see tap for the touch and the tap for the mouse
    if (currentAnimationRequest !== undefined && button.key === currentAnimationRequest.key)
        return;

    console.log('runAnimation: ' + button.key + ' ' + button.name);

    // Throttle calls to server when user is swiping around
    if (currentAnimationRequest !== undefined) {
        // Remember button to run when last call completes
        nextAnimationRequest = button;
    } else {
        currentAnimationRequest = button;
        let req = new XMLHttpRequest();
        req.open("POST", "/RunAnimations/");
        req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
        req.send(JSON.stringify(button));
        req.onreadystatechange = function () {
            if (req.readyState === 4) {
                console.info('runAnimation returned teensy connection state: ' + req.responseText);
                req = null;
            }

            // Run any pending animations so we always have the last one
            currentAnimationRequest = undefined;
            if (nextAnimationRequest !== undefined) {
                runAnimation(nextAnimationRequest);
                nextAnimationRequest = undefined;
            }
        };
    }
}
