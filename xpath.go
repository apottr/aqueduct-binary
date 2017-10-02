package main

import (
  "github.com/clbanning/mxj"
  "gopkg.in/xmlpath.v2"
  "encoding/json"
  "strings"
  "bytes"
  "fmt"
)

func makeJSONXML(j string) ([]byte,error) {
  m,err := mxj.NewMapJson([]byte(j))
  if err != nil {
    return []byte(""),err
  }
  x,xrr := m.Xml()
  if xrr != nil {
    return []byte(""),xrr
  }
  x = bytes.Replace(x,[]byte("&"),[]byte("&amp;"),-1)
  return x,nil
}

func execXpath(xml []byte, xpath string, id string) ([]string,error) {
  path := xmlpath.MustCompile(xpath)
  root, err := xmlpath.Parse(bytes.NewReader(xml))
  if err != nil {
    return []string{},err
  }
  out := []string{}
  iter := path.Iter(root)
  idx := 0
  for iter.Next() {
    out = append(out,fmt.Sprintf("%d-%s @$^* %s",idx,id,iter.Node().String()))
    idx = idx + 1
  }
  return out,nil
}

func xpath_module(arg map[string]interface{}, last Transfer) (Transfer,error) {
  var xml []byte
  var err error
  if strings.HasPrefix(last.data,"{"){
    xml,err = makeJSONXML(last.data)
    if err != nil {
      return Transfer{},err
    }
  }else{
    xml = []byte(last.data)
  }
  if val, ok := arg["debug"]; ok {
    if val == "xml" {
      fmt.Printf("%s\n",xml)
    }
  }
  out,xrr := execXpath(xml,arg["q"].(string),arg["argv"].(string))
  if xrr != nil {
    return Transfer{},xrr
  }
  return Transfer{data: strings.Join(out," !#%& ")},nil
}

func main_standby(){
  var parsed struct{Id string `json:"_id"`
                    Rev string `json:"_rev"`
                    Data string `json:"data"`}
  doc,_ := getDoc("http://localhost:5984/","testgraph_out","ed520db4-582a-4fa2-a073-e134b49e5a0a")
  json.Unmarshal(doc,&parsed)
  xml,err := makeJSONXML(parsed.Data)
  if err != nil {
    panic(err)
  }
  out,xrr := execXpath(xml,"/doc/statuses/*/text","testId")
  if xrr != nil {
    panic(xrr)
  }
  fmt.Printf("%q\n",out)
}

func init(){
  Functions["xpath"] = xpath_module
}
