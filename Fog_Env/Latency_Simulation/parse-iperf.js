var iperf = require('./iperf.json');

var bandwidth = "";

for (var i = 0; i < iperf.intervals.length; i++ ) {
	var num = iperf.intervals[i].sum.bits_per_second / 1024;
	num = Number((num).toFixed(2));
	bandwidth += num + "\n";
}

console.log(bandwidth);
var fs = require('fs');

fs.writeFile('iperf.dat',bandwidth,function(err) {
	if(err)
		return console.log(err);
	console.log('finish');
});
