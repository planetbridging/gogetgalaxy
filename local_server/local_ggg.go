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
	"bufio"
	"io"
	"net"
	"strings"
)


func main(){
	//--------------------------------------------load
	//load_country("Australia")
	go tcp_server()
	//--------------------------------------------load

	http.HandleFunc("/", handler)
	http.Handle("/web/", http.FileServer(http.Dir(*root)))
	log.Println("Serving at localhost:4848...")
	log.Fatal(http.ListenAndServe(":4848", nil))
}

func tcp_server() {
	listener, err := net.Listen("tcp", "0.0.0.0:4849")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
 
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
 
		// If you want, you can increment a counter here and inject to handleClientRequest below as client identifier
		go handleClientRequest(con)
	}
}
 
func handleClientRequest(con net.Conn) {
	defer con.Close()
 
	clientReader := bufio.NewReader(con)
 
	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')
 
		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			} else {
				log.Println(clientRequest)
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
 
		// Responding to the client request
		if _, err = con.Write([]byte("GOT IT!\n")); err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}
	}
}