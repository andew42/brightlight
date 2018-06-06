export default class Colour {

    // CSS RGB Color string (e.g. #00ff00) or
    // Decimal RGB Colour number (e.g. 200) or
    // RGB Colour object ({r:100, g:200, b:300}) or
    // HSV Colour object ({h:100, s:200, v:50})
    constructor(c) {

        // CSS RGB Color string (e.g. #00ff00)?
        if (typeof c === 'string') {
            if (c[0] !== '#')
                throw new Error('Invalid string color: ' + c);
            c = c.slice(1);
            let n = parseInt(c, 16);
            if (isNaN(n))
                throw Error('Invalid string color: ' + c);
            // Drop through to process as RGB number
            c = n;
        }

        // Decimal RGB Colour number (e.g. 200)?
        if (typeof c === 'number') {
            c >>>= 0;
            this.r = (c & 0xFF0000) >>> 16;
            this.g = (c & 0xFF00) >>> 8;
            this.b = c & 0xFF;
            let hsv = Colour.rgbToHsv(this.r, this.g, this.b);
            this.h = hsv.h;
            this.s = hsv.s;
            this.v = hsv.v;
            return;
        }

        // Colour object?
        if (c.r !== undefined) {
            Colour.throwIfInvalidRgb(c);
            this.r = c.r;
            this.g = c.g;
            this.b = c.b;
            // Keep original values if possible to avoid accumulating conversion errors
            if (c.h !== undefined) {
                Colour.throwIfInvalidHsv(c);
            } else {
                c = Colour.rgbToHsv(c);
            }
            this.h = c.h;
            this.s = c.s;
            this.v = c.v;
        } else if (c.h !== undefined) {
            Colour.throwIfInvalidHsv(c);
            this.h = c.h;
            this.s = c.s;
            this.v = c.v;
            // Keep original values if possible to avoid accumulating conversion errors
            if (c.r !== undefined) {
                Colour.throwIfInvalidRgb(c);
            } else {
                c = Colour.hsvToRgb(c);
            }
            this.r = c.r;
            this.g = c.g;
            this.b = c.b;
        } else {
            throw Error('Invalid colour: ' + c);
        }
    }

    static throwIfInvalidRgb(c) {
        if (c.r < 0 || c.r > 255 ||
            c.g < 0 || c.g > 255 ||
            c.b < 0 || c.b > 255)
            throw Error('Invalid colour: ' + c);
    }

    static throwIfInvalidHsv(c) {
        if (c.h < 0 || c.h > 360 ||
            c.s < 0 || c.s > 100 ||
            c.v < 0 || c.v > 100)
            throw Error('Invalid colour: ' + c);
    }

    static rgbToHsv(r, g, b) {
        if (arguments.length === 1) {
            g = r.g;
            b = r.b;
            r = r.r;
        }
        let max = Math.max(r, g, b), min = Math.min(r, g, b),
            d = max - min,
            h,
            s = (max === 0 ? 0 : d / max),
            v = max / 255;
        switch (max) {
            case min:
                h = 0;
                break;
            case r:
                h = (g - b) + d * (g < b ? 6 : 0);
                h /= 6 * d;
                break;
            case g:
                h = (b - r) + d * 2;
                h /= 6 * d;
                break;
            case b:
            default:
                h = (r - g) + d * 4;
                h /= 6 * d;
                break;
        }
        return {
            h: h * 360,
            s: s * 100,
            v: v * 100
        };
    }

    static hsvToRgb(h, s, v) {
        if (arguments.length === 1) {
            s = h.s;
            v = h.v;
            h = h.h;
        }
        let r, g, b, i, f, p, q, t;
        h = h / 360;
        s = s / 100;
        v = v / 100;
        i = Math.floor(h * 6);
        f = h * 6 - i;
        p = v * (1 - s);
        q = v * (1 - f * s);
        t = v * (1 - (1 - f) * s);
        switch (i % 6) {
            case 0:
                r = v;
                g = t;
                b = p;
                break;
            case 1:
                r = q;
                g = v;
                b = p;
                break;
            case 2:
                r = p;
                g = v;
                b = t;
                break;
            case 3:
                r = p;
                g = q;
                b = v;
                break;
            case 4:
                r = t;
                g = p;
                b = v;
                break;
            case 5:
            default:
                r = v;
                g = p;
                b = q;
                break;
        }
        return {
            r: Math.round(r * 255),
            g: Math.round(g * 255),
            b: Math.round(b * 255)
        };
    }

    asColourString() {
        let c = (this.r << 16) + (this.g << 8) + this.b;
        return "#" + ('000000' + ((c) >>> 0).toString(16)).slice(-6);
    }

    asHsvDescription() {
        return `{h:${this.h}, s:${this.s}, v:${this.v}}`;
    }
}
