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
)

//static nmap cmds
//ping
var icmp_pe = "nmap -PE -sn -oG - "
var icmp_pp = "nmap -PP -sn -oG - "
var icmp_pm = "nmap -PM -sn -oG - "
//port scan
var ports_ss = "nmap -sS -oG - --mtu 64 -p- -T4 --max-parallelism 100 --min-rate 10000"
var ports_st = "nmap -sT -oG - --mtu 64 -p- -T4 --max-parallelism 100 --min-rate 10000"
var ports_sa = "nmap -sA -oG - --mtu 64 -p- -T4 --max-parallelism 100 --min-rate 10000"
/*

nmap --mtu 64 -p- -T4 -sS --max-parallelism 100 --min-rate 10000 192.168.1.1 	
nmap -Pn -sS --mtu 64 -p - --max-parallelism 100 -vvv --min-rate 10000 pressback.space



c := "nmap -sS -p 53,80,443,8081,8443 -sV -T4 -oX - 192.168.1.1"

*/

func main(){

}


//----------------------------------------cmd

func multi_cmd_wait(){
	messages := make(chan string)
    var wg sync.WaitGroup

    // you can also add these one at 
    // a time if you need to 


	
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//time.Sleep(time.Second * 3)
			cmd_wait := cmd("nmap -p - -sS -oX - 192.168.1.1")
			messages <- cmd_wait
		}()
	}

	go func() {
        for i := range messages {
            fmt.Println(i)
        }
	}()

    wg.Wait()
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
