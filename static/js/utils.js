"use strict";

var attend = function (el, day, month, year) {
	var panel = document.getElementById('attend-data');
	var bounds = el.getBoundingClientRect();
	panel.style.display = 'block';
}

var closeAttend = function () {
	document.getElementById('attend-data').style.display = 'none';	
}
