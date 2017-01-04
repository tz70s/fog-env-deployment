/*
 * Latency simulation master node app
 *
 */

var net = require('net');
var util = require('util');

const PORT = 6632;
const IP = '192.168.1.133'

var query_data = {
	"method": "monitor",
	"id": 0,
	"params":["Open_vSwitch",
		null,
		{
			"Interface" : {
				"columns": ["name", "ingress_policing_rate", "ingress_policing_burst"]
			}
		}
	]
}

var update_ingress = {
	"method": "transact",
	"params": ["Open_vSwitch", {
		"op": "update",
		"table": "Interface",
		"where": [["_uuid", "==", ["uuid","068d1e36-1009-41a1-99ba-da0083119a58"]]],
		"row": {
			"ingress_policing_rate": 0,
			"ingress_policing_burst": 0
		}
	}],
	"id": 0
}

var client = new net.Socket();

client.connect(PORT, IP, function() {
	client.write(JSON.stringify(query_data));
});

client.on('data', function(data) {
	console.log('Received');
	var object = JSON.parse(data);
	console.log(util.inspect(object, {depth: null, colors: true}));
	client.destroy();
});
