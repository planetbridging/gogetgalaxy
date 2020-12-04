package main

import (
	"strings"
	"fmt"

	"strconv"
	//"time"
	"net/http"
	"io/ioutil"
  "image/jpeg"
  "bytes"
  "image"
  "flag"
  //"sync"
  "encoding/json"
)


type nodes struct{
	Name string `json:"name"`
}

type source struct{
	Source string `json:"source"`
	Target string `json:"target"`
	Value string `json:"value"`
}


type response2 struct {
	Links []source `json:"links"`
	Nodes []nodes `json:"nodes"`
}

func testjson(){

	s1 := source{
		Source: "Sylva",
		Target: "Compound",
		Value: "10.5",
	}

	s1_jsn, _ := json.Marshal(s1)
    fmt.Println(string(s1_jsn))


	lsts1 := []source{s1}

	res2D := response2{
        Links: lsts1}
    res2B, _ := json.Marshal(res2D)
    fmt.Println(string(res2B))
}

var root = flag.String("root", ".", "file system path")

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page: ", r.URL.Path)
	
	//icmp1, okicmp1 := r.URL.Query()["icmp1"]
	/*
	if r.URL.Path == "/scan"{
		s := "no idea what your doing lol"
		if okicmp1{
			//s =startICMP(icmp1,"1")
		}else if okicmp2{
			//s =startICMP(icmp2,"2")
		}else if okicmp3{
			//s =startICMP(icmp3,"3")
		}else if okicmpall{
			//s = startICMP(icmpall,"4")
		}

		fmt.Fprint(w,s)
		//page_loaded = true
	}*/

	if r.URL.Path == "/countries"{
		cards := "<div class='card-group'>"
		for _, s := range LstWorld {
			//fmt.Println(i, s.country)
			cards += country_card(s.country,s.count)
		}

		cards += "</div>"
		fmt.Fprint(w,cards)
	}else if r.URL.Path == "/ready"{
		//fmt.Fprint(w,get_available())
	}

}

func country_card(country string, count int) string{
	card := "<div class='card text-white bg-dark mb-3' style='width: 18rem;'>"
	card += "<div class='card-header'>" + country + "</div>"
	card += "<ul class='list-group list-group-flush'>";

	card += "<li class='list-group-item text-white bg-dark'>Count: "+strconv.Itoa(count)+"</li>"

	card += "</ul></div>";
	return card
}

func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		//log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		//log.Println("unable to write image.")
	}
}

func ReadFile(filename string) []string{
	data, err := ioutil.ReadFile(filename)
	var LstFile[]string
	if err != nil {
		fmt.Println("File reading error", err)
		return LstFile
	}
  
	for _, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
	  LstFile = append(LstFile, line)
	}
	return LstFile
  }