function checkLogin() {
	document.getElementById('search').addEventListener("keyup", function(event){doSearch(event)});
	if (document.cookie.indexOf("name") >= 0) {
		document.getElementById('status').innerHTML = readCookie('name');
		document.getElementById('fb_login').remove();
	} else {
		document.getElementById('logout').hidden='true';
		document.getElementById('ratings').hidden='true';
		return;				
	}

	// Fetch book vote
	if (document.getElementById("bookId") != null) {
		var bookid = document.getElementById("bookId").innerText;
		console.log("Trying to fetch book vote from server...");	
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
			if (xhttp.readyState === 4 && xhttp.status === 200) {
				var data = JSON.parse(xhttp.responseText);
				var clickVote = document.getElementById("clickVote");
				clickVote.onclick="";
				clickVote.innerText="You have voted:" + data.Rating;
			} else if (xhttp.readyState === 4 && xhttp.status === 400) {
				alert("Server replied bad request...");	
			}
		};
		xhttp.open("GET", "/user/vote/book/" + bookid, true);
		xhttp.send();
	}
}
function logOut() {
	deleteCookie('name');
	deleteCookie('id');
	deleteCookie('loginToken');
	window.location.reload(false);
}
function deleteCookie(name) {
	document.cookie = name + '=;expires=Thu, 01 Jan 1970 00:00:01 GMT;path=/;';
}
function readCookie(name) {
	var nameEQ = name + "=";
	var ca = document.cookie.split(';');
	for (var i=0;i < ca.length;i++) {
		var c = ca[i];
		while (c.charAt(0)==' ') c = c.substring(1,c.length);
		if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
	}
	return null;
}
function vote(bookid, rating) {
	var userid = readCookie('id')
	if (undefined == userid || '' == userid) {
		window.location.href = "/login";
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
        xhttp.status === 201) {
            document.getElementById("clickVote").innerText = "You have voted:" + rating;
        } else if (xhttp.readyState === 4 && xhttp.status === 400) {
			alert("Server replied bad request...");	
		}
	};
 	xhttp.open("POST", "/user/vote/book/", true);
	xhttp.setRequestHeader("Content-type", "application/json");
	xhttp.send("{\"bookid\":\"" + bookid + "\", \"userid\":\"" + userid + "\",\"rating\":" + rating + "}");
}
var resultDiv = document.getElementById("searchresults");
function doSearch(e) {
	var keynum;
	if (window.event) { // IE                    
		keynum = e.keyCode;
	} else if (e.which) { // Netscape/Firefox/Opera                   
		keynum = e.which;
	}
	//alert(keynum === 13);

	var query = document.getElementById("search").value;
	if (query === '') {
		clearSearchResults();		
		return;
	}
	console.log("Trying to send search to server...");	
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
        if (xhttp.readyState === 4 && xhttp.status === 200) {
			var data = JSON.parse(xhttp.responseText);
			clearSearchResults();
			for(j=0; j<data.length; j++){
				var result = document.createElement("article");
				result.innerText = data[j].Title;
				resultDiv.appendChild(result);		
			}	
		} else if (xhttp.readyState === 4 && xhttp.status === 400) {
			alert("Server replied bad request...");	
		}
	};
 	xhttp.open("GET", "/search?q=" + query, true);
	xhttp.setRequestHeader("Content-type", "application/json");
	xhttp.send();
}

function clearSearchResults() {
	while (resultDiv.firstChild) {
		resultDiv.removeChild(resultDiv.firstChild);
	}
}
