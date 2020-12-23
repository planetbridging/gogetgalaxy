package main

import (
	//test routes max
    "flag"
    //"fmt"
    "os"
    "runtime"
    "time"
    //"math"
    //"os/user"
    
	//test nmap max with cpu
	"io/ioutil"
    "strconv"
	"strings"
	
	//cmd
	"log"
	"os/exec"
    "sync"
    "context"
    "regexp"
    

    //tcp client
    "bufio"
	"io"
    "net"
)

type obj_connection struct{
	computer_type string
    ip string
    status string
}

type result struct {
    data string
    err string // you must use error type here
}

type Response struct {
    data   interface{}
    status bool
  }


var lst_obj_connections [] obj_connection

var os_type = ""
var found_server = "192.168.1.2"

//static nmap cmds
//ping
var icmp_pe = "nmap -PE -sn -oG - "
var icmp_pp = "nmap -PP -sn -oG - "
var icmp_pm = "nmap -PM -sn -oG - "
//port scan
var ports_ss = "nmap -Pn -sS -oG - --min-rate 10000 --top-ports 1000 "
var ports_st = "nmap -Pn -sT -oG - --min-rate 10000 --top-ports 1000 "
var ports_sa = "nmap -Pn -sA -oG - --min-rate 10000 --top-ports 1000 "
//windows get cpu usage
var windwows_cpu = "wmic cpu get loadpercentage"
var max_nmap_scans = 0

//furious -w 65535 -s connect --ports 1-65535 declair.in

func main(){
    /*scan_ip := "192.168.1.1"
    found := ping(scan_ip)
    if found{
       port_scan(scan_ip)
    }*/
    
   //log.Println(cmd(ports_sa+scan_ip))
   //tcp_client()
   
   
   
   //cpu_load_test()
    //multi_cmd_wait
    /*var lst_tmp_cmds [] string
   for n := 1; n < 5; n++{
       //go cmd("furious -s -w 65535 connect --ports 1-65535 declair.in")
       lst_tmp_cmds = append(lst_tmp_cmds,"furious -s connect -w 65535 --ports 1-65535 declair.in")
   }

   multi_tmp_results := multi_cmd_wait(lst_tmp_cmds)

   for _,i := range multi_tmp_results{
       log.Println(i)
   }*/


    //tmp_cmd := cmd("furious -s connect -w 65535 --ports 1-65535 192.168.1.1")

    //get_furious_ports(tmp_cmd)
   //time.Sleep(20 * time.Second)
   //log.Println("yay")
   //log.Println(get_cpu())
   //search_for_servers()
   //time.Sleep(5 * time.Second)
   
   //log.Println("Starting cpu nmap load test")
   //
   setup_ggg_profile()
}

func setup_ggg_profile(){
    ggg_profile_location := get_home_dir() + "/ggg/"


    if _, err := os.Stat(ggg_profile_location + "ggg_settings.txt"); os.IsNotExist(err) {
        cpu_load_test()
    }

    if _, err := os.Stat(ggg_profile_location); !os.IsNotExist(err) {
        read_ggg_profile(ggg_profile_location + "ggg_settings.txt")
    }else{
        err := os.Mkdir(ggg_profile_location, 0755)
        if err != nil {
            //log.Fatal(err)
            log.Println("failed to create profile folder")
        }else{
            create_new_profile(ggg_profile_location)
        }
    }

    if _, err := os.Stat(ggg_profile_location + "ggg_settings.txt"); !os.IsNotExist(err) {
    }else{
        create_new_profile(ggg_profile_location)
    }
}

func create_new_profile(home_path string){
    f, err := os.Create(home_path + "ggg_settings.txt")

    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()
    tmp_nmap_max := strconv.Itoa(max_nmap_scans)
    words := []string{"max_nmap: " + tmp_nmap_max}

    for _, word := range words {

        _, err := f.WriteString(word + "\n")

        if err != nil {
            log.Fatal(err)
            log.Println("failed to create ggg local profile")
        }
    }

    log.Println("ggg profile created at: " + home_path)
}

func read_ggg_profile(home_ggg string){
    file, err := os.Open(home_ggg)
    if err != nil {
        //log.Fatal(err)
        log.Println("failed to import ggg profile")
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if strings.Contains(scanner.Text(),"max_nmap:"){
            tmp_nmap_max := strings.Split(scanner.Text(),"max_nmap:")[1]
            i, _ := strconv.Atoi(tmp_nmap_max)
            max_nmap_scans = i
            log.Println("max nmap scans: " + tmp_nmap_max)
        }
    }

    if err := scanner.Err(); err != nil {
        //log.Fatal(err)
        log.Println("failed to import ggg profile")
    }
}



func get_home_dir() string {
    if runtime.GOOS == "windows" {
        home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")  + "\\"
        if home == "" {
            home = os.Getenv("USERPROFILE")
        }
        return home
    } else if runtime.GOOS == "linux" || runtime.GOOS == "android"{
        home := "/home/" + os.Getenv("SUDO_USER") + "/"
        if home != "" {
            return home
        }
    }
    return os.Getenv("HOME")
}

func get_furious_ports(tmp string){
    scanner := bufio.NewScanner(strings.NewReader(tmp))
    for scanner.Scan() {
        
        if strings.Contains(scanner.Text(), "/tcp"){
            s := strings.Split(scanner.Text(), "/tcp")
            reg, err := regexp.Compile("[^0-9]+")
            if err != nil {
                log.Fatal(err)
            }
            tmp_port := reg.ReplaceAllString(s[0], "")
            log.Println(tmp_port)
        }
    }
}

func cpu_load_test(){
    cpu_load := ""
    if runtime.GOOS == "windows" {
        os_type = "windows"
        cpu_load = "get_cpu_usage.bat"
    }else if runtime.GOOS == "linux"{
        os_type = "linux"
        log.Println(os_type)
    }

    

    max_cpu_limit := 70
    tested_cpu_limit := 0
    nmap_max_scans := 0
    break_loop := false

    for n := 1; n < 20; n++{
        try_amount := n * 20
        
        tmp_scan_ip := found_server

        if tmp_scan_ip == ""{
            tmp_scan_ip = "127.0.0.1"
        }
        //go test_time_out(try_amount)
        for c := 1; c <= try_amount; c++{
            go cmd("nmap -sT -p- " + tmp_scan_ip)
        }
        
    
        for i := 0; i < 5; i++{
            time.Sleep(2 * time.Second)

            tmp_amount := 0

            if os_type == "windows"{
                cpu_finished := cmd(cpu_load)

                //log.Println(cpu_finished)
                //clean := strings.Replace(cpu_finished, "LoadPercentage", "", -1)
                //cpu_amount, _ := strconv.Atoi(clean)
                
                //log.Println(cpu_amount)
                //log.Println("yay"+cpu_finished + "wtf")
                scanner := bufio.NewScanner(strings.NewReader(cpu_finished))
                for scanner.Scan() {
        
                    reg, err := regexp.Compile("[^0-9]+")
                    if err != nil {
                        log.Fatal(err)
                    }
                    processedString := reg.ReplaceAllString(scanner.Text(), "")
        
                    
                    if processedString != ""{
                        log.Println(processedString+" used with ",try_amount)
                        tmp_amount,_ = strconv.Atoi(processedString)
                        break
                    }
                }
            }else if os_type =="linux"{
                tmp_amount = get_cpu()
            }
            

            
            if tmp_amount >= max_cpu_limit || tested_cpu_limit >= max_cpu_limit{
                break_loop = true
                break
            }else if tmp_amount > tested_cpu_limit{
                tested_cpu_limit = tmp_amount
                nmap_max_scans = try_amount
            }
        } 
    
        time.Sleep(5 * time.Second)
        if os_type == "windows"{
            log.Println("killing all nmap tasks")
            //keep crashing
            //cmd("taskkill /im nmap.exe /t /f")
            go cmd("kill_all_nmap.bat")
            time.Sleep(5 * time.Second)
        }else if os_type == "linux"{
            log.Println("killing all nmap tasks")
            //keep crashing
            //cmd("taskkill /im nmap.exe /t /f")
            go cmd("killall nmap")
            time.Sleep(5 * time.Second)
        }
        if break_loop{
            break
        }
        log.Println("nmap load test limit: ",try_amount)
    }
   
    log.Println("Cpu nmap limit completed: ",nmap_max_scans)
    max_nmap_scans = nmap_max_scans
}
//----------------------------------------cpu usage testing
func test_time_out(test_amount int){
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    ch := make(chan Response, 1)

    go func() {

        var lst_test_cmds [] string

        for i := 0; i < test_amount; i++{
            lst_test_cmds = append(lst_test_cmds,"nmap -sT -p- 127.0.0.1")
        }
        multi_cmd_wait(lst_test_cmds)
        //time.Sleep(10 * time.Second)

        select {
        default:
            ch <- Response{data: "data", status: true}
        case <-ctx.Done():
            log.Println("Canceled by timeout")
            return
        }
    }()

    select {
    case <-ch:
        log.Println("Read from ch")
    case <-time.After(20 * time.Second):
        log.Println("Timed out")
    }
}
//----------------------------------------tcp connection

func search_for_servers(){
    for{
        if len(lst_obj_connections) == 0{
            find_local_server()
        }else{
            found_connection := false
            for i, _ := range lst_obj_connections {
                if lst_obj_connections[i].status == "connected"{
                    found_connection = true
                    found_server = lst_obj_connections[i].ip
                }
            }
            if !found_connection{
                find_local_server()
            }
        }
        time.Sleep(10 * time.Second)
    }
}

func get_local_ip()string{
    ifaces, _ := net.Interfaces()
    // handle err
    for _, i := range ifaces {
        addrs, _ := i.Addrs()
        // handle err
        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                    ip = v.IP
            case *net.IPAddr:
                    ip = v.IP
            }
            // process IP address
           
            if ip.String() != "127.0.0.1" && strings.Contains(ip.String(),"."){
                //log.Println(ip)
                return ip.String()
            }
        }
    }
    return ""
}
func find_local_server(){
    local_ip := get_local_ip()
    ip_trimed := strings.Split(local_ip,".")
    ip_joined := ip_trimed[0] + "." + ip_trimed[1] + "." + ip_trimed[2] + "."
    for i := 1; i <= 255; i++ {
        //log.Println(ip_joined + strconv.Itoa(i))
        tmp_ip := ip_joined + strconv.Itoa(i)
        //log.Println(tmp_ip)
        go try_connection(tmp_ip)
        
    }
}

func try_connection(tmp_ip string){
    con, err := net.Dial("tcp", tmp_ip+":4849")
    if err != nil {
        //log.Println("Failed: " + tmp_ip)
    }else{
        log.Println("Found: " + tmp_ip)
        tcp_client(con,tmp_ip)
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

func tcp_client(con net.Conn,ip_cleaned string) {
	//con, err := net.Dial("tcp", "192.168.1.201:4849")
	/*if err != nil {
        log.Println(err)
        return
    }*/

    
    update_connection_status(ip_cleaned)
    

	defer con.Close()
 
	clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(con)
 
	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')
 
		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
				log.Printf("failed to send the client request: %v\n", err)
			}
		case io.EOF:
			log.Println("client closed the connection")
			return
		default:
			log.Printf("client error: %v\n", err)
			return
		}
 
		// Waiting for the server response
		serverResponse, err := serverReader.ReadString('\n')
 
		switch err {
		case nil:
			log.Println(strings.TrimSpace(serverResponse))
		case io.EOF:
            log.Println("server closed the connection")
            close_connection(ip_cleaned)
			return
		default:
            log.Printf("server error: %v\n", err)
            close_connection(ip_cleaned)
			return
		}
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




//----------------------------------------cmd
func ping(ip string)bool{

	pe := icmp_pe + ip
	pp := icmp_pp + ip
	pm := icmp_pm + ip
	
	

	lst_tmp_ping:= []string{pe,pp,pm}
	log.Println(lst_tmp_ping)

    results := multi_cmd_wait(lst_tmp_ping)
	for _, r := range results{
        if strings.Contains(r,"Status: Up"){
            return true
        }
        //log.Println(r)
    }
    return false
}

func port_scan(ip string){
    ss := ports_ss + ip
	st := ports_st + ip
	sa := ports_sa + ip
    

	lst_tmp_port:= []string{ss,st,sa}
	//log.Println(lst_tmp_port)

    results := multi_cmd_wait(lst_tmp_port)
	for _, r := range results{
        log.Println(r)
    }
}

func multi_cmd_wait(lstcmds []string)[]string{
	

    var lst_results [] result
    var wg sync.WaitGroup
    ch := make(chan result)
	
	
    //log.Println(<-nums) // Read the value from unbuffered channel
    
    for _, c := range lstcmds {
        //channel := make(chan string)
        //lst_channels = append(lst_channels, channel)
        wg.Add(1)
        go cmd_wait(c,&wg, ch)
    }

    

    go func() {
        for v := range ch {
            lst_results = append(lst_results, v)
        }
    }()

    wg.Wait()
    close(ch)
	
    var lst_data [] string
    for d := range lst_results {
        lst_data = append(lst_data,lst_results[d].data )
    }

	return lst_data
	
}

func cmd_wait(c string,wg *sync.WaitGroup, results chan result){
	defer wg.Done()
	//log.Println(c + " starting")
	command := strings.Split(c, " ")
	if len(command) < 2 {
		// TODO: handle error
	}
	cmd := exec.Command(command[0], command[1:]...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		// TODO: handle error more gracefully
        //log.Fatal(err)
        results <- result{err: c}
        return
	}
	// do something with output
	//log.Println(c + " done")
    results <- result{data:  string(stdoutStderr)}
}


func cmd(c string) string{

	command := strings.Split(c, " ")
	if len(command) < 2 {
		// TODO: handle error
	}
	cmd := exec.Command(command[0], command[1:]...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		// TODO: handle error more gracefully
        //log.Fatal(err)
        //log.Println("cmd broke")
	}
	// do something with output
	//log.Println(c + " done")
	return string(stdoutStderr)
}


func get_cpu() int{
    idle0, total0 := get_CPU_Sample()
    time.Sleep(3 * time.Second)
    idle1, total1 := get_CPU_Sample()

    idleTicks := float64(idle1 - idle0)
    totalTicks := float64(total1 - total0)
    cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

    log.Printf("CPU usage is %f%% [busy: %f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
    return int(cpuUsage)
}

//----------------------------------------check cpu

func get_CPU_Sample() (idle, total uint64) {
    contents, err := ioutil.ReadFile("/proc/stat")
    if err != nil {
        return
    }
    lines := strings.Split(string(contents), "\n")
    for _, line := range(lines) {
        fields := strings.Fields(line)
        if fields[0] == "cpu" {
            numFields := len(fields)
            for i := 1; i < numFields; i++ {
                val, err := strconv.ParseUint(fields[i], 10, 64)
                if err != nil {
                    log.Println("Error: ", i, fields[i], err)
                }
                total += val // tally up all the numbers to get total ticks
                if i == 4 {  // idle is the 5th field in the cpu line
                    idle = val
                }
            }
            return
        }
    }
    return
}

//------------------------------------routines
var n = flag.Int("n", 1e5, "Number of goroutines to create")

var ch = make(chan byte)
var counter = 0

func f() {
    counter++
    <-ch // Block this goroutine
}

func test_routes_main() {
    flag.Parse()
    if *n <= 0 {
            log.Println(os.Stderr, "invalid number of goroutines")
            os.Exit(1)
    }

    // Limit the number of spare OS threads to just 1
    runtime.GOMAXPROCS(1)

    // Make a copy of MemStats
    var m0 runtime.MemStats
    runtime.ReadMemStats(&m0)

    t0 := time.Now().UnixNano()
    for i := 0; i < *n; i++ {
            go f()
    }
    runtime.Gosched()
    t1 := time.Now().UnixNano()
    runtime.GC()

    // Make a copy of MemStats
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)

    if counter != *n {
            log.Println(os.Stderr, "failed to begin execution of all goroutines")
            os.Exit(1)
    }

    log.Printf("Number of goroutines: %d\n", *n)
    log.Printf("Per goroutine:\n")
    log.Printf("  Memory: %.2f bytes\n", float64(m1.Sys-m0.Sys)/float64(*n))
    log.Printf("  Time:   %f µs\n", float64(t1-t0)/float64(*n)/1e3)
}
