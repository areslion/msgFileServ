package util
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"runtime"
)

const (
	cst_cfgPath_linux = "/wsp/tsms/cfg/tsms_run.json"
	cst_cfgPath_windows = ".\\tsms_run.json"
)
var m_runcfg SxCfgAll

type sxCfg_serF struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
	Sep  string `json:"sep"`
	Path string `json:"path"`
	LogA string `json:"logA"`
	LogB string `json:"logB"`
	LogM int `json:"logM"`
	LogLev int `json:"logLevel"`
	LogObj int `json:"logObj"`
}
func (p *sxCfg_serF) GetDownloadUlrPre() (r_pre string) {
	r_pre = "http://" + p.Ip + ":" + p.Port + "/"
	return
}
type SxCfg_db struct{
	Ip string `json:"dbip"`
	Usr string `json:"usr"`
	Pwd string `json:"pwd"`
	Dbname string `json:"dbname"`
	Charset string `json:"charset"`
}
func (p* SxCfg_db)GetCntStr()string{
	cntstr := p.Usr+":"+p.Pwd+"@tcp("+p.Ip+")/"
	cntstr += p.Dbname+"?charset="+p.Charset

	return cntstr
}
type SxCfgAll struct {
	ServFile sxCfg_serF `json:"depotSft"`
	Db SxCfg_db `json:"runcfg"`
}

func init(){
	getRunCfg()
}

func getRunCfg() {
	log.Println(runtime.GOOS)

	cfgPath := getCfgPath()
	bts, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		log.Println("Fail to ReadFile " + cfgPath + " " + err.Error())
		return
	}

	err = json.Unmarshal(bts, &m_runcfg)
	if err != nil {
		log.Println("Fail to Unmarshal " + cfgPath + " " + err.Error())
		return
	}

	log.Println("run cfg ", m_runcfg)
	var px *sxCfg_serF  = &m_runcfg.ServFile
	InitLog(px.LogA,px.LogB,px.LogLev,px.LogObj,px.LogM)
}

func getCfgPath()(r_path string){
	switch runtime.GOOS {
	case "windows":
		r_path = cst_cfgPath_windows
	case "linux":
		r_path = cst_cfgPath_linux
		
	default:
		r_path = cst_cfgPath_linux		
	}

	return
}

func GetSftCfg()(r_sftCfg *SxCfgAll){ 
	r_sftCfg = & m_runcfg

	return
}