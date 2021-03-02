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

type SSHCommand struct {
	ssin io.WriteCloser
}

var sshcon SSHCommand

func handleError(err error) {
	if err != nil {
		panic(err)
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

func write(cmd string) {
	_, err := sshcon.ssin.Write([]byte(cmd + "\r"))
	handleError(err)
}

func try_connection() {
	// create a new connection
	conn, err := cssh.Dial("tcp", "192.168.1.240:22", &cssh.ClientConfig{
		User:            "",
		Auth:            []cssh.AuthMethod{cssh.Password("")},
		HostKeyCallback: cssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	})

	if err != nil {
		fmt.Println(err)
	}
	session, err := conn.NewSession()
	handleError(err)
	sshOut, err := session.StdoutPipe()
	handleError(err)
	sshIn, err := session.StdinPipe()

	err = session.Shell()
	handleError(err)


	sshcon := &SSHCommand{
		ssin:  sshIn,
	}

	_ = sshcon

	fmt.Println(reflect.TypeOf(sshIn))
	fmt.Fprintf(sshIn, "%s\n", "ping google.com")
	//write("configure", sshIn)
	go readBuffForString(sshOut)

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