package main

import (

    "log"
	//"os/exec"
	//"strings"
	"fmt"
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
	//"reflect"
    "github.com/gorilla/websocket"
)

//sockets

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

//sockets

type obj_connection struct{
	computer_type string
	ip string
	status string
}

var lst_obj_connections [] obj_connection

func main(){

	go try_connection()

	

	//--------------------------------------------load
	//load_country("Australia")
	//go tcp_server()
	//aussie count = 63,232,000
	//--------------------------------------------load

	bewear_setup()
	//go setup_websockets()

	http.HandleFunc("/", handler)

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
        conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
        for {
            // Read message from browser
            msgType, msg, err := conn.ReadMessage()
            if err != nil {
                return
            }

			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg));

			cmd := string(msg);

			if(strings.HasPrefix(cmd,"refresh")){
				//

				if(findAll_status == "update"){
					
					for _, n := range network {
						fmt.Println(n.ip);
						/*if err = conn.WriteMessage(msgType, "pc: " + n.ip + ","); err != nil {

						}*/
						//xType := reflect.TypeOf(msg);
						//fmt.Println(xType);
						newmsg := []byte("pc: " + n.ip + ",")
						if err = conn.WriteMessage(msgType, newmsg); err != nil {
							
						}
					}

					findAll_status = "ready";
				}
			}

            // Print the message to the console
            
			

            // Write message back to browser
            if err = conn.WriteMessage(msgType, msg); err != nil {
                return
            }
        }
    })
	fmt.Println("Web sockets started");
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
		
		strRemoteAddr := con.RemoteAddr().String()
		remote_addr_clean := strings.Split(strRemoteAddr,":")[0]
		found := false
		for i := range lst_obj_connections{
			if remote_addr_clean == lst_obj_connections[i].ip && lst_obj_connections[i].status == "connected"{
				found = true
				break
			}
			log.Println(lst_obj_connections[i].ip)
		}

		if found{
			con.Close()
			log.Println("duplicate connection: " + strRemoteAddr)
		}else{
			go handleClientRequest(con)
		}
		// If you want, you can increment a counter here and inject to handleClientRequest below as client identifier
		
	}
}

func update_connection_status(ip_cleaned string){
    found_obj_connection := false
    for i, _ := range lst_obj_connections {
        if lst_obj_connections[i].ip == ip_cleaned{
            found_obj_connection = true
            lst_obj_connections[i].status = "connected"
            break
        }
    }

    if !found_obj_connection{
        new_obj_connection := obj_connection{ip: ip_cleaned,status: "connected"}
	    lst_obj_connections = append(lst_obj_connections,new_obj_connection)
    }
}

func close_connection(ip_cleaned string){
    for i, _ := range lst_obj_connections {
        if lst_obj_connections[i].ip == ip_cleaned{
            lst_obj_connections[i].status = "disconnected"
            break
        }
    }
}
 
func handleClientRequest(con net.Conn) {

	clientReader := bufio.NewReader(con)
	
	defer con.Close()
 
	strRemoteAddr := con.RemoteAddr().String()
	remote_addr_clean := strings.Split(strRemoteAddr,":")[0]
	log.Println("connected: " + strRemoteAddr)

	update_connection_status(remote_addr_clean)

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
			close_connection(remote_addr_clean)
			return
		default:
			log.Printf("error: %v\n", err)
			close_connection(remote_addr_clean)
			return
		}
 
		// Responding to the client request
		if _, err = con.Write([]byte("GOT IT!\n")); err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}
	}
}