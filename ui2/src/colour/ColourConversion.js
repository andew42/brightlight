export function HSVtoRGB(h, s, v) {
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

export function RGBtoHSV(r, g, b) {
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

// Converts #n, n or {r:x,g:y,b:z}
export function asColourObject(c) {

    // #00ff00
    if (typeof c === 'string') {
        if (c[0] === '#') {
            c = c.slice(1);
        }
        c = parseInt(c, 16)
    }

    // 200
    if (typeof c === 'number') {
        c >>>= 0;
        c = {
            r: (c & 0xFF0000) >>> 16,
            g: (c & 0xFF00) >>> 8,
            b: c & 0xFF
        };
    }

    // {r:255, g:0, b:240}
    return c;
}

export function asColourString(c) {

    c = asColourObject(c);
    c = (c.r << 16) + (c.g << 8) + c.b;
    return "#" + ('000000' + ((c) >>> 0).toString(16)).slice(-6);
}
