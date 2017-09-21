package main

import (
  "net/http"
  "encoding/json"
  "io/ioutil"
  "fmt"
)

type Config struct {
  Id string `json:"_id"`
  Rev string `json:"_rev"`
  Graph string `json:"graph"`
}

func getGraph(server string, ident string) (string,error) {
  var parsed Config
  resp,err := http.Get(fmt.Sprintf("%sgraph_configs/%s",server,ident))
  if err != nil {
    return "",err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "",err
  }
  json.Unmarshal(body,&parsed)
  return parsed.Graph,nil
}

func doesDatabaseExist(server string, database string) (string,error){
  return "",nil
}

func addDoc(server string, ident string, doc string) (string,error) {
  return "",nil
}
