package main

import (
  "encoding/json"
  "log"
  "os"
  "fmt"
)

type MyEncoder struct {
  *json.Encoder
}

func (enc *MyEncoder) encode_print(av ...interface{}) {
  for v := range av {
    if err := enc.Encoder.Encode(&v); err != nil {
      log.Println(err)
    }
  }
}

func main() {
  dec := json.NewDecoder(os.Stdin)
  //myenc := &MyEncoder{json.NewEncoder(os.Stdout)}
  for {
    var v interface{}
    if err := dec.Decode(&v); err != nil {
      if err.Error() != "EOF" {
        log.Println(err)
      }
      return
    }
    av := v.([]interface{})
    for _, v := range av {
      m := v.(map[string]interface{})["value"].(map[string]interface{})["fields"].(map[string]interface{})
      fmt.Println(m["id"], "(", m["version"], ") has currency_code:", m["currency_code"], "with value of:", m["rate"], "for date:", m["date"])
    }
    
    /*m := v.([]interface{})
    for k, v := range m {
      switch vv := v.(type) {
      case string:
        fmt.Println(k, "is string", vv)
      case int:
        fmt.Println(k, "is int", vv)
      case []interface{}:
        fmt.Println(k, "is an array:")
        for i, u := range vv {
          fmt.Println(i, u)
        }
      case map[string]interface{}:
        for k := range vv {
          if k != "key" {
            delete(vv, k)
          }
        }
        for _, u := range vv {
          myenc.encode_print(u)
        }
      default:
        fmt.Println(k, "is of a type I don't know how to handle")
      }
    }*/
  }
}
