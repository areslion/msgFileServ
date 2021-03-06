package servFile

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)
import (
	"fmt"
	"os"
	"software"
	"util"
)

func delApp(t_res http.ResponseWriter, t_ask *http.Request) {
	var bDel = false
	var nret int
	log.Println("DelApp called")
	log.Println(t_ask)
	if t_ask.Method == "POST" {
		bts, err := ioutil.ReadAll(t_ask.Body)
		if err == nil {
			var sftDel software.SxSftDel
			err = json.Unmarshal(bts, &sftDel)
			if err == nil {
				log.Println(sftDel.Mx() + " will be removed")
				sft, _, bret := software.GetSft(sftDel.NamexA)
				if bret {
					err = os.RemoveAll(sft.GetFolderPath(software.CfgSft, true))
					if err == nil {
						log.Println("remove folder ", sft.GetFolderPath(software.CfgSft, true))
						bDel = software.DelSft(sft)
						if bDel {
							util.L3I(sftDel.Mx() + " removed successfully")
						} else {
							util.L4E(sftDel.Mx() + " removed faild")
						}

					} else {
						log.Println("Fail to remove folder ", sft.GetFolderPath(software.CfgSft, true)+" "+err.Error())
					}
				} else {
					log.Println(sftDel.Mx() + " is not exist in server")
				}
			} else {
				util.L5F("Fail to parse json " + err.Error() + "  " + string(bts))
				var sftx software.SxSftDel
				sftx.NamexA = "tst1"
				sftx.Md5x = "123abc"

				jx, _ := json.Marshal(sftx)
				util.L3I("C---" + string(bts))
				util.L3I("S---" + string(jx))
			}

			if bDel == true {
				nret = http.StatusFound
			} else {
				nret = http.StatusInternalServerError
			}

			util.L3I("DelApp res=" + fmt.Sprintf("%d", nret))
			//http.Redirect(t_res,t_ask,"./View?id=",nret)
		} else {
			util.L4E("fail to read body data  " + err.Error())
		}
	} else {
		util.L4E("DelApp  undefined method=" + t_ask.Method)
	}
}

func downFileHandler(t_res http.ResponseWriter, t_ask *http.Request) {
	util.L3I("downFileHandler %s",t_ask.Method)
	
	if t_ask.Method=="GET" {util.NewFileServ(t_ask,&t_res,software.CfgSft.ServFile.PathSft)}
}

func getlstApp(t_res http.ResponseWriter, t_ask *http.Request) {
	util.L3I("GetlstApp called")

	if t_ask.Method == "GET" {
		_, strJson, _ := software.GetSftLst()
		t_res.Header().Set("Content-Type", "application/json; charset=utf-8")
		t_res.Write([]byte(strJson))
	}
}


func parseAsk(t_ask *http.Request) bool {
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

	util.L3I("Upload sft " + sft.Msgx())
	sft.FolderID, bfileSave = software.InsertDB(&sft, software.M_dbCfg)
	if buf.Len() > 0 {
		if bfileSave {
			_, bfileSave = saveFileBytes(&sft, buf.Bytes())
			if bfileSave {
				util.L3I("Upload sft " + sft.Msgx() + " successfully")
			} else {
				util.L4E("Upload sft " + sft.Msgx() + " failed")
			}

		} else {
			util.L5F("Upload sft " + sft.Msgx() + " Insert into db failed")
		}
	}

	return bfileSave
}

func saveFileBytes(t_f *software.SxSoft, buf []byte) (r_path string, b_ret bool) {
	folder := t_f.GetFolderPath(software.CfgSft, true)
	fileServer := folder + t_f.Namexf

	os.MkdirAll(folder, 0711)
	return util.SaveFileBytes(fileServer, buf)
}

func StarFileServ() {
	http.HandleFunc("/uploadx", upload)            //POST upload software
	http.HandleFunc("/download/", downFileHandler) //GET download file
	http.HandleFunc("/getlstapp", getlstApp)       //GET app list
	http.HandleFunc("/delsoft", delApp)            //POST delete software

	//err := http.ListenAndServe(":1234", nil)
	err := http.ListenAndServe(":"+software.CfgSft.ServFile.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
		fmt.Println("ListenAndServe 启动服务器失败 ", err.Error())
	} else {
		log.Fatal("ListenAndServe 重启动服务器")
		fmt.Println("ListenAndServe 重启动服务器")
	}
}

func upload(t_res http.ResponseWriter, t_ask *http.Request) {
	util.L3I(t_ask.Method)
	if t_ask.Method == "POST" {

		var nret int
		if bret := parseAsk(t_ask); bret == true {
			nret = http.StatusFound
		} else {
			nret = http.StatusInternalServerError
		}

		util.L3I("upload res=" + fmt.Sprintf("%d", nret))
		//http.Redirect(t_res, t_ask, "./View?id=", nret)
	}
}
