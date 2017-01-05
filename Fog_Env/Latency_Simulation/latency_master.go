/*
Latency Simulation master node 
*/

package main

import(
	"fmt"
	"encoding/json"
	"net"
	"bytes"
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
					"_uuid",
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
			"row": map[string]int{
				"ingress_policing_rate": 0,
				"ingress_policing_burst": 0,
			},

		},
	},
	"id": 0,
}

/*
	Monitor current ovsdb-server configuration view
*/

func TopoViewOVSDB () []byte {
	conn, _ := net.Dial("tcp", IP + ":" + PORT);
	out, _ := json.Marshal(Query_data);
	fmt.Fprintf(conn, string(out));

	response := make([]byte,1024)
	conn.Read(response)

	fmt.Printf("%s\n",response);
	//res,_ := PrettyPrint(response);
	//fmt.Printf("%s\n", res);
	return response;
}

/*
	Parse the configuration json structure && Load parameters
*/

func UpdateParse(source_uuid string, ingress_rate int, ingress_burst int) []byte {
	
	tmp := Update_Ingress_Template;
	
	/* Configure uuid */

	params := tmp["params"].([]interface{});
	conf := params[1].(map[string]interface{});
	where := conf["where"].([][]interface{});
	uuid := where[0][2].([]string);
	uuid[1] = source_uuid;
	
	/* Configure ingress policing */
	row := conf["row"].(map[string]int);
	row["ingress_policing_rate"] = ingress_rate;
	row["ingress_policing_burst"] = ingress_burst;

	out, _ := json.Marshal(tmp);
	
	return out;
}

/*
	Update Interfaces configuration over Socket 
*/

func UpdateInterfaceOVSDB() {
	
	conn, err := net.Dial("tcp", IP + ":" + PORT);
	if err != nil {
		fmt.Printf("Can't access ovsdb-server, please check out network configuration");
		return;
	}

	data := UpdateParse("3d0f0cb0-2752-4dde-a6bc-cb6bf07b1785",0,0);

	fmt.Fprintf(conn, string(data))
	response := make([]byte, 1024)
	conn.Read(response)
	fmt.Printf("%s\n", response)
	TopoViewOVSDB();
}


func PrettyPrint(origin []byte) ([]byte, error) {

	var out bytes.Buffer;
	err := json.Indent(&out, []byte(origin), "", "\t");
	fmt.Printf("%s\n", out)
	return out.Bytes(), err;
}

func main() {
	UpdateInterfaceOVSDB();
}
