package main

import (
	"flag"
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "log"
)

const (
    CstAddr = "http://10.20.11.17:1234/"
    CstDownload = "http://10.20.11.17:1234/download/"
)


// sample usage
func main() {
	nsel := flag.Int("sel",0,"choice functon")

    flag.Parse()
    log.Println("choice sel=",*nsel)
    switch *nsel {
    case 0:
        delSft()
    case 1:
        getlstAPP()
    default:
        panic("undefine parameter")
    }
}



func getlstAPP(){
    response,_ := http.Get(CstAddr+"getlstapp/")
    defer response.Body.Close()
    body,_ := ioutil.ReadAll(response.Body)
    fmt.Println(string(body))
}

func delSft(){
    tmp := `{"appName":"tst1","appMd5":"123abc"}`
    req := bytes.NewBuffer([]byte(tmp))

    body_type := "application/json;charset=utf-8"    
    resp, _ := http.Post(CstAddr+"delsoft/", body_type, req)
    //http.NewRequest("POST", CstDownload+"delsoft/", req_new)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))

    resp.Body.Close()
}