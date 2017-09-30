package main

import (
  "net/http"
  "encoding/json"
  "io/ioutil"
  "fmt"
  "bytes"
  "errors"
)

type Config struct {
  Id string `json:"_id"`
  Rev string `json:"_rev"`
  Graph string `json:"graph"`
}

type AddResponse struct {
  Id string `json:"id"`
  Ok string `json:"ok"`
  Rev string `json:"rev"`
}

var client = &http.Client{}

func getGraph(server string, ident string) (string,error) {
  var parsed Config
  resp,err := client.Get(fmt.Sprintf("%sgraph_configs/%s",server,ident))
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

func getDoc(server string, database string, ident string) ([]byte,error) {
  resp,err := client.Get(fmt.Sprintf("%s%s/%s",server,database,ident))
  if err != nil {
    return []byte(""),err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return []byte(""),err
  }
  return body,nil
}

func doesDatabaseExist(server string, database string) (bool,error){
  resp,err := client.Head(fmt.Sprintf("%s%s",server,database))
  if err != nil {
    return false,err
  }
  defer resp.Body.Close()
  if resp.StatusCode != 200 {
    return false,nil
  }else{
    return true,nil
  }
}

func createDatabase(server string, database string) error {
  req,err := http.NewRequest("PUT",fmt.Sprintf("%s%s",server,database),nil)
  if err != nil {
    return err
  }
  res, er := client.Do(req)
  if er != nil {
    return er
  }
  if res.StatusCode > 202 && res.StatusCode != 412 {
    return errors.New(res.Status)
  }
  defer res.Body.Close()
  return nil
}

func addDoc(server string, database string, ident string, doc string) (string,error) {
  err := createDatabase(server,database)
  if err != nil {
    return "",err
  }
  body := bytes.NewReader([]byte(doc))
  req,e := http.NewRequest("PUT",fmt.Sprintf("%s%s/%s",server,database,ident),body)
  if e != nil {
    return "",e
  }
  res, er := client.Do(req)
  if er != nil {
    return "",er
  }
  defer res.Body.Close()
  if res.StatusCode > 203 {
    return "",errors.New(res.Status)
  }else{
    var r AddResponse
    body, erro := ioutil.ReadAll(res.Body)
    if erro != nil {
      return "",erro
    }
    err = json.Unmarshal(body,&r)
    return r.Id,nil
  }
}
