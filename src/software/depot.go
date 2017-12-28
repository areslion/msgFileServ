package software
import (
	"log"
	"dbbase"	
)


type SxSoft struct{
	namex string
	ver string
	pathx string
	flgSft uint32
	md5x string
}
func (p *SxSoft)msgx()(string){
	strRet := p.namex+"("+string(p.flgSft)+" "+p.ver+" "+p.pathx+")"
	return strRet
}
func (SxSoft)getKey()(r_key[] string){
	strLst := []string{"file","filename","appname","appversion","apptype","appdescription"}
	return strLst
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

	sqlcmd := "REPLACE INTO depotSft(namex,ver,pathx,flagSft,md5x) "
	sqlcmd +="VALUES('?','?','?',?,'?')"
	if _,err := cnt.Prepare(sqlcmd);err!=nil{
		logx("InsertDB  fail to Prepare "+err.Error())
		return false
	}
	if _,err := cnt.Exec(sft.namex,sft.ver,sft.pathx,sft.flgSft,sft.md5x);err!=nil{
		logx("saveNewSft  "+err.Error())
		return false
	}

	
	return true
}

func logx(t_msg string){
	log.Println("depot  "+t_msg)
}
