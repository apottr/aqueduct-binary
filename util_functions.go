package main

import (
  "fmt"
  "github.com/satori/go.uuid"
  "encoding/json"
)

/*func addConfig(cfg map[string]string, function func(map[string]interface{}, Transfer) Transfer){
  for _,doc := range allDocs("http://localhost:5984/","function_store") {
    if doc == cfg {
      break
    }
  }
}*/

func out_function (arg map[string]interface{}, last Transfer) (Transfer,error) {
    fmt.Println(last.data)
    return Transfer{},nil
}

func debug_print(arg map[string]interface{}, last Transfer) (Transfer,error) {
  fmt.Println(last.data)
  return last,nil
}

func debug_args(arg map[string]interface{}, last Transfer) (Transfer,error){
  fmt.Println(arg)
  return last,nil
}

func in_couchdb(arg map[string]interface{}, last Transfer) (Transfer,error) {
  doc, err := getDoc(arg["server"].(string),arg["database"].(string),arg["id"].(string))
  if err != nil {
    return Transfer{},err
  }
  var parsed struct{
    Contents string `json:"contents"`
    Type string `json:"type"`
  }
  err = json.Unmarshal(doc,&parsed)
  if err != nil {
    return Transfer{},err
  }
  return Transfer{data: parsed.Contents},nil
}

func in_couchdb_args(arg map[string]interface{}, last Transfer) (Transfer,error){
  doc, err := getDoc(arg["server"].(string),arg["database"].(string),arg["argv"].(string))
  fmt.Println(string(doc))
  if err != nil {
    return Transfer{},err
  }
  var parsed struct{
    Contents string `json:"contents"`
    Type string `json:"type"`
  }
  err = json.Unmarshal(doc,&parsed)
  if err != nil {
    return Transfer{},err
  }
  return Transfer{data: parsed.Contents},nil
}

func out_couchdb(arg map[string]interface{}, last Transfer) (Transfer,error) {
  var doc = map[string]string{}
  doc["data"] = string(last.data)
  s_doc,e := json.Marshal(doc)
  if e != nil {
    return Transfer{},e
  }
  _,err := addDoc(arg["server"].(string),arg["database"].(string),uuid.NewV4().String(),string(s_doc))
  if err != nil {
    return Transfer{},err
  }
  return Transfer{},nil
}

func init(){
  Functions["in.couchdb"] = in_couchdb
  Functions["in.couchdb.args"] = in_couchdb_args
  Functions["out.print"] = out_function
  Functions["dbg.data"] = debug_print
  Functions["dbg.args"] = debug_args
  Functions["out.couchdb"] = out_couchdb
}
