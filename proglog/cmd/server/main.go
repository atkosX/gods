package main

import(
    "log"
    "github.com/beediBiceps/proglog/internal/server"
)

func main(){
    svr:=server.NewHttpServer(":8080")
    log.Fatal(svr.ListenAndServe())
}