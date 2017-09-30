package main

import (
  "encoding/base64"
)

func base64_decode(args map[string]interface{}, last Transfer) (Transfer,error) {
  t := Transfer{}
  b,err := base64.StdEncoding.DecodeString(last.data)
  if err != nil {
    return Transfer{},err
  }
  t.data = string(b)
  return t,nil
}

func substring_function(args map[string]interface{}, last Transfer) (Transfer,error) {
  t := Transfer{}
  t.data = last.data[int(args["start"].(float64)):len(last.data)-int(args["end"].(float64))]
  return t,nil
}

func init(){
  Functions["base64.decode"] = base64_decode
  Functions["substr"] = substring_function
}
