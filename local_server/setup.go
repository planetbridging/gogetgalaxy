package main

import(
	"fmt"
	"runtime"
	"os"
    "os/user"
	"bufio"
	"strings"
)

var current_os string
var current_ssh_path string

func setup_os(){

	user, err := user.Current()
    if err != nil {
        panic(err)
    }

	if runtime.GOOS == "windows" {
		current_os = "windows"
		current_ssh_path =  user.HomeDir + "\\onedrive\\ggg\\ssh.txt"

		if _, err := os.Stat(user.HomeDir + "\\onedrive\\ggg\\"); os.IsNotExist(err) {
			e := os.Mkdir(user.HomeDir + "\\onedrive\\ggg\\", 0755)
			if e != nil {
				fmt.Println("failed to create folder")
			}else{
				fmt.Println("created")
			}
		}

		if _, err := os.Stat(current_ssh_path); !os.IsNotExist(err) {
			import_ssh_connections()
		}
	}

	

	

	/*

    // Current User
    fmt.Println("Hi " + user.Name + " (id: " + user.Uid + ")")
    fmt.Println("Username: " + user.Username)
    fmt.Println("Home Dir: " + user.HomeDir)

    // Get "Real" User under sudo.
    // More Info: https://stackoverflow.com/q/29733575/402585
    fmt.Println("Real User: " + os.Getenv("SUDO_USER"))*/
}

func import_ssh_connections(){
	file, err := os.Open(current_ssh_path) 
  
    if err != nil { 
        fmt.Println("failed to open ssh import") 
    }else{
		scanner := bufio.NewScanner(file) 

		scanner.Split(bufio.ScanLines) 
		var text []string 
	
		for scanner.Scan() { 
			text = append(text, scanner.Text()) 
		} 

		file.Close() 

		for _, each_ln := range text { 
			if(strings.Contains(each_ln, ",")){
				l := strings.Split(each_ln, ",")
				if(len(l) == 3){
					found := false
					
					for _, n := range network {
						if n.ip == l[0]{
							found = true
							break
						}
					}
					if !found{
						tmp_ObjPC := ObjPC{
							ip: l[0],
							ssh_status: "readytotry",
							ssh_username: l[1],
							ssh_password: l[2],
						}
						network = append(network, tmp_ObjPC)
					}
				}
			}
		} 
	}
    
}

func load_static_ssh(){
	for _, n := range network{
		go try_ssh_connection(n.ip,n.ssh_username,n.ssh_password)
	}
}