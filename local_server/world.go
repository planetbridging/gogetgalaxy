package main

import (
    "fmt"
    "os"
    "path/filepath"
    //"flag"
    "log"
    "bufio"
    "strings"
    "strconv"
    "time"
)

type ObjPort struct{
  port int
  cpe[] string
  details[] string
}

type ObjComputer struct{
  ip string
  o_port []ObjPort
}

type ObjCountry struct{
  country string
  LstComputers[] ObjComputer
  count int
}

var LstWorld[] ObjCountry



/*func main() {
  //all ip in world 4,294,967,296
  //read_File()
  //1.0.4.0  	1.0.7.255       1024
  //1.120.0.0	1.159.255.255   2,621,440

  //ipadd := ip_range("0.0.0.0","255.255.255.255")
  //fmt.Println(len(ipadd))

  //strVar := "100"
	//intVar, _ := strconv.Atoi(strVar)
  //fmt.Println(intVar)
}*/

func ip_range(start string, stop string)[]string{

  var lstip []string

  start_split := strings.Split(start, ".")
  stop_split := strings.Split(stop, ".")

  //fmt.Println(start_split)
  //fmt.Println(stop_split)

  it1, _ := strconv.Atoi(start_split[0])
  it2, _ := strconv.Atoi(start_split[1])
  it3, _ := strconv.Atoi(start_split[2])
  it4, _ := strconv.Atoi(start_split[3])

  it5, _ := strconv.Atoi(stop_split[0])
  it6, _ := strconv.Atoi(stop_split[1])
  it7, _ := strconv.Atoi(stop_split[2])
  it8, _ := strconv.Atoi(stop_split[3])

  itb1 := false
  itb2 := false
  itb3 := false
  itb4 := false

  for it1 < 256 {

    for it2 < 256 {


      for it3 < 256 {


        for it4 < 256 {

          //fmt.Println(it1,it2,it3,it4)

          ip1 := strconv.Itoa(it1)
          ip2 := strconv.Itoa(it2)
          ip3 := strconv.Itoa(it3)
          ip4 := strconv.Itoa(it4)

          ipadd := ip1 + "." + ip2 + "." + ip3 + "." + ip4


          lstip = append(lstip,ipadd)

          if it1 == it5{
            if it2 == it6{
              if it3 == it7{
                if it4 == it8{
                  itb1 = true
                  itb2 = true
                  itb3 = true
                  itb4 = true

                  break
                }
              }
            }
          }
          if itb4{
            break
          }
          it4 += 1
        }
        if itb3{
          break
        }
        it4 = 0
        it3 += 1
      }
      if itb2{
        break
      }
      it3 = 0
      it2 += 1
    }
    if itb1{
      break
    }
    it2 = 0
    it1 += 1
  }
  return lstip
}

func get_All_Files(){
  var files []string

  dirroot := "countries/"
  err := filepath.Walk(dirroot, func(path string, info os.FileInfo, err error) error {
      files = append(files, path)
      return nil
  })
  if err != nil {
      panic(err)
  }
  for _, file := range files {
      fmt.Println(file)
  }
}

func read_File_test(){

  obj_country := ObjCountry{

  }

  file, err := os.Open("countries/Australia.txt")
  if err != nil {
      log.Fatal(err)
  }
  defer file.Close()
  total_count := 0
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
      //fmt.Println(scanner.Text())
      line := strings.Split(scanner.Text(),":")
      between := strings.Split(line[1],"-")
      //fmt.Println(between[0] + "---" + between[1])

      ipadd := ip_range(between[0],between[1])

      obj_port := ObjPort{
        port: 80,
        cpe: []string{"a:lighttpd:lighttpd:1.4.2","a:lighttpd:lighttpd:1.4.2","a:lighttpd:lighttpd:1.4.2"},
      }

      tmp_ports := []ObjPort{obj_port,obj_port,obj_port,obj_port,obj_port,obj_port,obj_port}

      for i := 0; i < 30000; i++ {
        tmp_ports = append(tmp_ports,obj_port)
      }

      for _, i := range ipadd {
        tmp_ip := ObjComputer{
          ip: i,
          o_port: tmp_ports,
        }
        obj_country.LstComputers = append(obj_country.LstComputers,tmp_ip )
      }

      total_count += len(ipadd)
      //fmt.Println(len(ipadd))
  }

  obj_country.count = total_count

  if err := scanner.Err(); err != nil {
      log.Fatal(err)
  }

  fmt.Println("Total: " + strconv.Itoa(obj_country.count))
  time.Sleep(10 * time.Second)
}


func load_country(tmp_country string){

  obj_country := ObjCountry{

  }

  file, err := os.Open("countries/"+tmp_country+".txt")
  if err != nil {
      log.Fatal(err)
  }
  defer file.Close()
  total_count := 0
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
      //fmt.Println(scanner.Text())
      line := strings.Split(scanner.Text(),":")
      between := strings.Split(line[1],"-")
      //fmt.Println(between[0] + "---" + between[1])

      ipadd := ip_range(between[0],between[1])

      for _, i := range ipadd {
        tmp_ip := ObjComputer{
          ip: i,
          //o_port: tmp_ports,
        }
        obj_country.LstComputers = append(obj_country.LstComputers,tmp_ip )
      }

      total_count += len(ipadd)
      //fmt.Println(len(ipadd))
  }

  obj_country.count = total_count

  if err := scanner.Err(); err != nil {
      log.Fatal(err)
  }

  fmt.Println("Total: " + strconv.Itoa(obj_country.count))
  time.Sleep(10 * time.Second)
}