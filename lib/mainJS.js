function testAjax() {
    document.getElementById("art").style.color = document.getElementById("textcolor").value
    document.getElementById("art").style.backgroundColor = document.getElementById("bgcolor").value
    var IT = document.getElementById("inputText").value
    var FL = document.getElementById("fontList").value
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var str = this.responseText;
            var str1 = str.substring(str.indexOf('id="art">'), );
            var str2 = str1.substring(0, str1.indexOf('</tex'));
            str3 = str2.replace('id="art">', "")
            document.getElementById("art").innerHTML = str3;
        }
    };
    xhttp.open("POST", "/?inputText="+IT+"&fontList="+FL, true);
    xhttp.send();
}

$(function() {
    $('#saveToFile').click(function(e) {
      var data = document.getElementById('art').value;
      var data = 'data:application/txt;charset=utf-8,' + encodeURIComponent(data);
      var el = e.currentTarget;
      el.href = data;
      el.target = '_blank';
      el.download = 'data.txt';
    });
  });