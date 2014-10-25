function sliderChanged(value) {
    var v = parseInt(value);
    v = (v * 65536) + (v * 256) + v;
    allLights(v.toString(16));
}

function allLights(colour) {
    console.log(colour);
    var req = new XMLHttpRequest();
    req.open('PUT', '/AllLights/' + colour, true);
    req.send();
    req.onreadystatechange = function() {
        if (req.readyState == 4) {
            setConnectionStatus(req.responseText)
            req = null
        }
    }
}

function animation(name) {
    var req = new XMLHttpRequest();
    req.open('PUT', '/Animation/' + name, true);
    req.send();
    req.onreadystatechange = function() {
        if (req.readyState == 4) {
            setConnectionStatus(req.responseText)
            req = null
        }
    }
}

function setConnectionStatus(status) {
    document.getElementById("status").innerHTML =  JSON.parse(status) ? "OK" : "Not Connected";
}
