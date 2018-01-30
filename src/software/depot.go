package software

import (
	"bytes"
	"container/list"
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
)

var CfgSft *util.SxCfgAll = util.GetSftCfg()


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

func (p *SxSoft) SetUlrF(t_cfg *util.SxCfgAll, t_f string) string {
	p.Pathx = t_cfg.ServFile.GetDownloadUlrPre("download") + p.FolderID + util.Cst_sept + t_f
	return p.Pathx
}
func (p *SxSoft) Msgx() string {
	strRet := "SxSoft(" + fmt.Sprint(p.FlgSft) + " " + p.Namexf + " " + p.Namexa + " " + string(p.FlgSft) + " " + p.Ver + " " + p.Pathx + ")"
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
		log.Println(valx)
		intval, _ := strconv.Atoi(valx)
		p.FlgSft = uint(intval)
	} else if t_part.FormName() == cst_6des {
		p.Desc = valx
	} else if t_part.FormName() == cst_7md5 {
		p.Md5x = valx
	} else {
		util.L4E("SxSoft  undefed part " + t_part.FormName())
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
func (p *SxSoft) GetFolderPath(t_cfg *util.SxCfgAll, t_endsep bool) (r_folderPath string) { //获取文件仓库文件夹的路径
	folder := t_cfg.ServFile.PathSft + t_cfg.ServFile.Sep + p.FolderID
	if t_endsep == true {
		folder += t_cfg.ServFile.Sep
	}

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

func (p *SxSftDel) Mx() (r_msg string) {
	return "SxSftDel(" + p.NamexA + " " + p.Md5x + ")"
}

var (
	M_dbCfg *util.SxCfgAll
)

func init() {
	M_dbCfg = util.GetSftCfg()
}

func InsertDB(sft *SxSoft, cfg *util.SxCfgAll) (r_folderId string, b_ret bool) {
	dbopt,bret := dbbase.NewSxDB(&cfg.Db,""); if !bret {return}
	defer dbopt.Close()

	sft.FolderID,bret = getFolderID(sft.Namexf); if !bret {return}
	sft.SetUlrF(CfgSft, sft.Namexf)
	r_folderId = sft.FolderID
	dbopt.Sqlcmd = "REPLACE INTO depotSft(namexf,namexa,ver,pathx,pathIcon,flagSft,md5x,folderId,descx) "
	dbopt.Sqlcmd += "VALUES(?,?,?,?,?,?,?,?,?)"

	if !dbopt.Exc(sft.Namexf, sft.Namexa, sft.Ver, sft.Pathx, sft.PathIcon, sft.FlgSft, sft.Md5x, sft.FolderID, sft.Desc){return}

	b_ret = true
	return
}

func getFolderID(t_fileNmae string) (r_id string,b_ret bool) {
	dbopt,bret := dbbase.NewSxDB(&M_dbCfg.Db,"getFolderID"); if !bret {return}
	defer dbopt.Close()

	dbopt.Sqlcmd = "SELECT folderId FROM depotSft WHERE namexf = ?"
	if !dbopt.Query(t_fileNmae) {return}

	if dbopt.Rows.Next() {
		dbopt.Rows.Scan(&r_id)
	} else {
		r_id = util.Guid()
	}

	b_ret = true
	return
}

func GetSft(t_name string) (r_sft *SxSoft, r_folderid string, b_ret bool) {
	dbopt,bret := dbbase.NewSxDB(&M_dbCfg.Db,"GetSft");if !bret {return}
	defer dbopt.Close()

	r_sft = new(SxSoft)
	dbopt.Sqlcmd = "SELECT namexf,namexa,ver,pathx,pathIcon,flagSft,md5x,folderId "
	dbopt.Sqlcmd += "FROM depotSft "
	dbopt.Sqlcmd += "WHERE namexa = ? "
	if !dbopt.Query(t_name) {return}
	if dbopt.Rows.Next() {
		dbopt.Rows.Scan(&r_sft.Namexf, &r_sft.Namexa, &r_sft.Ver, &r_sft.Pathx, &r_sft.PathIcon, &r_sft.FlgSft, &r_sft.Md5x, &r_sft.FolderID)
		r_folderid = r_sft.FolderID
	} else {
		r_folderid = util.Guid()
		return
	}

	b_ret = true
	return
}

func DelSft(t_sft *SxSoft) (b_ret bool) {
	dbopt,bret := dbbase.NewSxDB(&M_dbCfg.Db,"DelSft"); if !bret {return}
	defer dbopt.Close()

	dbopt.Sqlcmd = "DELETE FROM depotSft WHERE namexa = ? "
	if !dbopt.Exc(t_sft.Namexa){return}

	b_ret = true
	util.L3I("software %s %s deleted",t_sft.Namexa,t_sft.FolderID)
	return true
}

func GetSftLst() (r_lst *list.List, r_json string, b_ret bool) {
	dbopt,bret := dbbase.NewSxDB(&M_dbCfg.Db,"GetSftLst"); if !bret {return}
	defer dbopt.Close()

	dbopt.Sqlcmd = "SELECT namexf,namexa,ver,pathx,pathIcon,flagSft,md5x,folderId,descx FROM depotSft "
	if !dbopt.Query(){return}

	r_lst = list.New()
	var lstar []SxSftListEle
	for dbopt.Rows.Next() {
		var sft SxSoft
		var jxe SxSftListEle
		dbopt.Rows.Scan(&sft.Namexf, &sft.Namexa, &sft.Ver, &sft.Pathx,
			&sft.PathIcon, &sft.FlgSft, &sft.Md5x, &sft.FolderID, &sft.Desc)
		r_lst.PushBack(sft)

		jxe.NamexF = sft.Namexf
		jxe.NamexA = sft.Namexa
		jxe.Ver = sft.Ver
		jxe.UlrF = sft.Pathx
		jxe.UlrIcon = sft.PathIcon
		jxe.Typex = fmt.Sprint(sft.FlgSft)
		jxe.Descx = sft.Desc
		lstar = append(lstar, jxe)
	}

	jx, _ := json.Marshal(lstar)
	r_json = "{\"repoAppList\":" + string(jx) + "}"

	util.L5F("GetSft num=" + fmt.Sprintf("%v", len(lstar)))
	b_ret = true
	return
}

func strx(t_ii string) string {
	var ret string
	ret = "'" + t_ii + "'"
	return ret
}
