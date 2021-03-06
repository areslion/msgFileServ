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
	"strconv"
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
	Enable int `json:"enable"`//在策略中是否被启用
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
type sxStrategyItem struct {
	Num string `json:"num"`//策略的ID
	Name string `json:"name"`//策略名称
	Cls int `json:"class"`//策略类型 1 为软件白名单 2为软件黑名单
}
type sxStrategyList struct{
	List []sxStrategyItem `json:"list"`
}
//软件策略返回结果
type sxRes struct{
	ListSft sxSftItemList `json:"listSoft"`//策略软件清单
	ListDev sxDevItemList `json:"listDev"`//策略应用于的终端清单
	Strategy sxStrategyItem `json:"strategy"`//策略的基本信息
}

//软件过滤器
type sxFilter struct{
	strFilter []string//需要要被过滤的软件关键字
	strKeeper string//放置重叠的关键字系列
	strPlacer []string//需要被替换的关键字
}

//func (p *sxEmp) GetLstDepat(t_path, t_sep string) (r_lst []string, r_json string) 
func (p * sxFilter) init(){	
	cfg := util.GetSftCfg()
	for _,itx:= range cfg.Strategy.Software.Filter {
		p.strFilter = append(p.strFilter,itx.Name)
	}
	for _,itx := range cfg.Strategy.Software.Placer {
		p.strPlacer = append(p.strPlacer,itx.Name)
	}
	
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

//删除指定的软件策略，并返删除失败的策略的清单
func (p *StrategySoft)deleteStrategy(t_bts[] byte)(r_bts[] byte){
	var list sxStrategyList
	err := json.Unmarshal(t_bts,&list); if err!=nil {
		util.L4E("json.Unmarshal(t_bts,&list) "+err.Error())
		r_bts = append(r_bts,t_bts...)
		return
	}

	p.deleteStrxFromDB(list)
	p.deleteStrxFromFile(list)

	// var strAll string
	// for _,itx := range lstDB.List {
	// 	strAll += itx.Num + " "
	// }


	// var listFailed sxStrategyList
	// for _,itx := range list.List{
	// 	if strings.Index(strAll,itx.Num)==-1 {
	// 		var ele sxStrategyItem
	// 		ele.Num = itx.Num
	// 		listFailed.List = append(listFailed.List,ele)
	// 	}
	// }

	// r_bts,err = json.Marshal(&listFailed);if err!=nil{
	// 	util.L4E("json.Marshal(&listFailed) "+err.Error())
	// 	return
	// }

	return
}
//从数据库删除指定的软件策略，并返删除失败的策略的清单
func (p *StrategySoft)deleteStrxFromDB(t_list sxStrategyList)(r_list sxStrategyList){
	dbopt, bret := dbbase.NewSxDB(&util.GetSftCfg().Db, "delete strategy from db")
	if !bret {return}
	defer dbopt.Close()

	for _,itx := range t_list.List{
		dbopt.Sqlcmd = "DELETE FROM sftRuleAbstract WHERE num=? "
		dbopt.Exc(itx.Num)
		dbopt.Sqlcmd = "DELETE FROM sftRuleSend WHERE numRule=? "
		dbopt.Exc(itx.Num)
		util.L3I("delete strategy from db "+itx.Num)
	}

	return
}
//从数据库删除指定的软件策略，并返删除失败的策略的清单
func (p *StrategySoft)deleteStrxFromFile(t_list sxStrategyList)(r_list sxStrategyList){
	for _,itx := range t_list.List{
		var pathx = p.getStrategyPath(itx.Num)
		util.RemoveAll(pathx)
		util.L3I("delete strategy file "+pathx)
	}
	
	return
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
	r_bts,err = ioutil.ReadFile(path);if err!=nil {
		util.L4E(path+" "+err.Error())
		return
	}

	err = json.Unmarshal(r_bts,&r_res);if err!=nil{
		util.L4E("json.Unmarshal(bts,r_res) "+err.Error())
		return
	}

	//从本地获取策略的软件清单之后 从数据库全局再次获取并补充到未选中的软件清单中
	listDB,_ := p.getSoftListFromDB()
	var keperx string
	for _,itx := range r_res.ListSft.List {
		keperx += itx.Namex
	}
	for _,itx := range listDB.ListSft.List {
		if strings.Contains(keperx,itx.Namex)==true {continue}
		itx.Enable = 0
		r_res.ListSft.List = append(r_res.ListSft.List,itx)
	}
	r_bts,err = json.Marshal(r_res);if err!=nil{//将最新的数据重新整理为Bts格式
		util.L4E("json.Marshal(bts,r_res) "+err.Error())
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
		ele.Enable = 0
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
func (p *StrategySoft)getStrategy(t_num string,t_opt int)(r_res sxRes,r_bts []byte,r_str string){
	sfgCfg := util.GetSftCfg()

	
	util.L3I("get strategy opt=%d num="+t_num,t_opt)
	switch t_opt {
	case 0://获取可以进行筛选的软件清单
		r_res,r_bts = p.getSoftListFromDB()
	case 1://获取所有策略清单
		_,r_str,r_bts = p.getStrategyList()
	case 2://获取软件策略忽略和过替换的关键字系列
		var err error
			r_bts,err = json.Marshal(sfgCfg.Strategy.Software); if err!=nil{
				util.L4E("json.Marshal(r_bts,&sfgCfg.Strategy.Software "+err.Error())
				return
			}
	case 3://获取指定终端软件策略，此时num为终端的ID
		r_bts = p.getDevSoftStrategy(t_num)
	case 4://获取指定GUID的软件策略的详情，此时num为策略的guid
		r_res, r_bts = p.getSoftList(t_num)
	default:
		util.L4E("undefine value soft?opt=%d",t_opt)
		return
	}


	// if t_opt==0{
	// 	r_res,r_bts = p.getSoftListFromDB()
	// }else if t_opt==1 {
	// 	num,err := strconv.Atoi(t_num) ; if err!=nil{
	// 		util.L4E("strconv.Atoi(%s) "+err.Error(),t_num)
	// 		return
	// 	}

	// 	switch num {
	// 	case 1:
	// 		var err error
	// 		r_bts,err = json.Marshal(sfgCfg.Strategy.Software); if err!=nil{
	// 			util.L4E("json.Marshal(r_bts,&sfgCfg.Strategy.Software "+err.Error())
	// 			return
	// 		}
	// 	default:
	// 		util.L4E("undefine value soft?num=%d",num)
	// 		return
	// 	}

	// }else if len(t_num)==2 {
	// 	cls,err:=strconv.Atoi(t_num);if err!=nil || cls!=-1{
	// 		util.L4E("strconv.Atoi(%s) %s",t_num,err.Error())
	// 		return 
	// 	}
	// 	_,r_str,r_bts = p.getStrategyList()
	// }else {
	// 	r_res, r_bts = p.getSoftList(t_num)
	// }

	r_str = string(r_bts)
	return
}
//获取策略清单
func (p *StrategySoft)getStrategyList()(r_list sxStrategyList,r_json string,r_bts []byte){
	dbopt, bret := dbbase.NewSxDB(&util.GetSftCfg().Db, "read software startegy list")
	if !bret {return}
	defer dbopt.Close()

	dbopt.Sqlcmd = "SELECT num,namex,cls FROM sftRuleAbstract "
	if !dbopt.Query() { return }
	for dbopt.Next(){
		var ele sxStrategyItem
		dbopt.Scan(&ele.Num,&ele.Name,&ele.Cls)
		r_list.List = append(r_list.List,ele)
	}

	var err error
	r_bts,err = json.Marshal(r_list);if err!=nil{
		util.L4E("json.Marshal(r_list) "+err.Error())
	}
	r_json  = string(r_bts)
	
	util.L3I("read software startegy list=%d",len(r_list.List))
	return
}
//获取t_num指定的终端的软件策略
func (p *StrategySoft)getDevSoftStrategy(t_num string)(r_bts []byte){
	dbopt, bret := dbbase.NewSxDB(&util.GetSftCfg().Db, "read software startegy list")
	if !bret {return}
	defer dbopt.Close()

	var rulestr string
	dbopt.Sqlcmd = "SELECT numRule FROM sftRuleSend WHERE numUser=?"
	dbopt.Query(t_num)
	for dbopt.Next(){
		dbopt.Scan(&rulestr)
		break
	}
	if(len(rulestr)<=0){
		util.L4E("fail to get rule id from db")
		return
	}

	_,r_bts = p.getSoftList(rulestr)

	return
}

//存储软件策略
func (p *StrategySoft)saveStrategy(t_bts []byte)(b_ret bool){
	var strx sxRes
	err := json.Unmarshal(t_bts,&strx);if err!=nil{//change byte to struct
		util.L4E("json.Unmarshal(t_bts,strx) "+err.Error())
	}

	if len(strx.Strategy.Num)!=36 {
		util.L4E("invalid num(%s) length=%d need=%d",strx.Strategy.Num,len(strx.Strategy.Num),36)
		return
	}

	path := p.getStrategyPath(strx.Strategy.Num)
	_, b_ret = util.SaveFileBytes(path, t_bts);if !b_ret {
		util.L4E("write software strategy(%s) failed(%s)",strx.Strategy.Num,path)
	}
	util.L4E("software(%s) saved as",strx.Strategy.Num,path)
	
	b_ret = p.insertDB(strx)
	
	return
}

//将软件策略存储到数据库中的sftRuleAbstract和sftRuleSend两张表中
func (p *StrategySoft)insertDB(t_stra sxRes)(b_ret bool){
	dbopt, bret := dbbase.NewSxDB(&util.GetSftCfg().Db, "insert software startegy into db")
	if !bret {return}
	defer dbopt.Close()

	dbopt.Sqlcmd = " REPLACE INTO sftRuleAbstract(num,namex,cls) VALUES(?,?,?)"
	b_ret = dbopt.Exc(t_stra.Strategy.Num,t_stra.Strategy.Name,t_stra.Strategy.Cls);if !b_ret{return}
	util.L3I("%s replaced into sftRuleAbstract",t_stra.Strategy.Num)

	dbopt.Sqlcmd = "DELETE FROM sftRuleSend WHERE numRule = ? "
	b_ret = dbopt.Exc(t_stra.Strategy.Num);if !b_ret{return}
	util.L3I("%s deleted into sftRuleAbstract",t_stra.Strategy.Num)

	//将属于该用户的所有软件策略删除 仅保留最新的软件策略
	dbopt.Sqlcmd = "DELETE FROM sftRuleSend WHERE numUser = ? "
	for _,itx:= range t_stra.ListDev.List{
		b_ret = dbopt.Exc(itx.NumUser)
	}

	dbopt.Sqlcmd = "INSERT INTO sftRuleSend(numRule,numUser,namex,cls,path) VALUES(?,?,?,?,?)"
	for _,itx:= range t_stra.ListDev.List{
		bret := dbopt.Exc(t_stra.Strategy.Num,itx.NumUser,itx.Namex,t_stra.Strategy.Cls,itx.Path);
		b_ret = bret||b_ret
	}
	util.L3I("%s inserted into sftRuleSend num=%d",t_stra.Strategy.Num,len(t_stra.ListDev.List))

	return
}





var sftMgr StrategySoft
//软件策略详情获取和上传路由
func SoftManager(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I(t_ask.Method)

	var bret = true
	if t_ask.Method=="GET"{
		num:=t_ask.FormValue("num")
		opt:=t_ask.FormValue("opt")
		util.L3I("opt=%s num=%s",opt,num)
		nopt,err := strconv.Atoi(opt); if err!=nil{
			util.L4E("strconv.Atoi(%s)",opt)
			bret = false
		} else {
			_,r_bts,_ := sftMgr.getStrategy(num,nopt)
			t_res.Write(r_bts)
		}
	} else if t_ask.Method=="POST" {
		opt := t_ask.FormValue("opt")
		bts, err := ioutil.ReadAll(t_ask.Body);if err!=nil{
			util.L4E("ioutil.ReadAll(t_ask.Body) "+err.Error())
			bret = false
		} else {
			if opt=="del"{
				r_bts := sftMgr.deleteStrategy(bts)
				t_res.Write(r_bts)
			} else {
				bret = sftMgr.saveStrategy(bts)	
			}
		}
	}

	if !bret {t_res.WriteHeader(http.StatusNotAcceptable)}
}

func init(){
	http.HandleFunc("/strategy/soft",SoftManager)//软件策略路由
}
