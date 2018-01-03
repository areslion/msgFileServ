package software

import (
	"bytes"
	"container/list"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)
import (
	"dbbase"
	"util"
)

const (
	//"file","filename","appname","appversion","apptype","appdescription"
	cst_1file string = "file"
	cst_2fnam string = "filename"
	cst_3anam string = "appname"
	cst_4ver  string = "appversion"
	cst_5type string = "apptype"
	cst_6des  string = "appdescription"
	cst_7md5  string = "md5"

	CstAddr     = "http://localhost:1234/"
	CstDownload = "http://localhost:1234/download/"

	CSTPathSep = "\\"
	CSTUpdate_dir = "E:"+CSTPathSep+"workspace"+CSTPathSep+"005.XNKJ"+CSTPathSep+"tstFilServ"
)

type SxSoft struct {
	Namexf   string //file name
	Namexa   string //app name
	Ver      string
	Pathx    string
	PathIcon string //icon path
	Desc     string //description
	Md5x     string
	FlgSft   uint
	FolderID string //folder id
}

func (p *SxSoft) SetUlrF(t_f string) string {
	p.Pathx = CstDownload + p.FolderID + util.Cst_sept + t_f
	return p.Pathx
}
func (p *SxSoft) Msgx() string {
	strRet :=  "SxSoft("+p.Namexf + " " + p.Namexa + " " + string(p.FlgSft) + " " + p.Ver + " " + p.Pathx + ")"
	return strRet
}
func (SxSoft) getKey() (r_key []string) {
	strLst := []string{cst_1file, cst_2fnam, cst_3anam, cst_4ver, cst_5type, cst_6des}
	return strLst
}
func (p *SxSoft) Set(t_part *multipart.Part) bool {
	buf := new(bytes.Buffer)
	buf.ReadFrom(t_part)
	valx := buf.String()

	if t_part.FormName() == cst_2fnam {
		p.Namexf = valx
	} else if t_part.FormName() == cst_3anam {
		p.Namexa = valx
	} else if t_part.FormName() == cst_4ver {
		p.Ver = valx
	} else if t_part.FormName() == cst_5type {
		intval, _ := strconv.Atoi(valx)
		p.FlgSft = uint(intval)
	} else if t_part.FormName() == cst_6des {
		p.Desc = valx
	} else if t_part.FormName() == cst_7md5 {
		p.Md5x = valx
	} else {
		logx("SxSoft  undefed part " + t_part.FormName())
		return false
	}
	return true
}
func (p *SxSoft) SetNameFile(t_name string) {
	p.Namexf = t_name
}
func (p *SxSoft) SetPath(t_x string) {
	p.Pathx = t_x
}
func (p *SxSoft) GetFolderPath()(r_folderPath string){//获取文件仓库文件夹的路径
	folder := CSTUpdate_dir+CSTPathSep + p.FolderID + CSTPathSep

	return folder
}

//获取软件仓库中单个软件的基本信息
type SxSftListEle struct {
	NamexA  string `json:"appDisplayName"`
	NamexF  string `json:"appSetupPackageName"`
	Ver     string `json:"appVersion"`
	UlrF    string `json:"appDownLoadUrl"`
	UlrIcon string `json:"appDisplayIcon"`
	Typex   string `json:"appType"`
	Md5x    string `json:"appMd5"`
	Descx   string `json:"appDescription"`
	Sizex   string `json:"appSize"`
}

type SxSftDel struct {
	NamexA string `json:"appName"`
	Md5x   string `json:"appMd5"`
}
func (p *SxSftDel)Mx()(r_msg string){
	return "SxSftDel("+p.NamexA+" "+p.Md5x+")"
}

var (
	M_dbCfg dbbase.SCfg
)

func InsertDB(sft *SxSoft, cfg *dbbase.SCfg) (r_folderId string, b_ret bool) {
	cnt, bret := dbbase.Open(cfg)
	if bret == false {
		return "", false
	}
	defer dbbase.Close()
	defer cnt.Close()

	sft.FolderID = getFolderID(cnt, sft.Namexf)
	sft.SetUlrF(sft.Namexf)
	sqlcmd := "REPLACE INTO depotSft(namexf,namexa,ver,pathx,pathIcon,flagSft,md5x,folderId,descx) "
	sqlcmd += "VALUES(?,?,?,?,?,?,?,?,?)"
	smt, err := cnt.Prepare(sqlcmd)
	if err != nil {
		logx("InsertDB  fail to Prepare " + err.Error())
		return "", false
	}

	_, err = smt.Exec(sft.Namexf, sft.Namexa, sft.Ver, sft.Pathx, sft.PathIcon, sft.FlgSft, sft.Md5x, sft.FolderID, sft.Desc)
	if err != nil {
		logx("saveNewSft  " + err.Error())
		return "", false
	}

	return sft.FolderID, true
}

func getFolderID(t_db *sql.DB, t_fileNmae string) (r_id string) {
	var folderid string
	sqlcmd := "SELECT folderId FROM depotSft WHERE namexf = ?"
	//sqlcmd += strx(t_fileNmae)

	smt, err := t_db.Prepare(sqlcmd)
	if err != nil {
		logx("getFolderID  fail to Prepare " + err.Error())
		var folderid string
		return folderid
	}

	rows, err := smt.Query(t_fileNmae)
	if err != nil {
		logx("getFolderID  " + err.Error())
		return util.Guid()
	}

	if rows.Next() {
		rows.Scan(&folderid)
	} else {
		folderid = util.Guid()
	}

	return folderid
}

func GetSft(t_name string) (r_sft *SxSoft, r_folderid string, b_ret bool) {
	var folderid string
	sft := new(SxSoft)
	sqlcmd := "SELECT namexf,namexa,ver,pathx,pathIcon,flagSft,md5x,folderId "
	sqlcmd += "FROM depotSft "
	sqlcmd += "WHERE namexa = ? "

	cnt, bret := dbbase.Open(&M_dbCfg)
	if bret == false {
		return nil, "", false
	}
	defer dbbase.Close()
	defer cnt.Close()

	smt, err := cnt.Prepare(sqlcmd)
	if err != nil {
		logx("GetSft  fail to Prepare " + err.Error())
		return nil, "", false
	}

	rows, err := smt.Query(t_name)
	if err != nil {
		logx("GetSft  " + err.Error())
		return nil, "", false
	}

	if rows.Next() {
		rows.Scan(&sft.Namexf, &sft.Namexa, &sft.Ver, &sft.Pathx, &sft.PathIcon, &sft.FlgSft, &sft.Md5x, &sft.FolderID)
		folderid = sft.FolderID
	} else {
		folderid = util.Guid()
		return sft, folderid, false
	}

	return sft, folderid, true
}

func DelSft(t_sft *SxSoft) bool {
	cnt, bret := dbbase.Open(&M_dbCfg)
	if bret == false {
		return false
	}
	defer dbbase.Close()
	defer cnt.Close()

	sqlcmd := "DELETE FROM depotSft WHERE namexa = ? "
	smt, err := cnt.Prepare(sqlcmd)
	if err != nil {
		logx("DelSft  fail to Prepare " + err.Error())
		return false
	}

	_, err = smt.Exec(t_sft.Namexa)
	if err != nil {
		logx("DelSft  " + err.Error())
		return false
	}

	return true
}

func GetSftLst() (r_lst *list.List,r_json string, b_ret bool) {
	sqlcmd := "SELECT namexf,namexa,ver,pathx,pathIcon,flagSft,md5x,folderId,descx FROM depotSft "

	cnt, bret := dbbase.Open(&M_dbCfg)
	if bret == false {
		return nil,"", false
	}
	defer dbbase.Close()
	defer cnt.Close()

	smt, err := cnt.Prepare(sqlcmd)
	if err != nil {
		logx("GetSftLst  fail to Prepare " + err.Error())
		return nil,"", false
	}

	rows, err := smt.Query()
	if err != nil {
		logx("GetSft  " + err.Error())
		return nil, "",false
	}

	lstSft := list.New()
	var lstar [] SxSftListEle
	for rows.Next() {
		var sft SxSoft
		var jxe SxSftListEle
		rows.Scan(&sft.Namexf, &sft.Namexa, &sft.Ver, &sft.Pathx,
			&sft.PathIcon, &sft.FlgSft, &sft.Md5x, &sft.FolderID, &sft.Desc)
		lstSft.PushBack(sft)

		jxe.NamexF = sft.Namexf;jxe.NamexA = sft.Namexf;jxe.Ver=sft.Ver
		jxe.UlrF = sft.Pathx; jxe.UlrIcon = sft.PathIcon; jxe.Typex = fmt.Sprint(sft.FlgSft); jxe.Descx = sft.Desc
		//jx, _ := json.Marshal(jxe)

		lstar = append(lstar,jxe)
		//logx(string(jx))
	}

	jx, _ := json.Marshal(lstar)
	strRetJson := "{\"repoAppList\":"+string(jx)+"}"

	return lstSft, strRetJson,true
}

func strx(t_ii string) string {
	var ret string
	ret = "'" + t_ii + "'"
	return ret
}

func logx(t_msg string) {
	log.Println("depot  " + t_msg)
}
