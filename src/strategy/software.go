//功能：软件策略服务端支撑
package strategy

import (
	// "bytes"
	// "container/list"
	// "encoding/json"
	// "fmt"
	// "log"
	// "mime/multipart"
	// "strconv"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
)
import (	
	//"os"
	"dbbase"
	"util"
)


//软件策略中的软件基本属性
type sxSftItem struct{
	Namex string `json:"name"`//软件的名称
	Enable bool `json:"enable"`//在策略中是否被启用
}
type sxSftItemList struct{
	List []sxSftItem `json:"list"`
}

//软件策略应用于设备的基本信息
type sxDevItem struct{
	Namex string `json:"name"`//终端归属人的姓名
	Path string `json:"path"`//归属部门
	NumDev string `json:"numDev"`//终端ID
	NumUser string `json:"numUser"`//用户Key
}
type sxDevItemList struct{
	List []sxDevItem `json:"list"`
}
//软件策略返回结果
type sxRes struct{
	ListSft sxSftItemList `json:"listSoft"`//策略软件清单
	ListDev sxDevItemList `json:"listDev"`//策略应用于的终端清单
	Num string `json:"num"`//策略的ID
	Name string `json:"name"`//策略名称
	Cls int `json:"class"`//策略类型 1 为软件白名单 2为软件黑名单
}

//软件过滤器
type sxFilter struct{
	strFilter []string//需要要被过滤的软件关键字
	strKeeper string//放置重叠的关键字系列
	strPlacer []string//需要被替换的关键字
}

//func (p *sxEmp) GetLstDepat(t_path, t_sep string) (r_lst []string, r_json string) 
func (p * sxFilter) init(){	
	filter := []string{"_CRT_",
		"AMD ",
		"Catalyst","Cisco","CryptoKit",
		"HP ",
		"IntelliJ",
		"IntelliJ",
		"NVIDIA",
		"Runtime",
		"Support",
		"tools",
		"vs_",
		" (KB","Chinese (Simplified)","Framework","Visual C++ Compiler",
		"NET Framework",".NET","Update","for Visual Studio","Intel(R) Processor Graphics",
		"Microsoft Visual","Microsoft System","Microsoft Web ","Microsoft SQL Server ",
		"Visual S","Visual F","Visual C","Visual B",
		"Windows ",
		"icecap_collection_neutral","tools-freebsd","WinRT Intellisense","Kit SDK ","Windows Runtime ","语言包","Pack"," Help "," (SP",
		" Component","Java","amd64","Silverlight","Utility","Language","Library","Service","Adobe","Realtek","Sensor","MSI","Tools",
		"Realtek","(x86)","Toolbar","Intel(R)","icbc","ICBC","Oracle","MySQL","C++"}	
		p.strFilter = append(p.strFilter,filter...)

		xreplacer :=[]string{"0","1","2","3","4","5","6","7","8","9",".","-","、",
			"正式版","版本"}
		p.strPlacer = append(p.strPlacer,xreplacer...)


}
//查看输入串是否被过滤
func (p *sxFilter) filtered(t_name string)(b_ret bool){
		b_ret = false
		for _,itx := range p.strFilter{
			if strings.Contains(t_name,itx){
				//util.L3I("soft filtered "+itx)
				return true
			}
		}

	return;
}
//在输入串中如果存在指定的关键字则过滤
func (p *sxFilter) replace(t_str string)(r_str string){
		r_str = t_str;
		for _,itx:=range p.strPlacer {
			r_str = strings.Replace(r_str,itx,"",-1)
		}

	return
}
//先对输入的数据整理 然后过滤  过滤之后看是否已经被录入过
//如果 返回为true 则表明是有效数据需要被录入  否则直接丢弃即可
func (p *sxFilter) filterx(t_strIn string)(r_str string,b_filed bool){
	b_filed = false
	r_str = p.replace(t_strIn)
	b_filed = p.filtered(r_str); if b_filed {
		return
	}

	var stritem = strings.Replace(r_str," ","",-1)
			stritem = stritem+" "
			if strings.Index(p.strKeeper,stritem)==-1 {p.strKeeper = p.strKeeper + stritem+" "
			} else {b_filed = true
			}

	return
}




type StrategySoft struct{

}

//按照策略的编号获取策略文件所在的绝对路径
func (p *StrategySoft)getStrategyPath(t_num string)(r_path string){
	pcfg := util.GetSftCfg()
	r_path = pcfg.ServFile.PathSrategy + pcfg.ServFile.Sep + t_num + ".stra"
	return
}

//根据软件策略编号获取软件清单信息
//如果num为空 则从数据库获取软件清单 否则从磁盘的json文件获取策略的详情
func (p *StrategySoft)getSoftList(t_num string)(r_res sxRes,r_bts []byte){
	path := p.getStrategyPath(t_num)
	var err error
	r_bts,err = ioutil.ReadFile(path);if err==nil {
		util.L4E(path+" "+err.Error())
	}

	err = json.Unmarshal(r_bts,r_res);if err!=nil{
		util.L4E("json.Unmarshal(bts,r_res) "+err.Error())
		return
	}
	
	return
}
//从数据库获取清单信息
func (p *StrategySoft)getSoftListFromDB()(r_res sxRes,r_bts []byte){
	dbopt, bret := dbbase.NewSxDB(&util.GetSftCfg().Db, "read software list")
	if !bret {return}
	defer dbopt.Close()

	var filterx sxFilter
	filterx.init()

	//var stringx string
	dbopt.Sqlcmd = "SELECT DISTINCT namex FROM softInstalled ORDER BY namex "
	if !dbopt.Query() { return }
	for dbopt.Next() {
		var strName sql.NullString
		var ele sxSftItem		
		dbopt.Scan(&strName)
		strRet,bret := filterx.filterx(strName.String);if bret{ continue
		}
		ele.Namex = strRet
		r_res.ListSft.List = append(r_res.ListSft.List,ele)
	}

	var err error
	r_bts,err = json.Marshal(r_res);if err!=nil{
		util.L4E("json.Marshal(r_res) "+err.Error())
		return
	}

	util.L3I("softwaare num=%d", len(r_res.ListSft.List))
	return
}
//获取软件策略
func (p *StrategySoft)getStrategy(t_num string)(r_res sxRes,r_bts []byte,r_str string){
	if len(t_num)==0{
		r_res,r_bts = p.getSoftListFromDB()
	} else {
		r_res, r_bts = p.getSoftList(t_num)
	}

	r_str = string(r_bts)
	return
}
//存储软件策略
func (p *StrategySoft)saveStrategy(t_bts []byte)(b_ret bool){
	var strx sxRes
	err := json.Unmarshal(t_bts,strx);if err!=nil{//change byte to struct
		util.L4E("json.Unmarshal(t_bts,strx) "+err.Error())
	}

	if len(strx.Num)!=36 {
		util.L4E("invalid num(%s) length=%d need=%d",strx.Num,len(strx.Num),36)
		return
	}

	path := p.getStrategyPath(strx.Num)
	_, b_ret = util.SaveFileBytes(path, t_bts);if !b_ret {
		util.L4E("write software strategy(%s) failed(%s)",strx.Num,path)
	}
	
	b_ret = p.insertDB(strx)
	
	return
}

//将软件策略存储到数据库中的sftRuleAbstract和sftRuleSend两张表中
func (p *StrategySoft)insertDB(t_stra sxRes)(b_ret bool){
	dbopt, bret := dbbase.NewSxDB(&util.GetSftCfg().Db, "insert software startegy into db")
	if !bret {return}
	defer dbopt.Close()

	dbopt.Sqlcmd = " REPLACE INTO sftRuleAbstract(num,namex,cls) VALUES(?,?,?)"
	b_ret = dbopt.ExcAlone(t_stra.Num,t_stra.Name,t_stra.Cls);if !b_ret{return}
	util.L3I("%s replaced into sftRuleAbstract",t_stra.Num)

	dbopt.Sqlcmd = "DELETE FROM sftRuleSend WHERE numRule = ? "
	b_ret = dbopt.ExcAlone(t_stra.Num);if !b_ret{return}
	util.L3I("%s deleted into sftRuleAbstract",t_stra.Num)

	dbopt.Sqlcmd = "INSERT INTO sftRuleSend(numRule,numUser,namex,cls,path) VALUES(?,?,?,?,?)"
	for _,itx:= range t_stra.ListDev.List{
		bret := dbopt.Exc(t_stra.Num,itx.NumUser,itx.Namex,t_stra.Cls,itx.Path);
		b_ret = bret||b_ret
	}
	util.L3I("%s inserted into sftRuleSend num=%d",t_stra.Num,len(t_stra.ListDev.List))

	return
}





var sftMgr StrategySoft
//软件策略详情获取和上传路由
func SoftManager(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I(t_ask.Method)

	var bret = true
	if t_ask.Method=="GET"{
		num:=t_ask.FormValue("num")
		_,r_bts,_ := sftMgr.getStrategy(num)
		t_res.Write(r_bts)
	} else if t_ask.Method=="POST" {
		bts, err := ioutil.ReadAll(t_ask.Body);if err!=nil{
			util.L4E("ioutil.ReadAll(t_ask.Body) "+err.Error())
			bret = false
		} else {
			bret = sftMgr.saveStrategy(bts)
		}
	}

	if !bret {t_res.WriteHeader(http.StatusNotAcceptable)}
}

func init(){
	http.HandleFunc("/strategy/soft",SoftManager)//软件策略路由
}
