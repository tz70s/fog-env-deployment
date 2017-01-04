/*
 * Latency simulation master node app
 *
 */

var net = require('net');
var util = require('util');
var jsonfile = require('jsonfile');
const file = './node.json';

jsonfile.spaces= 4;

const PORT = 6632;
const IP = '192.168.1.133'

const query_data = {
	"method": "monitor",
	"id": 0,
	"params":["Open_vSwitch",
		null,
		{
			"Interface" : {
				"columns": ["_uuid","name", "ingress_policing_rate", "ingress_policing_burst"]
			}
		}
	]
}

var update_ingress_template = {
	"method": "transact",
	"params": ["Open_vSwitch", {
		"op": "update",
		"table": "Interface",
		"where": [["_uuid", "==", ["uuid",""]]],
		"row": {
			"ingress_policing_rate": 0,
			"ingress_policing_burst": 0
		}
	}],
	"id": 0
}
function LoadIngressFromJSON () {
	var new_update_ingress = [];
	var new_obj = jsonfile.readFileSync(file)

	var loop = function(i) {
		if (i < new_obj.length) {
			update_ingress_template.params[1]['where'][0][2][1] = new_obj[i]['_uuid'][1];
			console.log(new_obj[i]['_uuid'][1]);
			console.log(update_ingress_template.params[1]['where'][0][2][1]);
			update_ingress_template.params[1].row.ingress_policing_rate = new_obj[i].ingress_policing_rate;
			update_ingress_template.params[1].row.ingress_policing_burst = new_obj[i].ingress_policing_burst;
			new_update_ingress.push(update_ingress_template);
			loop(i+1);
		}	
	}

	loop(0);

	return new_update_ingress;
}

function UpdateFromOVSDB (action_data) {

	var client = new net.Socket();

	client.connect(PORT, IP, function() {
		client.write(JSON.stringify(action_data));
	});

	client.on('data', function(data) {
		console.log('Received ovsdb response');
		var object = JSON.parse(data);
		console.log(util.inspect(object.result, {depth: null, colors: true}));
		jsonfile.writeFile(file, ObjectHandler(object.result.Interface), function(err) {
			if(err) {
				console.error(err);
			}
		});
		client.destroy();
	});
}

function ObjectHandler (object) {
	var new_obj = [];
	Object.keys(object).forEach(function(key) {
		var val = object[key];
		new_obj.push(val.new);
	})
	return new_obj;
}

//UpdateFromOVSDB(query_data);

var anotherfile = 'metrics.json'
var last_obj = LoadIngressFromJSON();
jsonfile.writeFileSync(anotherfile, last_obj);
console.log(util.inspect(last_obj,{depth:null,colors:true}));
