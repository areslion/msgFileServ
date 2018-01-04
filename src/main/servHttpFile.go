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
	"dbbase"
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
	staticHandler http.Handler
	wg            sync.WaitGroup
)

func init() {
	staticHandler = http.StripPrefix("/download/", http.FileServer(http.Dir("download")))
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
	dbbase.Tstmysql()
	lstSft, strJson, _ := software.GetSftLst()
	for ix := lstSft.Front(); ix != nil; ix = ix.Next() {
		sft := ix.Value.(software.SxSoft)
		log.Println(sft.Msgx())
	}
	log.Println(strJson)
	wg.Done()
}

func inithttp() {
	http.HandleFunc("/", index)
	http.HandleFunc("/view", ViewHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/uploadx", upload)                //POST upload software
	http.HandleFunc("/download/", downFileHandlerEx)     //GET download file
	http.HandleFunc("/getlstapp", httpserv.GetlstApp) //GET app list
	http.HandleFunc("/delsoft", httpserv.DelApp)      //POST delete software

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

func uploadEz(t_res http.ResponseWriter, t_ask *http.Request) {
	if t_ask.Method == "GET" {
		t, err := template.ParseFiles("/upload.html")
		if err != nil {
			http.Error(t_res, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error(), http.StatusInternalServerError)
			log.Fatal("upload", err.Error())
			return
		}

		t.Execute(t_res, nil)
		return
	}
	if t_ask.Method == "POST" {

		//ParseAskx(t_ask)

		f, h, err := t_ask.FormFile("file")
		//fmt.Println(t_ask,"\r\nPostForm\r\n",t_ask.PostForm,"\r\n","\r\nForm\r\n",t_ask.Form)
		if err != nil {
			http.Error(t_res, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error(), http.StatusInternalServerError)
			log.Fatal("upload", err.Error())
			http.Redirect(t_res, t_ask, "./View?id="+"", http.StatusFound)
			return
		}

		filename := h.Filename
		defer f.Close()

		log.Println("接到上传文件请求 name=", h.Filename, "size=", h.Size)
		//open files
		fileServer := software.CSTUpdate_dir + software.CSTPathSep + filename
		t, err := os.Create(fileServer)
		if err != nil {
			http.Error(t_res, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error(), http.StatusInternalServerError)
			log.Fatal("upload 创建服务端文件失败", err.Error())
			http.Redirect(t_res, t_ask, "./View?id="+filename, http.StatusFound)
			return
		}
		defer t.Close()

		log.Println("成功创本地文件 ", fileServer)
		if _, err := io.Copy(t, f); err != nil {
			http.Error(t_res, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error(), http.StatusInternalServerError)
			log.Fatal("upload", err.Error())
			http.Redirect(t_res, t_ask, "./View?id="+filename, http.StatusFound)
			return
		}

		log.Println("文件 ", fileServer, " 上传成功，将跳转到文件浏览页面")
		//http.Redirect(t_res, t_ask, "./View?id="+filename, http.StatusFound)
	}
}

func upload(t_res http.ResponseWriter, t_ask *http.Request) {
	if t_ask.Method == "POST" {

		var nret int
		if bret := ParseAsk(t_ask); bret == true {
			nret = http.StatusFound
		} else {
			nret = http.StatusInternalServerError
		}

		logx("upload res="+fmt.Sprintf("%d",nret))
		//http.Redirect(t_res, t_ask, "./View?id=", nret)
	}
}

func ParseAsk(t_ask *http.Request) bool {
	muti_reader, _err := t_ask.MultipartReader()
	var sft software.SxSoft
	var bfileSave bool = false
	buf := new(bytes.Buffer)

	if _err == nil {
		for {
			part, err := muti_reader.NextPart()
			if err == io.EOF {
				break
			}

			if part.FormName() == "file" {
				buf.ReadFrom(part)
				sft.SetNameFile(part.FileName())
			} else {
				sft.Set(part)
			}

		}
	}

	if buf.Len() > 0 {
		logx("Upload sft " + sft.Msgx())
		sft.FolderID, bfileSave = software.InsertDB(&sft, &software.M_dbCfg)
		if bfileSave {
			_, bfileSave = saveFileBytes(&sft, buf.Bytes())
			if bfileSave {
				logx("Upload sft " + sft.Msgx() + " successfully")
			} else {
				logx("Upload sft " + sft.Msgx() + " failed")
			}

		} else {
			logx("Upload sft " + sft.Msgx() + " Insert into db failed")
		}
	}

	return bfileSave
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
	if bExist := isExists(imagepath); !bExist {
		http.NotFound(t_res, t_ask)
		return
	}
}

func uploadEnd(t_res http.ResponseWriter, t_ask *http.Request) {
	return
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return os.IsExist(err)
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

func downFileHandler(t_res http.ResponseWriter, t_ask *http.Request) {
	log.Println("path:" + t_ask.URL.Path)
	_,_,_ = util.GetPathEle(t_ask.URL.Path)
	staticHandler.ServeHTTP(t_res, t_ask)
}

func downFileHandlerEx(t_res http.ResponseWriter, t_ask *http.Request) {
	log.Println("path:" + t_ask.URL.Path)

	var sft software.SxSoft
	folder,fldID,_ := util.GetPathEle(t_ask.URL.Path)
	prefix := "/"+folder+"/"
	sft.FolderID = fldID
	logx("start fileservr("+prefix+" "+sft.GetFolderPath(false)+")")
	staticFServ := http.StripPrefix(prefix, http.FileServer(http.Dir(sft.GetFolderPath(false))))
	staticFServ.ServeHTTP(t_res, t_ask)
}

func saveFile(filename string, file *multipart.File) bool {
	fileServer := software.CSTUpdate_dir + software.CSTPathSep + filename
	t, err := os.Create(fileServer)
	if err != nil {
		log.Println("saveFile Fail to create file "+fileServer+" ", err.Error())
		return false
	}
	defer t.Close()

	if _, err := io.Copy(t, *file); err != nil {
		log.Println("saveFile fail to save file " + fileServer + " " + err.Error())
		return false
	}

	return true
}

func saveFileBytes(t_f *software.SxSoft, buf []byte) (r_path string, b_ret bool) {
	folder := t_f.GetFolderPath(true)
	fileServer := folder + t_f.Namexf

	os.MkdirAll(folder, 0711)
	return util.SaveFileBytes(fileServer, buf)
}

func logx(t_msg string) {
	log.Println("serHttpFile  " + t_msg)
}
