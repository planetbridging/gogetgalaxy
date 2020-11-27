package main

import (

    "log"
	//"os/exec"
	//"strings"
	//"fmt"
	//"net"
	//"encoding/xml"
	//"encoding/xml"
	//"strconv"
	//"time"
	"net/http"
	//"sync"
)

func main(){
	//--------------------------------------------load
	load_country("Australia")

	//--------------------------------------------load

	http.HandleFunc("/", handler)
	http.Handle("/web/", http.FileServer(http.Dir(*root)))
	log.Fatal(http.ListenAndServe(":4848", nil))
}