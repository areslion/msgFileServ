package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"sync"
)
import (
	"httpserv"
	"software"
	"util"
)

const (
	//CSTUpdate_dir = "/wsp/gotst/upload"
	tpl = `<html>  
	<head>  
	<title>上传文件</title>  
	</head>  
	<body>  
	<form enctype="multipart/form-data" action="/upload" method="post">  
	<input type="file" name="uploadfile" />  
	<input type="hidden" name="token" value="{...{.}...}"/>  
	<input type="submit" value="upload" />  
	</form>  
	</body>  
	</html>`

	tpl2 = `<html>  
	<head>  
	<title>上传文件</title>  
	</head>  
	<body>  
	<form enctype="multipart/form-data" action="/uploadx" method="post">  
	<!--input type="file" name="file" /-->  
	<input type="file" name="file"/>  
	<input type="submit" value="upload" />  
	</form>  
	</body>  
	</html>`
)

var (
	wg sync.WaitGroup
)

func init() {
	software.M_dbCfg.Init("10.20.10.101", "root", "123456", "deskSafe", "utf8")
}

func main() {
	//inithttp()
	//tstdatabase()
	mutiRun()
}

func mutiRun() {
	runtime.GOMAXPROCS(1)
	wg.Add(2)

	go inithttp()
	go tstdatabase()

	wg.Wait()
}

func tstdatabase() {
	// dbbase.Tstmysql()
	// lstSft, strJson, _ := software.GetSftLst()
	// for ix := lstSft.Front(); ix != nil; ix = ix.Next() {
	// 	sft := ix.Value.(software.SxSoft)
	// 	log.Println(sft.Msgx())
	// }
	// log.Println(strJson)
	wg.Done()
}

func inithttp() {
	http.HandleFunc("/", index)
	http.HandleFunc("/view", ViewHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/uploadx", httpserv.Upload)            //POST upload software
	http.HandleFunc("/download/", httpserv.DownFileHandler) //GET download file
	http.HandleFunc("/getlstapp", httpserv.GetlstApp)       //GET app list
	http.HandleFunc("/delsoft", httpserv.DelApp)            //POST delete software

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
		fmt.Println("ListenAndServe 启动服务器失败 ", err.Error())
	} else {
		log.Fatal("ListenAndServe 重启动服务器")
		fmt.Println("ListenAndServe 重启动服务器")
	}

	wg.Done()
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(tpl2))
}

func showPart(t_part *multipart.Part) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(t_part)
	log.Println("\r\n\r\nForma name=", t_part.FormName(),
		"\r\n", "File name=", t_part.FileName(), "\r\n", "size=", buf.Len(), "\r\n", "content=", buf.String())
}

func ViewHandler(t_res http.ResponseWriter, t_ask *http.Request) {
	imageid := t_ask.FormValue("id")
	imagepath := software.CSTUpdate_dir + "/" + imageid
	if bExist := util.IsExists(imagepath); !bExist {
		http.NotFound(t_res, t_ask)
		return
	}
}

func helloHandler(t_res http.ResponseWriter, t_ask *http.Request) {
	log.Println(t_ask)
	body, _ := ioutil.ReadAll(t_ask.Body)

	log.Println(string(body))
	log.Println(t_ask.RequestURI)
	log.Println(t_ask.URL.Path)

	if t_ask.Method == "GET" {
		log.Println("hello------", 2)
		t, err := template.ParseFiles("." + software.CSTPathSep + "html" + software.CSTPathSep + "hello.html")
		if err != nil {
			http.Error(t_res, err.Error(), http.StatusInternalServerError)
			log.Println(err.Error(), http.StatusInternalServerError)
			log.Fatal("upload", err.Error())
			log.Println("hello------", 3)
			return
		}

		log.Println("hello------", 4)
		t.Execute(t_res, nil)
		log.Println("hello------", 5)
		return
	}
}

func tstDownload() {
	res, err := http.Get("http://localhost:1234/test/client.cpp")
	if err != nil {
		panic(err)
	}
	f, err := os.Create("qq.exe")
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body)
}

func logx(t_msg string) {
	log.Println("serHttpFile  " + t_msg)
}
