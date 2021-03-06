package main

import (
	"flag"
"bytes"
"fmt"
"io"
"mime/multipart"
"io/ioutil"
"net/http"
"os"
"log"
"strconv"
)

const (
    // CstAddr = "http://10.20.10.101:1234/"
    // CstDownload = CstAddr+"download/"

    CstAddr = "http://10.20.11.17:1234/"
    CstDownload = CstAddr+"download/"

    cstmsg = CstAddr+"msgfile/"
    cstmsgnew = cstmsg+"newmsg"
)

// sample usage
func main() {
    nsel := flag.Int("sel",0,"choice functon")

    flag.Parse()
    log.Println("choice sel=",*nsel)

    switch *nsel {
    case 0:
        tstPostFile()
    case 1:
        delSft()
    case 2:
        getlstAPP()
    case 3:
        getFile("go1.9.2.windows-amd64.msi","7d68d66c-e201-4aea-9a2f-7a677e984ed9/go1.9.2.windows-amd64.msi")
    case 20:
        postfileMsg()
    default:
        panic("undefine parameter")
    }
    
}

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


func postMutiFile(filename string, targetUrl string,path string) error {
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

    fileWriter2, err := bodyWriter.CreateFormFile("file", "XX_"+filename)
    fh2, err := os.Open(path+filename)
    _, err = io.Copy(fileWriter2, fh2)

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

func newFileFormx(t_formname,t_filename ,t_path string,wter *multipart.Writer){
    fileWriter2, _ := wter.CreateFormFile(t_formname, t_filename)
    fh2, _ := os.Open(t_path)
    _, _ = io.Copy(fileWriter2, fh2)
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


    // fmName,_ :=bodyWriter.CreateFormFile("desc","descx")
    // strfmt :=fmx.Encode()
    // fmName.Write([]byte(strfmt))    
    // fmName,_ =bodyWriter.CreateFormFile("appname","visx")
    // fmName.Write([]byte("visx-------------1"))
    newForm(bodyWriter,"appname","","tst1")
    newForm(bodyWriter,"appversion","","1.9.1")
    newForm(bodyWriter,"appname","","tst1")
    newForm(bodyWriter,"apptype","","7")
    newForm(bodyWriter,"appdescription","","This is a form test 1.0.1")
    newForm(bodyWriter,"md5","","123abc")

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

func newForm(t_wter *multipart.Writer,t_filed,t_name,t_cxt string){
    ioW,_ := t_wter.CreateFormFile(t_filed,t_name)
    ioW.Write([]byte(t_cxt))
}

func postMutiForm(){
    //var mutiRes bytes.Buffer
    //muti_writer :=multipart.NewWriter(&mutiRes)

    
}

func postfileMsg(){
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)

    btsDesc ,_:= ioutil.ReadFile(".\\cfg\\tsms_msg.json")
    newForm(bodyWriter,"description","",string(btsDesc))
    newFileFormx("attachment","a.txt",".\\cfg\\tsms_msg.json",bodyWriter)
    newFileFormx("attachment","b.bat",".\\cfg\\tsms_msg.json",bodyWriter)
    newFileFormx("attachment","c.docx",".\\cfg\\tsms_msg.json",bodyWriter)

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()
    resp, err := http.Post(cstmsgnew, contentType, bodyBuf)
    //resp, err := http.Post(cstmsgnew, "", nil)
    if err != nil {
        return
    }
    defer resp.Body.Close()
    
    resp_body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
    }
    fmt.Println(resp.Status)
    fmt.Println(string(resp_body))
    return
}


func tstPostFile(){
    target_url := CstAddr+"uploadx"
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

func getFile(local,t_file string){

    f, err := os.OpenFile(local, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

func getlstAPP(){
    response,_ := http.Get(CstAddr+"getlstapp")
    defer response.Body.Close()
    body,_ := ioutil.ReadAll(response.Body)
    fmt.Println(string(body))
}

func delSft(){
    tmp := `{"appName":"tst1","appMd5":"123abc"}`
    req := bytes.NewBuffer([]byte(tmp))

    body_type := "application/json;charset=utf-8"    
    resp, _ := http.Post(CstAddr+"delsoft", body_type, req)
    //http.NewRequest("POST", CstDownload+"delsoft/", req_new)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))

    resp.Body.Close()
}