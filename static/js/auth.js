var login = function () {
  var email = document.forms["login-form"]["email"].value;
  var password = document.forms["login-form"]["password"].value;

  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function () {
    if (this.readyState == 4 && this.status == 200) {
      location.reload();
    }
  };
  xhttp.open("POST", "/auth/login", true);
  xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  xhttp.send("email=" + email + "&password=" + password);
}

var refresh = function () {
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function () {
    if (this.readyState == 4 && this.status == 200) {
      console.log("refresh done")
    }
  };
  xhttp.open("POST", "/auth/refresh", true);
  xhttp.send();
}

var removeJWT = function () {
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function () {
    if (this.readyState == 4 && this.status == 200) {
      console.log("logout done")
    } else {
      document.cookie = 'jwt=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/;';
      document.cookie = 'jwt_refresh=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/;';
    }
  };
  xhttp.open("POST", "/auth/logout", true);
  xhttp.send();
  
  location.replace("/blog");
}

function getCookie(cname) {
  var name = cname + "=";
  var decodedCookie = decodeURIComponent(document.cookie);
  var ca = decodedCookie.split(';');
  for (var i = 0; i < ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

setInterval(refresh, 30 * 1000);
