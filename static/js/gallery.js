var numRows = 3
var numCols = 4

var divs=document.getElementsByTagName('div');

for(var i = 0, len = divs.length;i<len;i++){
    var newDiv = document.createElement('div');
    divs[i].appendChild(newDiv);
}