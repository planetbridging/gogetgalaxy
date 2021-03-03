package main

import (

    "log"
	"os/exec"
	"strings"
	"fmt"
	"net"
	//"encoding/xml"
	"encoding/xml"
	"strconv"
	//"time"
	//"net/http"
)


var network[] ObjPC

var findAll_status string

func bewear_setup(){
	findAll_status = "ready"
}


func test(){
	c := "nmap -sS -p 53,80,443,8081,8443 -sV -vvv -oX - 192.168.1.1"
	cback := cmd(c)
	

	var result NmapRun

	err := xml.Unmarshal([]byte(cback), &result)

    if err != nil {
    	fmt.Println(err)
	}

	var (
		osName     string
		osAccuracy = 0
	)
	
	for _, host := range result.Hosts {
		if host.Status.State == "up" {

			for _, osMatch := range host.Os.OsMatches {
				tempOsAccuracy, _ := strconv.Atoi(osMatch.Accuracy)
				if tempOsAccuracy >= osAccuracy {
					osName = osMatch.Name
					osAccuracy = tempOsAccuracy
				}
			}

			ipAddr := host.Addresses[0].Addr
			for _, port := range host.Ports {
				portStr := strconv.Itoa(port.PortId)
				servicesStr := port.Service.CPEs
				fmt.Println(ipAddr, portStr, servicesStr, osName)
			}
		}
	}
}

//----------------------------------------------default scans

func getLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

func findAll(){
	findAll_status = "busy"
	localip := getLocalIP()

	icmp_pe := "nmap -PE -sn -oG - "+localip+"/24"
	icmp_pp := "nmap -PP -sn -oG - "+localip+"/24"
	icmp_pm := "nmap -PM -sn -oG - "+localip+"/24"

	icmp_pe_c := cmd(icmp_pe)
	icmp_pp_c := cmd(icmp_pp)
	icmp_pm_c := cmd(icmp_pm)
	
	lst_icmp_pe_c := sortOGICMPScan(icmp_pe_c)
	lst_icmp_pp_c := sortOGICMPScan(icmp_pp_c)
	lst_icmp_pm_c := sortOGICMPScan(icmp_pm_c)

	addLstIp(lst_icmp_pe_c)
	addLstIp(lst_icmp_pp_c)
	addLstIp(lst_icmp_pm_c)
	findAll_status = "update"
}

func findAllPorts(){
	for _, n := range network {
		ports_ss := "nmap -sS -oG - " + n.ip
		ports_st := "nmap -sT -oG - " + n.ip
		ports_sa := "nmap -sA -oG - " + n.ip

		ports_ss_c := cmd(ports_ss)
		ports_st_c := cmd(ports_st)
		ports_sa_c := cmd(ports_sa)

		lst_ports_ss_c := sortOGPortScan(ports_ss_c)
		lst_ports_st_c := sortOGPortScan(ports_st_c)
		lst_ports_sa_c := sortOGPortScan(ports_sa_c)

		tmp_l_ports := ""

		tmp_l_ports += strings.Join(lst_ports_ss_c, ",") 
		tmp_l_ports += "," + strings.Join(lst_ports_st_c, ",") 
		tmp_l_ports += "," + strings.Join(lst_ports_sa_c, ",") 

		l_ports := strings.Split(tmp_l_ports,",")

		uniq_ports := unique(l_ports)

    	fmt.Println(uniq_ports)
	}
}

//----------------------------------------------default scans

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

func addLstIp(ip[] string){
	for _, i := range ip {
		if len(network) == 0{
			tmp_ObjPC := ObjPC{
				ip: i,
				ssh_status: "notready",
			}
			//fmt.Println(len(tmp_ObjPC.lstport))
			network = append(network, tmp_ObjPC)
		}else{
			found := false
			for _, n := range network {
				if n.ip == i{
					found = true
					break
				}
			}
			if !found{
				tmp_ObjPC := ObjPC{
					ip: i,
					ssh_status: "notready",
				}
				network = append(network, tmp_ObjPC)
			}
		}
	}
}

func unique(e []string) []string {
    r := []string{}

    for _, s := range e {
        if !contains(r[:], s) {
            r = append(r, s)
        }
    }
    return r
}

func contains(e []string, c string) bool {
    for _, s := range e {
        if s == c {
            return true
        }
    }
    return false
}

//----------------------------------------------scan extract
func sortOGPortScan(d string)[] string{
	l_ports := []string{}
	for _, line := range strings.Split(strings.TrimSuffix(d, "\n"), "\n") {
		//fmt.Println(line)
		if strings.Contains(line, "Ports:"){
			l_split := strings.Split(line," ")
			for _, l := range l_split{
				if strings.Contains(l, "/"){
					l_split_ports := strings.Split(l,"/")
					//fmt.Println(l_split_ports[0])
					l_ports = append(l_ports, l_split_ports[0])
				}
			}
		}
	}
	return l_ports
}

func sortOGICMPScan(d string)[] string{
	l_ip := []string{}
	for _, line := range strings.Split(strings.TrimSuffix(d, "\n"), "\n") {
		//fmt.Println(line)
		if strings.Contains(line, "Up"){
			l_split := strings.Split(line," ")
			//fmt.Println(l_split[1])
			l_ip = append(l_ip, l_split[1])
		}
	}
	return l_ip
}


//----------------------------------------------scan extract