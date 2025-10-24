package main 

import (
    "fmt"
    "io"
    "github.com/fmo/resilienthttp"
)

func main() {
    res, err := resilienthttp.Get("http://localhost:8001/goals-list")
    if err != nil {
        panic(err)
    }
    body, _ := io.ReadAll(res.Body)
    
    fmt.Println(string(body))
}
