package main


import (
	"fmt"
	cssh "golang.org/x/crypto/ssh"
	"io"
	"time"
	"reflect"
	"log"
	"strings"
)



var sshcon ObjSSH

func handleError(err error, ip string) {
	if err != nil {
		fmt.Println(err)
		fmt.Println("connection failed: " + ip)
		for n := range network {
			if(network[n].ip == ip){
				network[n].ssh_status = "broke"
			}
		}
	}
}

func handleSessionError(err error, ip string) {
	if err != nil {
		fmt.Println(err)
		fmt.Println("connection failed: " + ip)
	}else{
		fmt.Println("connected to: " + ip)
	}
}

func readBuffForString(sshOut io.Reader) string {
	buf := make([]byte, 1000)
	n, err := sshOut.Read(buf) //this reads the ssh terminal
	waitingString := ""
	if err == nil {
		for _, v := range buf[:n] {
			fmt.Printf("%c", v)
		}
		waitingString = string(buf[:n])
	}
	for err == nil {
		// this loop will not end!!
		n, err = sshOut.Read(buf)
		waitingString += string(buf[:n])
		for _, v := range buf[:n] {
			fmt.Printf("%c", v)
		}
		if err != nil {
			fmt.Println(err)
		}

		current := string(buf[:n])
		//fmt.Println(current + "ffs")
		if(strings.Contains(current, "64")){
			fmt.Println("i found O____O")
		}
	}

	

	return waitingString
}

func write(cmd string,ip string) {
	//_, err := sshcon.ssin.Write([]byte(cmd + "\r"))
	//handleError(err)
	for n := range network {
		fmt.Println(network[n].ssh_status + ":::ready")
		fmt.Println(network[n].ip + ":::" + ip)
		if(network[n].ssh_status == "ready" && network[n].ip == ip){
			fmt.Println("sent--------------------------")
			fmt.Fprintf(network[n].ssh.ssin, "%s\n", cmd)
			
			break
		}
	}
}

func try_ssh_connection(ip string, u string, p string) {
	// create a new connection
	conn, err := cssh.Dial("tcp", ip + ":22", &cssh.ClientConfig{
		User:            u,
		Auth:            []cssh.AuthMethod{cssh.Password(p)},
		HostKeyCallback: cssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	})

	if err != nil {
		fmt.Println(err)
	}
	session, err := conn.NewSession()
	handleSessionError(err, ip)
	sshOut, Outerr := session.StdoutPipe()
	handleError(Outerr, ip)
	sshIn, Inerr := session.StdinPipe()
	handleError(Inerr, ip)
	err = session.Shell()
	handleError(err, ip)

	if(err != nil){
		
	}else{
		fmt.Println("CONNECTED----------------------------")
		tmp_Objssh := ObjSSH{
			ssin: sshIn,
		}
		for n := range network {
			if(network[n].ip == ip){
				network[n].ssh_status = "ready"
				network[n].ssh = tmp_Objssh
			}
		}
	}

	/*sshcon := &SSHCommand{
		ssin:  sshIn,
	}*/

	//_ = sshcon

	fmt.Println(reflect.TypeOf(sshIn))
	//fmt.Fprintf(sshIn, "%s\n", "ping google.com")
	//write("configure", sshIn)
	go readBuffForString(sshOut)

	

	
	//write("ping google.com", "192.168.1.240")
	//time.Sleep(5 * time.Second)

	//sshIn.Write([]byte("ping google.com" + "\r"))
	
	//session.Close()
	//conn.Close()
}


func try3(){
	quit := make(chan bool)

    /*config := &cssh.ClientConfig{
        User: "user",
        Auth: []cssh.AuthMethod{
            cssh.Password("password"),
        },
    }

    //client, err := cssh.Dial("tcp", "192.168.1.240:22", config)*/

	client, err := cssh.Dial("tcp", "192.168.1.240:22", &cssh.ClientConfig{
		User:            "pi",
		Auth:            []cssh.AuthMethod{cssh.Password("raspberry")},
		HostKeyCallback: cssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	})
    if err != nil {
        log.Fatal(err)
    }

    go func() {
        // check if connection is still alive

        // client.Wait() does not return if the network of the 
        // client is down or if the server is turned of
		fmt.Println("yay")
        err := client.Wait()
        if err != nil {
            log.Print(err)
        }else{
			fmt.Println("connected")
		}
        quit <- true
    }()

    <-quit

	fmt.Println("poo")
}