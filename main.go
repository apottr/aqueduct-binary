package main

import (
  "fmt"
  "strings"
  "regexp"
  "errors"
  "encoding/json"
  "os"
  )

type Transfer struct {
  data string
}
 
var Functions = map[string]func(map[string]interface{}, Transfer) (Transfer,error){}

func graphExec(graph string, id string){
  lines := strings.Split(graph," -> ")
  var last Transfer
  var err error
  re := regexp.MustCompile(`([\w.]+)(\(.+\))?`)
  for _,v := range lines {
    cmd := re.FindStringSubmatch(v)
    if cmd[2] != "" && cmd[2][:1] == "("{
      var parsed map[string]interface{}
      json.Unmarshal([]byte(cmd[2][1:len(cmd[2])-1]),&parsed)
      parsed["argv"] = id
      last,err = Functions[cmd[1]](parsed,last)
      if err != nil {
        panic(fmt.Sprintf("Node threw error: %v",err))
      }
    }else{
      last,err = Functions[cmd[1]](map[string]interface{}{"empty":"","argv": id},last)
      if err != nil {
        panic(fmt.Sprintf("Node threw error: %v",err))
      }
    }
  }
}

func parseFirstNode(graph string) map[string]interface{} {
  var parsed map[string]interface{}
  node := strings.Split(graph," -> ")[0]
  re := regexp.MustCompile(`([\w.]+)(\(.+\))?`)
  cmd := re.FindStringSubmatch(node)
  if cmd[2] != "" && cmd[2][:1] == "(" {
    json.Unmarshal([]byte(cmd[2][1:len(cmd[2])-1]),&parsed)
    return parsed
  }else{
    return map[string]interface{}{}
  }
}

func getEnv(key string) string {
  val, ok := os.LookupEnv(key)
  if !ok {
    return ""
  }else{
    return val
  }
}
func main(){
  svr := getEnv("SERVER")
  if svr == "" {
    panic(errors.New("Provide SERVER environment variable"))
  }
  gph := getEnv("G")
  if gph == "" {
    panic(errors.New("Provide G environment variable"))
  }
  graph,err := getGraph(svr,gph)
  if err != nil {
    panic(err)
  }
  fmt.Println(graph.Graph)
  node := parseFirstNode(graph.Graph)
  changes, err := getChanges(node["server"].(string),node["database"].(string),graph.LastSeq)
  if err != nil {
    panic(err)
  }
  fmt.Println(changes.Changes)
  for i,item := range changes.Changes {
    graphExec(graph.Graph,item)
  }
  graph.LastSeq = changes.Last
  err = saveGraph(svr,gph,graph)
  if err != nil {
    panic(err)
  }
}
