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
	/*
	for k, v := range r.Form {
		fmt.Printf("key: %s, value: %s\n", k,v);
	}
	*/
	fmt.Fprintf(w, "Linux tc delay configuration accessed");
	mod_delay(r.Form["action"][0], r.Form["port_name"][0], r.Form["delay_time"][0]);
}

func main () {
	http.HandleFunc("/", pathHandler)
	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}
}

func mod_delay (action string, port_name string, delay_time string) {
	var check_bit int = check_exist(port_name);
	if check_bit == 0 && action == "add"{
		action = "change";
	}

	if action == "del" {
		if check_bit == 1 {
			fmt.Println("Target port not exist")
			return;
		} else if check_bit == 0 {
			script_split := []string{"tc","qdisc","del","dev",port_name,"root"}
			script := strings.Join(script_split, " ");
			_, err := exec.Command("sh", "-c", script).Output();
			if err != nil {
				fmt.Printf("%s\n",err);
			}
			return;
		}
	}

	script_split := []string{"tc","qdisc",action ,"dev",port_name, "root netem delay", delay_time};
	script := strings.Join(script_split, " ");
	
	_, err := exec.Command("sh", "-c", script).Output();

	if err != nil {
		fmt.Printf("%s\n",err)
	}
}

func check_exist (port_name string) int {

	script_split := []string{"tc","qdisc","|","grep",port_name};
	
	script := strings.Join(script_split, " ");

	_, err := exec.Command("sh", "-c", script).Output()
	
	if err != nil {
		return 1;
	} else {
		return 0;
	}
}
