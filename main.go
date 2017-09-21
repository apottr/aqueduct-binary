package main

import (
  "fmt"
  "strings"
  "regexp"
  "encoding/json"
  "os"
  )

type Transfer struct {
  data string
}
 
var Functions = map[string]func(map[string]interface{}, Transfer) Transfer{}

func graphExec(graph string){
  lines := strings.Split(graph," -> ")
  var last Transfer
  re := regexp.MustCompile(`([\w.]+)(\(.+\))?`)
  for _,v := range lines {
    cmd := re.FindStringSubmatch(v)
    if cmd[2] != "" && cmd[2][:1] == "("{
      var parsed map[string]interface{}
      json.Unmarshal([]byte(cmd[2][1:len(cmd[2])-1]),&parsed)
      last = Functions[cmd[1]](parsed,last)
    }else{
      last = Functions[cmd[1]](map[string]interface{}{"empty":""},last)
    }
  }
}

func main(){
  graph,err := getGraph(os.Args[1],os.Args[2])
  if err != nil {
    panic(err)
  }
  fmt.Println(graph)
  graphExec(graph)
}
