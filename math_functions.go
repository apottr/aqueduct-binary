package main

import (
  "strconv"
)

func math_add (arg map[string]interface{}, last Transfer) (Transfer,error) {
  t := Transfer{}
  lst,err := strconv.ParseFloat(last.data,64)
  if err != nil {
    return Transfer{},err
  }
  t.data = strconv.FormatFloat(arg["data"].(float64)+lst,'f',-1,64)
  return t,nil

}

func math_subtract (arg map[string]interface{}, last Transfer) (Transfer,error) {
  t := Transfer{}
  lst,err := strconv.ParseFloat(last.data,64)
  if err != nil {
    return Transfer{},err
  }
  t.data = strconv.FormatFloat(lst-arg["data"].(float64),'f',-1,64)
  return t,nil

}

func math_in(arg map[string]interface{}, last Transfer) (Transfer,error) {
  t := Transfer{}
  t.data = strconv.FormatFloat(arg["data"].(float64),'f',-1,64)
  return t,nil
}

func init(){
  Functions["math.add"] = math_add
  Functions["math.subtract"] = math_subtract
  Functions["math.in"] = math_in
}
