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
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)
import (	
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

//根据软件策略编号获取软件清单信息
//如果num为空 则从数据库获取软件清单 否则从磁盘的json文件获取策略的详情
func getSoftList(t_num string)(r_list sxSftItemList){
	return
}
//从数据库获取清单信息
func (p *StrategySoft)getSoftListFromDB()(r_list sxSftItemList){
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
		r_list.List = append(r_list.List,ele)
	}

	util.L3I("softwaare num=%d", len(r_list.List))
	return
}
//从JSON文件获取软件策略详情
func getSoftRuelDetail(t_num string)(r_list sxSftItemList){
	return
}


var sftMgr StrategySoft
func softManager(t_res http.ResponseWriter,t_ask *http.Request){
	if t_ask.Method=="GET"{
		sftMgr.getSoftListFromDB()
	}
	return
}



func init(){
	http.HandleFunc("/strategy/soft",softManager)
}
