package software
import (
	"log"
	"mime/multipart"
	"bytes"
	"strconv"
)
import (
	"dbbase"	
)


const(
	//"file","filename","appname","appversion","apptype","appdescription"
	cst_1file string = "file"
	cst_2fnam string ="filename"
	cst_3anam string ="appname"
	cst_4ver string ="appversion"
	cst_5type string ="apptype"
	cst_6des string ="appdescription"
	cst_7md5 string ="md5"
)
type SxSoft struct{
	namexf string//file name 
	namexa string//app name
	ver string
	pathx string
	desc string//description
	md5x string
	flgSft uint	
}
func (p *SxSoft)Msgx()(string){
	strRet := p.namexf+"("+" "+p.namexa+" "+string(p.flgSft)+" "+p.ver+" "+p.pathx+")"
	return strRet
}
func (SxSoft)getKey()(r_key[] string){
	strLst := []string{cst_1file,cst_2fnam,cst_3anam,cst_4ver,cst_5type,cst_6des}
	return strLst
}
func (p* SxSoft)Set(t_part *multipart.Part)bool{
	buf := new(bytes.Buffer)
	buf.ReadFrom(t_part)
	valx := buf.String()

	if t_part.FormName() == cst_2fnam {
		p.namexf = valx
	} else if t_part.FormName() == cst_3anam {
		p.namexa = valx
	} else if t_part.FormName() == cst_4ver {
		p.ver = valx
	} else if t_part.FormName() == cst_5type {
		intval,_ := strconv.Atoi(valx)
		p.flgSft = uint(intval)
	} else if t_part.FormName() == cst_6des {
		p.desc = valx
	} else if t_part.FormName() == cst_7md5 {
		p.md5x = valx
	} else {
		logx("SxSoft  undefed part "+t_part.FormName())
		return false
	}
	return true
}
func (p* SxSoft)SetNameFile(t_name string){
	p.namexf = t_name
}
func (p* SxSoft)SetPath(t_x string){
	p.pathx = t_x
}


var(
	M_dbCfg dbbase.SCfg
)

func InsertDB(sft *SxSoft,cfg *dbbase.SCfg) bool {
	cnt,bret :=dbbase.Open(cfg)
	if bret==false{
		return false
	}
	defer dbbase.Close()
	defer cnt.Close()

	sqlcmd := "REPLACE INTO depotSft(namexf,namexa,ver,pathx,flagSft,md5x) "
	sqlcmd +="VALUES(?,?,?,?,?,?)"
	smt,err := cnt.Prepare(sqlcmd)
	if err!=nil{
		logx("InsertDB  fail to Prepare "+err.Error())
		return false
	}
	if _,err := smt.Exec(sft.namexf,sft.namexa,sft.ver,sft.pathx,sft.flgSft,sft.md5x);err!=nil{
		logx("saveNewSft  "+err.Error())
		return false
	}

	return true
}

func strx(t_ii string) (string){
	var ret string
	ret = "'"+t_ii+"'"
	return ret
}

func logx(t_msg string){
	log.Println("depot  "+t_msg)
}
