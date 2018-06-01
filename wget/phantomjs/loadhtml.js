"use strict";
var page = require('webpage').create(),
	system = require('system'),
	address;

page.settings = {
	userAgent: 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0',
	javascriptEnabled: true,
}

if (system.args.length === 1) {
	console.log('Usage: loadhtml.js <some URL>');
	phantom.exit(2);
} else {
	address = system.args[1];
	page.open(address, function (status) {
		if (status !== 'success') {
			console.log('Fail to load the address: ' + address);
			phantom.exit(1);
		} else {
			window.setTimeout(function () {
				console.log( page.content )
				phantom.exit();
			}, 2000); // Change timeout as required to allow sufficient time 
		}
	});
}

