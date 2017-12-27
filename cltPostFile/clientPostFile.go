package main

import (
"net/url"
"bytes"
"fmt"
"io"
"io/ioutil"
"mime/multipart"
"net/http"
"os"
"log"
"strconv"
)

const (
    CstAddr = "http://localhost:1234/"
    CstDownload = "http://localhost:1234/download/"
)

func postFile(filename string, targetUrl string,path string) error {
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)

    //关键的一步操作
    fileWriter, err := bodyWriter.CreateFormFile("file", filename)
    if err != nil {
        fmt.Println("error writing to buffer")
        return err
    }

    //打开文件句柄操作
    fh, err := os.Open(path+filename)
    if err != nil {
        fmt.Println("error opening file")
        return err
    }

    //iocopy
    _, err = io.Copy(fileWriter, fh)
    if err != nil {
        return err
    }

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    resp, err := http.Post(targetUrl, contentType, bodyBuf)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    resp_body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    fmt.Println(resp.Status)
    fmt.Println(string(resp_body))
    return nil
}



func postFileEx(filename string, targetUrl string,path string) error {
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)
    

    //关键的一步操作
    fileWriter, err := bodyWriter.CreateFormFile("file", filename)
    if err != nil {
        fmt.Println("error writing to buffer")
        return err
    }

    //打开文件句柄操作
    fh, err := os.Open(path+filename)
    if err != nil {
        fmt.Println("error opening file ",err)
        return err
    }

    //iocopy
    _, err = io.Copy(fileWriter, fh)
    if err != nil {
        return err
    }

    fmx := url.Values{
        "name":{"企控通"},
        "time":{"2017年12月27日14:00:52"},
    }

    fmName,_ :=bodyWriter.CreateFormFile("desc","descx")
    strfmt :=fmx.Encode()
    fmName.Write([]byte(strfmt))


    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    
    resp, err := http.Post(targetUrl, contentType, bodyBuf)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    resp_body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    fmt.Println(resp.Status)
    fmt.Println(string(resp_body))
    return nil
}

func postMutiForm(){
    //var mutiRes bytes.Buffer
    //muti_writer :=multipart.NewWriter(&mutiRes)

    
}


// sample usage
func main() {
    //tstHello()
    //getFile("a.txt")
    tstPostFile()
}

func tstPostFile(){
    target_url := "http://localhost:1234/uploadx"
    //filename := ".\\ax1.pdf"
    path := ".\\"//os.Args[1]
    filename := "a.txt"//os.Args[2]
    fmt.Println(path,filename)
    postFileEx(filename, target_url,path)
}

func tstHello(){
    acx := "hello"
    log.Println("ask------",1)
    res,err := http.Get(CstAddr+acx)   
    log.Println("ask------",2)

    if err != nil {
        log.Println("ask------",3)
        log.Println(acx," 请求失败 ",err)
        return
    }
    defer res.Body.Close()

    log.Println("ask------",4)
    body,_ := ioutil.ReadAll(res.Body)
    log.Println("ask------",5)
    log.Println(string(body))
}

func getFile(t_file string){

    f, err := os.OpenFile(t_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    stat, err := f.Stat() //获取文件状态
    if err != nil { panic(err) } //把文件指针指到文件末，当然你说为何不直接用 O_APPEND 模式打开，没错是可以。我这里只是试验。

    ask,err := http.NewRequest("GET",CstDownload+t_file,nil)
    ask.Header.Set("Range", "bytes=" + strconv.FormatInt(stat.Size(),10) + "-")
    resp, err := http.DefaultClient.Do(ask)
    if err != nil { panic(err) }
    written, err := io.Copy(f, resp.Body)
    if err != nil { panic(err) }
    println("written: ", written)
}