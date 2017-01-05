/*
Latency Simulation master node 
*/

package main

import(
	"fmt"
	"encoding/json"
	"net"
	"net/http"
	//"bytes"
	"io/ioutil"
	"os"
	"strconv"
)

const PORT string = "6632";
const IP string = "127.0.0.1";

var Query_data = map[string]interface{}{
	"method" : "monitor",
	"id" : 0,
	"params" : []interface{} {
		"Open_vSwitch",
		0,
		map[string]interface{} {
			"Interface": map[string]interface{} {
				"columns": []string {
					"name",
					"ingress_policing_rate",
					"ingress_policing_burst",
				},
			},
		},	
	},
}

var Update_Ingress_Template = map[string]interface{} {
	"method": "transact",
	"params": []interface{} {
		"Open_vSwitch",
		map[string]interface{} {
			"op": "update",
			"table": "Interface",
			"where" : [][]interface{}{
				{
					"_uuid",
					"==",
					[]string{
						"uuid",
						"",
					},
				},
			},
			"row": map[string]float64{
				"ingress_policing_rate": 0,
				"ingress_policing_burst": 0,
			},

		},
	},
	"id": 0,
}

type Interface_Conf struct {
	Uuid string
	Name string
	Ingress_policing_rate float64
	Ingress_policing_burst float64
	Netem_latency float64	
}

func TCRequest (action string, port_name string, delay float64) float64 {

	delay_time := strconv.FormatFloat(delay, 'f', -1, 64)
	
	var out []byte;
	
	resp, err := http.Get("http://" + IP + ":8989" + "/?action="+action+"&port_name="+port_name+"&delay_time="+delay_time);

	if err != nil {

		fmt.Printf("Request error : %v \n", err);

	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Printf("GET OK: %s\n", body);
		out = body;
	}
	
	flo,_ := strconv.ParseFloat(string(out), 64);
	return flo;
}

/*
	Load node_conf.json in structure
*/

func LoadJSONConf () {
	file, err := ioutil.ReadFile("./node_conf.json");
	if err != nil {
		fmt.Printf("Read json file error : %v \n", err);
		os.Exit(1);
	}
	var collects [] Interface_Conf;
	json.Unmarshal(file,&collects);
	
	//out,_ := json.Marshal(collects);
	//fmt.Println(string(out));
	
	var i int = 0;
	for i < len(collects) {
		UpdateIngress(collects[i]);
		i++;
	}
}

/*
	Update OVSDB Ingree policy
*/

func UpdateIngress (obj Interface_Conf) {
	UpdateInterfaceOVSDB(obj.Uuid, obj.Ingress_policing_rate, obj.Ingress_policing_burst);
	if (obj.Netem_latency == -1) {
		TCRequest("del", obj.Name, obj.Netem_latency);
	} else {
		TCRequest("add", obj.Name, obj.Netem_latency);
	}
}

/*
	INIT
*/

func Init() {
	SaveJSONConf(TopoViewOVSDB());
}

/*
	Save structures to node_conf.json
*/

func SaveJSONConf(out []byte) {
	err := ioutil.WriteFile("./node_conf.json", out, 0644);
	if err != nil {
		fmt.Printf("Can't successfully save json file %v\n", err);
	}
}

/*
	Monitor current ovsdb-server configuration view
*/

func TopoViewOVSDB () []byte {
	conn, _ := net.Dial("tcp", IP + ":" + PORT);
	out, _ := json.Marshal(Query_data);
	fmt.Fprintf(conn, string(out));

	response := make([]byte,1024)
	resplen, err := conn.Read(response);
	if err != nil {
		fmt.Printf("Monitoring failed : %v \n",err);
	}

	//fmt.Printf("%s\n",response[:resplen]);
	//res,_ := PrettyPrint(response[:resplen]);
	//fmt.Printf("%s\n", res);
	
	var rslice [] Interface_Conf = ParseMonitoring(response[:resplen]);
	disp, _ := json.MarshalIndent(rslice, "", "  ");
	fmt.Printf("%s\n",disp);
	return disp;
}

/*
	Parse Monitor response's json format
*/

func ParseMonitoring(input []byte) [] Interface_Conf {
	var rslice [] Interface_Conf;
	var parseObject interface{};
	err := json.Unmarshal(input, &parseObject);
	if err != nil {
		fmt.Printf("Parse Monitoring Object failed : %v\n", err);
		os.Exit(1);
	}
	result := parseObject.(map[string]interface{})
	interf := result["result"].(map[string]interface{})
	node := interf["Interface"].(map[string]interface{})
	for k, v := range node {
		nest := v.(map[string]interface{})
		last := nest["new"].(map[string]interface{})

		var tmp_node Interface_Conf;
		tmp_node.Uuid = k;
		tmp_node.Name = last["name"].(string);
		tmp_node.Ingress_policing_rate = last["ingress_policing_rate"].(float64);
		tmp_node.Ingress_policing_burst = last["ingress_policing_burst"].(float64);
		tmp_node.Netem_latency = TCRequest("get", tmp_node.Name, 0);
		rslice = append(rslice,tmp_node);
	}
	return rslice;
}

/*
	Parse the configuration json structure && Load parameters
*/

func UpdateParse(source_uuid string, ingress_rate float64, ingress_burst float64) []byte {
	
	tmp := Update_Ingress_Template;
	
	/* Configure uuid */

	params := tmp["params"].([]interface{});
	conf := params[1].(map[string]interface{});
	where := conf["where"].([][]interface{});
	uuid := where[0][2].([]string);
	uuid[1] = source_uuid;
	
	/* Configure ingress policing */
	row := conf["row"].(map[string]float64);
	row["ingress_policing_rate"] = ingress_rate;
	row["ingress_policing_burst"] = ingress_burst;

	out, _ := json.Marshal(tmp);
	
	return out;
}

/*
	Update Interfaces configuration over Socket 
*/

func UpdateInterfaceOVSDB(source_uuid string, ingress_rate float64, ingress_burst float64) {
	
	conn, err := net.Dial("tcp", IP + ":" + PORT);
	if err != nil {
		fmt.Printf("Can't access ovsdb-server, please check out network configuration");
		return;
	}

	data := UpdateParse(source_uuid, ingress_rate, ingress_burst);

	fmt.Fprintf(conn, string(data))
	response := make([]byte, 1024)
	conn.Read(response)
}

/*

func PrettyPrint(origin []byte) ([]byte, error) {

	var out bytes.Buffer;
	err := json.Indent(&out, []byte(origin), "", "  ");
	return out.Bytes(), err;
}

*/

func main() {
	if (len(os.Args) != 2) {
		fmt.Printf("USAGE:\n\tinit\t:\tInitialized Latencies Configuration & Rerwite node_conf.json\n\tupdate\t:\tUpdate Modification from node_conf.json\n\tview\t:\tCurrent Latencies View\n");
		os.Exit(0);
	}
	switch(os.Args[1]) {
		case "init":
			Init();
		case "update":
			LoadJSONConf();
			TopoViewOVSDB();
		case "view":
			TopoViewOVSDB();
		default:
			fmt.Printf("USAGE:\n\tinit\t:\tInitialized Latencies Configuration & Rerwite node_conf.json\n\tupdate\t:\tUpdate Modification from node_conf.json\n\tview\t:\tCurrent Latencies View\n")
			os.Exit(0);
	}

}
