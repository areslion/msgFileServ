package httpserv
import (
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
)
import (
	"os"
	"software"
)

var (
	staticHandler http.Handler
)

func DelApp(t_res http.ResponseWriter,t_ask *http.Request){
	var bDel = false
	var nret int
	log.Println("DelApp called")
	if t_ask.Method == "POST" {
		bts,err := ioutil.ReadAll(t_ask.Body)
			var sftDel  software.SxSftDel
			err = json.Unmarshal(bts,&sftDel)
			if err ==nil {
				log.Println(sftDel.Mx()+" will be removed")
				sft ,_,bret:= software.GetSft(sftDel.NamexA)
				if bret {
					err = os.RemoveAll(sft.GetFolderPath(true))
					if err ==nil {
						log.Println("remove folder ",sft.GetFolderPath(true))
						bDel = software.DelSft(sft)
						if bDel {
							logx(sftDel.Mx()+" removed successfully")
						} else {
							logx(sftDel.Mx()+" removed faild")
						}

					} else {
						log.Println("Fail to remove folder ",sft.GetFolderPath(true)+" "+err.Error())
					}
				} else {
					log.Println(sftDel.Mx()+" is not exist in server")
				}
			} else {
				logx("Fail to parse json "+err.Error()+"  "+string(bts))
				var sftx software.SxSftDel
				sftx.NamexA = "tst1"
				sftx.Md5x = "123abc"

				jx,_ := json.Marshal(sftx)
				logx("C---"+string(bts))
				logx("S---"+string(jx))
			}

			if bDel==true {
				nret = http.StatusFound
			} else {
				nret = http.StatusInternalServerError
			}
			
			http.Redirect(t_res,t_ask,"./View?id=",nret)
		}
}

func DownFileHandler(t_res http.ResponseWriter,t_ask *http.Request){ 
	log.Println("path:" + t_ask.URL.Path)
    staticHandler.ServeHTTP(t_res, t_ask)
}


func GetlstApp(t_res http.ResponseWriter,t_ask *http.Request){
	logx("GetlstApp called")

	if t_ask.Method =="GET" {
		_,strJson,_ := software.GetSftLst()
		t_res.Write([]byte(strJson))
	}
}



func logx(t_msg string){
	log.Println("fileserv  ",t_msg)
}