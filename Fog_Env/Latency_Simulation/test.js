var jsonfile = require('jsonfile');
var file = './node.json';

jsonfile.spaces = 4;

jsonfile.readFile(file, function(err, obj) {
	console.log(JSON.stringify(obj));
})

var object = {
	name : "jon"
}

jsonfile.writeFile(file,object, function(err) {
	if (err) {
		console.error(err);
	}
})
