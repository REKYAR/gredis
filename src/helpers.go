package main

import (
	"log"
	"strings"
)

func checkErr(err error) {

	if err != nil {

		log.Fatal(err)
	}
}

func parse_req(raw_req string) []string {
	arrcom := strings.Split(raw_req, " ")
	return arrcom
}

func execute_command(command []string) {
	switch command[0] {
	case "get":
		if len(command) == 2 {
			do_get()
		}
		break
	case "set":
		if len(command) == 3 {
			do_set()
		}
		break
	case "del":
		if len(command) == 2 {
			do_del()
		}
		break
	default:

		break
	}
}

func do_get() {

}
func do_set() {

}
func do_del() {

}
