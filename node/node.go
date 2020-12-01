package main

import (
	//test routes max
    "flag"
    "fmt"
    "os"
    "runtime"
	"time"

    
	//test nmap max with cpu
	"io/ioutil"
    "strconv"
	"strings"
	
	//cmd
	"log"
	"os/exec"
    "sync"

    //tcp client
    "bufio"
	"io"
	"net"

)

type result struct {
    data string
    err string // you must use error type here
}

//static nmap cmds
//ping
var icmp_pe = "nmap -PE -sn -oG - "
var icmp_pp = "nmap -PP -sn -oG - "
var icmp_pm = "nmap -PM -sn -oG - "
//port scan
var ports_ss = "nmap -Pn -sS -oG - --min-rate 10000 --top-ports 1000 "
var ports_st = "nmap -Pn -sT -oG - --min-rate 10000 --top-ports 1000 "
var ports_sa = "nmap -Pn -sA -oG - --min-rate 10000 --top-ports 1000 "

func main(){
    /*scan_ip := "192.168.1.1"
    found := ping(scan_ip)
    if found{
       port_scan(scan_ip)
    }*/
    
   //fmt.Println(cmd(ports_sa+scan_ip))
   tcp_client()
}

func tcp_client() {
	con, err := net.Dial("tcp", "0.0.0.0:4849")
	if err != nil {
		log.Fatalln(err)
	}
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
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}
	}
}




//----------------------------------------cmd
func ping(ip string)bool{

	pe := icmp_pe + ip
	pp := icmp_pp + ip
	pm := icmp_pm + ip
	
	

	lst_tmp_ping:= []string{pe,pp,pm}
	fmt.Println(lst_tmp_ping)

    results := multi_cmd_wait(lst_tmp_ping)
	for _, r := range results{
        if strings.Contains(r,"Status: Up"){
            return true
        }
        //fmt.Println(r)
    }
    return false
}

func port_scan(ip string){
    ss := ports_ss + ip
	st := ports_st + ip
	sa := ports_sa + ip
    

	lst_tmp_port:= []string{ss,st,sa}
	//fmt.Println(lst_tmp_port)

    results := multi_cmd_wait(lst_tmp_port)
	for _, r := range results{
        fmt.Println(r)
    }
}

func multi_cmd_wait(lstcmds []string)[]string{
	

    var lst_results [] result
    var wg sync.WaitGroup
    ch := make(chan result)
	
	
    //fmt.Println(<-nums) // Read the value from unbuffered channel
    
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
	fmt.Println(c + " starting")
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
	fmt.Println(c + " done")
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
		log.Fatal(err)
	}
	// do something with output
	fmt.Println(c + " done")
	return string(stdoutStderr)
}


func get_cpu() {
    idle0, total0 := get_CPU_Sample()
    time.Sleep(3 * time.Second)
    idle1, total1 := get_CPU_Sample()

    idleTicks := float64(idle1 - idle0)
    totalTicks := float64(total1 - total0)
    cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

    fmt.Printf("CPU usage is %f%% [busy: %f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
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
                    fmt.Println("Error: ", i, fields[i], err)
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
            fmt.Fprintf(os.Stderr, "invalid number of goroutines")
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
            fmt.Fprintf(os.Stderr, "failed to begin execution of all goroutines")
            os.Exit(1)
    }

    fmt.Printf("Number of goroutines: %d\n", *n)
    fmt.Printf("Per goroutine:\n")
    fmt.Printf("  Memory: %.2f bytes\n", float64(m1.Sys-m0.Sys)/float64(*n))
    fmt.Printf("  Time:   %f µs\n", float64(t1-t0)/float64(*n)/1e3)
}
