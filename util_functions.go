package main

import "fmt"

func out_function (arg map[string]interface{}, last Transfer) Transfer {
    fmt.Println(last.data)
    return Transfer{}
}

func init(){
  Functions["out"] = out_function
}
