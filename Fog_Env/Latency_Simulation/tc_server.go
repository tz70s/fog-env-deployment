package main

import (
	"fmt"
	"os/exec"
	"strings"
	"net/http"
	"log"
)
func pathHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm();
	fmt.Println(r.Form);
	var delay_time string = "";
	delay_time = mod_delay(r.Form["action"][0], r.Form["port_name"][0], r.Form["delay_time"][0]);
	fmt.Fprintf(w, delay_time);
}

func main () {
	http.HandleFunc("/", pathHandler)
	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}
}

func mod_delay (action string, port_name string, delay_time string) string {

	var check_bit int = check_exist(port_name);
	
	var delay_tmp string = "-1";
	
	var script string;
	
	switch(action) {
		case "get":
			if check_bit == 1 {
				fmt.Println("Target port not exist");
				return "-1";
			}
			//fmt.Println("inter");
			script_split := []string {"tc","qdisc","|","grep",port_name};
			script = strings.Join(script_split," ");

			out, err := exec.Command("sh", "-c", script).Output()
			if err != nil {
				fmt.Printf("Command execution failed : %v\n", err);
			}
			outt := strings.Split(string(out),"\n");
			var i = 0;
			var new_out string = string(out);
			for i < len(outt) {
				
				if strings.Contains(outt[i],"netem") {
					new_out = outt[i];
					//fmt.Println(new_out);
					break;
				}
				i++;
				
			}
			
			sub_strings := strings.Split(new_out," ");

			delay_tmp = sub_strings[len(sub_strings)-1];
			//fmt.Println(delay_tmp);
			return delay_tmp[:len(delay_tmp)-2];

		case "add":
			if check_bit == 0 {
				fmt.Println("Target port existed");
				action = "change";
			}
			script_split := []string{"tc","qdisc",action,"dev", port_name, "root netem delay", delay_time};
			script = strings.Join(script_split, " ");
			script = script + "ms";

			delay_tmp = delay_time;
			
			_, err := exec.Command("sh", "-c", script).Output();
			
			if err != nil {
				fmt.Printf("Execute command failed : %v\n", err);
			}
			
			break;

		case "del":
			if check_bit == 1 {
				fmt.Println("Target port not exist");
				return "-1";
			}

			script_split := []string{"tc","qdisc","del","dev",port_name,"root"}
			script = strings.Join(script_split, " ");
			delay_tmp = "-1";
			_, err := exec.Command("sh", "-c", script).Output();
			if err != nil {
				fmt.Printf("Execute command failed : %v\n", err);
			}
			break;

		default:
			return "-2";
	}

	return delay_tmp;
}

func check_exist (port_name string) int {

	script_split := []string{"tc","qdisc","|","grep",port_name};
	
	script := strings.Join(script_split, " ");

	out, _ := exec.Command("sh", "-c", script).Output()
	exist := strings.Contains(string(out), "netem");
	if exist == true {
		fmt.Println("exist");
		return 0
	} else {
		fmt.Println("not exist")
		return 1
	}
}
