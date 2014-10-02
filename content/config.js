function config() {
    var stripIndex = document.getElementById('si').value
    var stripLength = document.getElementById('sl').value
    var req = new XMLHttpRequest();
    req.open('PUT', '/Config/' + stripIndex + ',' + stripLength, true);
    req.send();
}
