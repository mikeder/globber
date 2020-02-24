function publish() {
    var title = document.forms["compose-form"]["title"].value;
    var markdown = document.forms["compose-form"]["markdown"].value;
  
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
      if (this.readyState == 4 && this.status == 200) {
        // alert(this.responseText);
        // redirect?
      }
    };
    xhttp.open("POST", "/blog/entry/new", true);
    xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhttp.send("title=" + title + "&markdown=" + markdown);
  }
  