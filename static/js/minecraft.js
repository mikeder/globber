var getStatus = function () {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var status = JSON.parse(this.responseText);
            updateStatus(status.result);
        }
    };
    xhttp.open("GET", "/minecraft/status", true);
    xhttp.send();
}

var periodicUpdate = setInterval(getStatus, 5000);

var updateStatus = function (server) {
    var banner = document.getElementById("mcstatus");

    if (server.online) {
        banner.classList.remove("alert-danger");
        banner.classList.remove("alert-warning");
        banner.classList.add("alert-success");
        banner.innerText = "Server Status: ONLINE"
    } else {
        banner.classList.remove("alert-success");
        banner.classList.remove("alert-warning");
        banner.classList.add("alert-danger");
        banner.innerText = "Server Status: OFFLINE"
    }

    document.getElementById("servername").innerHTML = server.motd;
    document.getElementById("address").innerHTML = "<b>Address:</b> " + server.address;
    document.getElementById("players").innerHTML = "<b>Players:</b> " + server.current_players + "/" + server.max_players;
    document.getElementById("latency").innerHTML = "<b>Latency:</b> " + server.latency;
    document.getElementById("version").innerHTML = "<b>Version:</b> " + server.version;

    updatePlayers(server.online_players)
}

var updatePlayers = function (players) {
    var player_list = document.getElementById("players_list")
    player_list.innerHTML = ""

    players.forEach(player => {
        var entry = document.createElement('li');
        entry.appendChild(document.createTextNode(player.name));
        player_list.appendChild(entry);
    });
}