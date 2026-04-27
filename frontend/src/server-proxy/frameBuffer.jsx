// Returns a web socket for a frame buffer connection
// The caller should close returned socket when done
// Calls cb with each frame buffer frame every 40ms
import {OpenWebSocket} from "./webSocket";

export function Open(cb) {

    // Convert a number to a 6 digit zero padded colour
    function numberToColourCode(n) {
        let c = n.toString(16);
        while (c.length < 6) {
            c = '0' + c;
        }
        return c;
    }

    return OpenWebSocket('FrameBuffer', fb => {

        // We get called with the framebuffer payload as JSON
        let strips = fb.Strips;

        // Rebuild led array with first few (20) leds from each strip
        let rc = [];
        for (let s = 0; s < strips.length; s++) {
            let fullStrip = strips[s];
            let leds = fullStrip.Leds;
            let displayStrip = [];
            let limitStripCount = 20;
            for (let l = 0; l < leds.length; l++) {
                if (limitStripCount-- <= 0) {
                    break;
                }
                displayStrip.push(numberToColourCode(leds[l]));
            }
            rc.push(displayStrip);
        }
        cb(rc);
    });
}
