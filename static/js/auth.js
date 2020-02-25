var login = function() {
  var email = document.forms["login-form"]["email"].value;
  var password = document.forms["login-form"]["password"].value;

  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      location.reload();
    }
  };
  xhttp.open("POST", "/auth/login", true);
  xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  xhttp.send("email=" + email + "&password=" + password);
}

var removeJWT = function() {
  document.cookie = 'jwt=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/;';
  location.replace("/")
  location.reload();
}
