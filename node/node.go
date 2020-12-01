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
var ports_ss = "nmap -Pn -sS -oG - -p- -T4 --max-parallelism 100 --min-rate 10000 "
var ports_st = "nmap -Pn -sT -oG - -p- -T4 --max-parallelism 100 --min-rate 10000 "
var ports_sa = "nmap -Pn -sA -oG - -p- -T4 --max-parallelism 100 --min-rate 10000 "
/*

nmap --mtu 64 -p- -T4 -sS --max-parallelism 100 --min-rate 10000 192.168.1.1 	
nmap -Pn -sS --mtu 64 -p - --max-parallelism 100 -vvv --min-rate 10000 pressback.space



c := "nmap -sS -p 53,80,443,8081,8443 -sV -T4 -oX - 192.168.1.1"

*/


func main(){
    scan_ip := "192.168.1.1"
    /*found := ping(scan_ip)
    if found{
       
    }*/
    port_scan(scan_ip)
   //fmt.Println(cmd(ports_sa+scan_ip))
}


func merge(channelList []chan string) <-chan string {
    var wg sync.WaitGroup
    out := make(chan string)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan string) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    wg.Add(len(channelList))
    for _, c := range channelList {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
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
	//fmt.Println(lstcmds)
	/*messages := make(chan string,len(lstcmds))
	lst_results := []string{}
    var wg sync.WaitGroup

	for _, i := range lstcmds {
	//for i := 0; i < len(lstcmds); i++{
		//fmt.Println(i + "ffs")
		tmp := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			//time.Sleep(time.Second * 3)
			
			cmd_wait := cmd(tmp)
			messages <- cmd_wait
		}()
	}

	go func() {
        for i := range messages {
			//fmt.Println(i)
			lst_results = append(lst_results,i)
        }
	}()

    wg.Wait()*/
    
    
    /*for _, i := range lstcmds {
        c1 := make(chan string, 1)
	
		go func() {
            // simulate spending time to do work to get answer
			cmd_wait := cmd(i)
			c1 <- cmd_wait
		}()

        x := <-c1
        lst_results = append(lst_results,x)
		fmt.Println(x)
    }*/

    /*lst_results := []string{}

    runtime.GOMAXPROCS(len(lstcmds))
    var processTest sync.WaitGroup
    processTest.Add(len(lstcmds))
     
    for _, i := range lstcmds {
        tmp_wait := make(chan string, 1)
        go func(tmp_i string) {
            defer processTest.Done()
            tmp_wait <- cmd(tmp_i)
           // lst_results = append(lst_results,cmd_wait)
        }(i)
        lst_results = append(lst_results, <- tmp_wait)
    }
    
 
    // Wait for the goroutines to finish.
    processTest.Wait() */

    /*var wg sync.WaitGroup
    var lst_results [] string
    var channelList []chan string

    for i, c := range lstcmds {
        wg.Add(1)
        //go cmd_wait(ch, c,&wg)
        go func(tmp_c string,chan_num int) {
            defer wg.Done()
            tmp_cmd := cmd(tmp_c)
            channelList[chan_num] <- tmp_cmd
            //time.Sleep(time.Second * 3)
            //cmd_wait(ch, tmp_c,&wg)
        }(c,i)
    }

    go func() {
        for i := range channelList {
            fmt.Println(i)
        }
	}()

    wg.Wait()*/
    
    /*lst_channels := make([]int, 5)
	for i, c := range lstcmds {
		go func(tmp_c string,num int) {
			//time.Sleep(time.Second * 3)
			cmd_wait := cmd(tmp_c)
			lst_channels[num] <- cmd_wait
        }(c,i)
        lst_results = append(lst_results,<-tmp_channel)
	}

	go func() {
        for _,i := range channelList {
            lst_results = append(lst_results,<-i)
        }
    }()*/
    /*var channelList []chan struct{}
	for i := 0; i < 10; i++ {
		channel := make(chan struct{})
		channelList = append(channelList, channel)
		go simulateProcessingThenClose(channel)

	}

	waitForChannelsToClose(channelList...)

	// Or use individual channels
	ch1 := make(chan struct{})
    ch2 := make(chan struct{})
    
    //var resultsList []chan string

	go simulateProcessingThenClose(ch1)
	go simulateProcessingThenClose(ch2)

	waitForChannelsToClose(ch1,ch1,ch1,ch1)

    fmt.Println("That's all for now!")*/
    /*var lst_results [] string
    var lst_channels []chan string
    //var wg sync.WaitGroup
	for _, c := range lstcmds {
        //wg.Add(1)
		channel := make(chan string)
		lst_channels = append(lst_channels, channel)
		/*go func(tmp_c string) {
			defer wg.Done()
			//time.Sleep(time.Second * 3)
			cmd_wait := cmd(tmp_c)
            channel <- cmd_wait
            //close(channel)
        }(c)*
        go simulateProcessingThenClose(channel,c)
	}

	t := time.Now()
	for _, v := range lst_channels {
        <-v
        //lst_results = append(lst_results,tmp_results)

		fmt.Printf("%v for scan to close\n", time.Since(t))
	}
    fmt.Printf("%v for channels to close\n", time.Since(t))*/
    
    /*var lst_results [] string
    var wg sync.WaitGroup
    var lst_channels []chan string

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(tmp_c string) {
        cmd_wait := cmd(tmp_c)
        channel := make(chan string)
        lst_channels = append(lst_channels, channel)
        channel <- cmd_wait
        wg.Done()
    }
    wg.Add(len(lstcmds))

    for _, c := range lstcmds {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        
    }()*/

    var lst_results [] result

   /* for _, c := range lstcmds {

        c1 := make(chan string)
        go func(tmp_c string) {
			//time.Sleep(time.Second * 3)
			tmp_cmd_wait := cmd(tmp_c)
            //messages <- cmd_wait
            //lst_results = append(lst_results,cmd_wait)
            //fmt.Println(tmp_cmd_wait)
            //c1 <- tmp_cmd_wait
        }(c)
        x := <-c1
		fmt.Println(x)
    }*/
    //var lst_channels []chan string
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
