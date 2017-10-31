function checkLogin() {
	if (document.cookie.indexOf("name") >= 0) {
		document.getElementById('status').innerHTML = readCookie('name');
		document.getElementById('fb_login').remove();
	} else {
		document.getElementById('logout').hidden='true';				
	}
}
function logOut() {
	deleteCookie('name');
	deleteCookie('id');
	deleteCookie('loginToken');
	window.location('/');
}
function deleteCookie(name) {
	document.cookie = name + '=;expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}
function readCookie(name) {
	var nameEQ = name + "=";
	var ca = document.cookie.split(';');
	for(var i=0;i < ca.length;i++) {
		var c = ca[i];
		while (c.charAt(0)==' ') c = c.substring(1,c.length);
		if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
	}
	return null;
}
function vote(bookid, rating) {
	var userid = readCookie('id')
	if (undefined == userid || '' == userid) {
		alert("Must login")
		return;
	} 
	console.log("Trying to send vote to server, userid:" + userid + " bookid:" + bookid + " rating:" + rating);	
	var xhttp;
	if (window.XMLHttpRequest) {
	    xhttp = new XMLHttpRequest();
	} else {
	    // code for IE6, IE5
	    xhttp = new ActiveXObject("Microsoft.XMLHTTP");
	}
	xhttp.onreadystatechange = function() {
        // inline function to check the status
        // of our request
        // this is called on every state change
        if (xhttp.readyState === 4 &&
        xhttp.status === 200) {
            console.log("Sent vote to server successfully...");
        }
	};
 	xhttp.open("POST", "/user/vote/book/", true);
	xhttp.setRequestHeader("Content-type", "application/json");
	xhttp.send("{\"bookid\":\"" + bookid + "\", \"userid\":\"" + userid + "\",\"rating\":\"" + rating + "\"}");
}
