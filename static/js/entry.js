function publishEntry() {
  var title = document.forms["compose-form"]["title"].value;
  // var markdown = document.forms["compose-form"]["markdown"].value;
  
  var html = $('#summernote').summernote('code');
  console.log(html)

  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      location.replace("/blog");
      location.reload();
    }
  };
  xhttp.open("POST", "/blog/compose", true);
  xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  xhttp.send("title=" + title + "&html=" + html);
}

function saveEntry(entryID) {
  var title = document.forms["compose-form"]["title"].value;
  // var markdown = document.forms["compose-form"]["markdown"].value;

  var html = $('#summernote').summernote('code');
  console.log(html)

  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      location.replace("/blog");
      location.reload();
    }
  };
  xhttp.open("POST", "/blog/compose?id="+entryID, true);
  xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  xhttp.send("title=" + title + "&html=" + html);
}
