//this package used for employee and orgnization manage
package employee

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	//"database/sql"
)
import (
	"dbbase"
	"fmt"
	"util"
)


type sxMan struct {
	Path     string   `json:"path"`
	Ukey     string   `json:"ukey"`
	NumDev   string   `json:"devnum"`
	Emial    string   `json:"email"`
	Name     string   `json:"name"`
	depart   []string `json:"depart"`
	Pwdlogin string   `json:"pwdlogin"`

	Gender   int `json:"gender"`
	Priviege int `json:"privilege"`

	brother []sxMan
}


type sxEmp struct{
	lstMan []sxMan
	org sxOrg
	bload bool
}

type sxOrg struct {
	Path    string `json:path`
	dicKey  []string
	Curkey  string
	Depth   int `json:depth`
	Brother []sxOrg
	Child   []*sxOrg
	pathclt [][]*sxOrg
	Men     []sxMan `json:Men`
}

type sxRetJsDep struct {
	Num   int      `json:"num"`
	Path  string   `json:"path"`
	Sep  string   `json:"sep"`
	Depth int      `json:"depth"`
	Lst   []string `json:"List"`
}
type sxRetJsMen struct {
	Num   int     `json:"num"`
	Path  string  `json:"path"`
	Sep  string   `json:"sep"`
	Depth int     `json:"depth"`
	Lst   []sxMan `json:"List"`
}


func (p *sxEmp)clear(){
	p.lstMan = p.lstMan[0:0]
	p.org.clear()
	p.bload = false
}

func (p *sxEmp)load(){
	p.clear()

	p.bload = p.readAllMan()
	for _,itm := range p.lstMan {
		p.org.insertChild(&itm)
	}

	util.L3I("man info has bee re-load,num=%d",len(p.lstMan))
}

func (p *sxEmp) GetLstDepat(t_path, t_sep string) (r_lst []string, r_json string) {
	if !p.bload {p.load()}
	r_lst,r_json = p.org.getLstDepat(t_path,t_sep)
	return 
}

func (p *sxEmp) GetLstMan(t_path, t_sep string) (r_lst []sxMan, r_json []byte) {
	if !p.bload {p.load()}

	r_lst,r_json =  p.org.getLstMan(t_path,t_sep)
	return
}



func (p *sxMan) getKeyPath(t_grad int) (s_ret, s_lst string, b_ret bool) {
	if t_grad+1 > len(p.depart) {
		return
	}

	for ix := 0; ix <= t_grad; ix++ {
		if ix == 0 {
			s_ret = s_ret + p.depart[ix]
		} else {
			s_ret = s_ret + getSep() + p.depart[ix]
		}
		s_lst = p.depart[ix]
	}

	b_ret = true
	return
}

// parse one man's department according path
func (p *sxMan) parse(t_sep string, t_LtoR bool) {
	fx := strings.Index
	pathx := p.Path
	if !t_LtoR {
		fx = strings.LastIndex
	}

	p.depart = p.depart[0:0]
	for true {
		var strOne string
		ix := fx(pathx, t_sep)
		if ix == -1 {
			p.depart = append(p.depart, pathx)
			break
		}
		if t_LtoR {
			strOne = pathx[:ix]
			pathx = pathx[ix+1:]
		} else {
			strOne = pathx[ix+1:]
			pathx = pathx[:ix]
		}
		p.depart = append(p.depart, strOne)
	}

	for ix, _ := range p.depart {
		var keyx string = ""
		for iy := 0; iy <= ix; iy++ {
			if len(keyx) > 0 {
				keyx = keyx + t_sep + p.depart[iy]
			} else {
				keyx = keyx + p.depart[iy]
			}
		}
		if len(keyx) < 1 {
			keyx = t_sep
		}
	}
}

//add one man info into list
func (p *sxEmp) push(t_man *sxMan) {
	p.lstMan = append(p.lstMan, *t_man)
}

func (p *sxEmp) readAllMan() (b_ret bool) {
	dbopt, bret := dbbase.NewSxDB(&util.GetSftCfg().Db, "readAllMan")
	if !bret {
		return
	}
	defer dbopt.Close()

	//dbopt.Sqlcmd = "SELECT pathx,ukey,email,namex,pwdlogin,gender,priviege FROM employee"
	p.lstMan = p.lstMan[0:0]
	dbopt.Sqlcmd = "SELECT a.pathx,a.ukey,a.email,a.namex,a.pwdlogin,a.gender,a.priviege,b.numDev "
	dbopt.Sqlcmd += "FROM employee a LEFT JOIN terDevBasicInfo b "
	dbopt.Sqlcmd += "ON a.ukey=b.numUsrKey"
	if !dbopt.Query() { return }
	for dbopt.Rows.Next() {
		var ele sxMan
		dbopt.Rows.Scan(&ele.Path, &ele.Ukey, &ele.Emial, &ele.Name, &ele.Pwdlogin, &ele.Gender, &ele.Priviege, &ele.NumDev)
		ele.parse(getSep(), true)
		p.push(&ele)
	}


	b_ret = true
	return
}


func (p *sxOrg) clear() {
	p.Brother = p.Brother[0:0]
	p.Child = p.Child[0:0]
	p.Curkey = ""
	p.Depth = 0
	p.dicKey = p.dicKey[0:0]
	p.Men = p.Men[0:0]
	p.Path = ""
}

func (p *sxOrg) tstJson() {
	var b1, b2, b3, c1, c2, cd1, root sxOrg

	root.Path = ""
	root.Curkey = "root"
	b1.Path = "b"
	b1.Curkey = "b1"
	b1.Depth = 1
	b2.Path = "b"
	b2.Curkey = "b2"
	b2.Depth = 1
	b3.Path = "b"
	b3.Curkey = "b3"
	b3.Depth = 1
	c1.Path = "b/c"
	c1.Curkey = "c1"
	c1.Depth = 2
	c2.Path = "b/c"
	c2.Curkey = "c2"
	c2.Depth = 2
	cd1.Path = "b/c/d"
	cd1.Curkey = "cd1"
	cd1.Depth = 3

	root.Brother = append(root.Brother, b1)
	root.Brother = append(root.Brother, b2)
	root.Brother = append(root.Brother, b3)

	c1.Child = append(c1.Child, &cd1)
	root.Child = append(root.Child, &c1)
	root.Child = append(root.Child, &c2)
	bts, err := json.Marshal(&root)
	if err != nil {
		util.L3I("json.Marshal(&root) %s", err.Error())
		return
	}
	util.L3I(string(bts))

	root.itorx()

	var lst []string
	lst = append(lst, "")
	lst = append(lst, "m")
	// lst = append(lst,"c")
	// lst = append(lst,"d")
	// lst = append(lst,"e")

	// for ix,_:= range lst {
	// 	root.matchFater(lst[:ix+1])
	// }

	var final *sxOrg = &root
	for ix := len(lst); ix > 0; ix-- {
		ret := root.matchFater(lst[:ix])
		if ret != nil {
			final = ret
			break
		}
	}

	util.L3I("Final %d %s %s", final.Depth, final.Path, final.Curkey)

	bx := false
	for ix, _ := range final.Brother {
		if lst[final.Depth+1] == final.Brother[ix].Curkey {
			bx = true
		}
	}
	if !bx {
		var ele sxOrg
		ele.Depth = final.Depth
		ele.Curkey = lst[final.Depth+1]
		final.Brother = append(final.Brother, ele)
	}

	final.itorx()
}

func (p *sxOrg) itorx() {
	util.L3I("%d %s %s ", p.Depth, p.Curkey, p.Path)
	for ix, _ := range p.Brother {
		p.Brother[ix].itorx()
	}

	for ix, _ := range p.Child {
		p.Child[ix].itorx()
	}
}

func (p *sxOrg) itorMan(t_lst *[]sxMan) {
	util.L1T("%d %s %s num=%d", p.Depth, p.Curkey, p.Path, len(p.Men))

	for _, itm := range p.Men {
		*t_lst = append(*t_lst, itm)
	}

	// for ix, _ := range p.Brother {
	// 	p.Brother[ix].itorMan(t_lst)
	// }

	for ix, _ := range p.Child {
		p.Child[ix].itorMan(t_lst)
	}
}

func (p *sxOrg) matchFater(t_path []string) (r_fater *sxOrg) {
	var keyx, cuKey string
	if p == nil {
		return
	}
	if p.Depth > len(t_path) {
		return
	}

	for ix := 0; ix < len(t_path); ix++ {
		if ix == 0 {
			keyx = keyx + t_path[ix]
		} else {
			keyx = keyx + getSep() + t_path[ix]
		}
	}
	for ix := 0; ix < p.Depth; ix++ {
		//util.L3I("%d %d %d",ix,p.Depth,len(t_path))
		if ix == 0 {
			cuKey = cuKey + t_path[ix]
		} else {
			cuKey = cuKey + getSep() + t_path[ix]
		}
	}

	util.L1T("%d [%s]---for---[%s]", p.Depth, p.Path, keyx)
	if keyx == p.Path {
		util.L1T("found-------------------------------%d %s %s", p.Depth, p.Path, p.Curkey)
		return p
	}

	if p.Depth == len(t_path) && keyx != cuKey {
		for ix, _ := range p.Brother {
			r_fater = p.Brother[ix].matchFater(t_path)
			if r_fater != nil {
				return
			}
		}
	} else if p.Depth < len(t_path) {
		for ix, _ := range p.Child {
			r_fater = p.Child[ix].matchFater(t_path)
			if r_fater != nil {
				return
			}
		}
	}

	return
}

func (p *sxOrg) saveJson(t_path string) {
	bts, _, bret := p.toJson()
	if !bret {
		return
	}

	ioutil.WriteFile(t_path, bts, os.ModePerm)
}

func (p *sxOrg) string() (r_str string) {
	r_str = fmt.Sprintf("%d %s %s", p.Depth, p.Curkey, p.Path)
	return
}

func (p *sxOrg) toJson() (r_bts []byte, r_json string, b_ret bool) {
	var err error
	r_bts, err = json.Marshal(p)
	if err != nil {
		util.L4E("sxOrg toJson Marshal" + err.Error())
		return
	}

	b_ret = true
	r_json = string(r_bts)
	return
}

func (p *sxOrg) getAddr(t_path []string) (t_addr *sxOrg) {
	var pathkey string
	for ix, itm := range t_path {
		if ix > 0 {
			pathkey = pathkey + "/" + itm
		} else {
			pathkey = pathkey + itm
		}

		if ix == 0 {
			//for iy,ity := range p.Brother{
			//if ity.Curkey == t_path[ix]
		}
	}
	return
}

func (p *sxOrg) getLstDepat(t_path, t_sep string) (r_lst []string, r_json string) {
	util.L3I("get " + t_path + " sep=" + t_sep)
	lst := strings.Split(t_path, t_sep)
	px := p.matchFater(lst)
	if px == nil {
		return
	}

	util.L1T("%v", *px)
	for _, itm := range px.Child {
		str := itm.Curkey //+" "+itm.Path+fmt.Sprintf(" %d",itm.Depth)
		r_lst = append(r_lst, str)
		//util.L3I(itm.string())
	}

	var depx sxRetJsDep
	depx.Lst = r_lst
	depx.Num = len(r_lst)
	depx.Depth = px.Depth
	depx.Path = t_path
	depx.Sep = getSep()
	bts, err := json.Marshal(&depx)
	if err != nil {
		util.L4E("json.Marshal(&depx) " + err.Error())
		return
	}
	r_json = string(bts)

	util.L3I("res %s num=%d",t_path,depx.Num)
	return
}

func (p *sxOrg) getLstMan(t_path, t_sep string) (r_lst []sxMan, r_json []byte) {
	util.L3I("get " + t_path + " sep=" + t_sep)
	lst := strings.Split(t_path, t_sep)
	px := p.matchFater(lst)
	if px == nil {
		util.L3I("find nothing")
		return
	}

	px.itorMan(&r_lst)
	var retJs sxRetJsMen
	retJs.Depth = px.Depth
	retJs.Num = len(r_lst)
	retJs.Path = t_path
	retJs.Lst = r_lst
	retJs.Sep = getSep()

	var err error
	r_json, err = json.Marshal(&retJs)
	if err != nil {
		util.L3I("json.Marshal(&retJs) " + err.Error())
		return
	}

	util.L3I("%s num=%d", retJs.Path, retJs.Num)
	return
}

func (p *sxOrg) insertBrother(t_man *sxMan) {
	pfx := p.matchFater(t_man.depart)
	util.L3I("%v", pfx)
}

func (p *sxOrg) insertChild(t_man *sxMan) {
	if p.Path == t_man.Path {
		p.Men = append(p.Men, *t_man)
	}
	if p.Depth >= len(t_man.depart) {
		return
	}
	if p == nil {
		return
	}

	var px *sxOrg = p
	for ix := len(t_man.depart); ix >= 0; ix-- {
		ret := px.matchFater(t_man.depart[0:ix])
		if ret != nil {
			px = ret
			break
		}
	}

	var childx sxOrg
	childx.Curkey = t_man.depart[p.Depth]
	childx.Path, childx.Curkey, _ = t_man.getKeyPath(px.Depth)
	if childx.Path == p.Path {
		if px.Path == t_man.Path {
			px.Men = append(px.Men, *t_man)
		}
		//util.L3I("---------------------1.4")
		util.L1T("childx.Path=%s p.Depth=%d p.Path=%s  t_man.Path=%s px.Path=%s", childx.Path, p.Depth, p.Path, t_man.Path, px.Path)
		//util.L3I("---------------------1.4.1")
		return
	}
	childx.Depth = px.Depth + 1
	px.Child = append(px.Child, &childx)
	//util.L3I(childx.string())

	childx.insertChild(t_man)
}

func getSep()(t_sep string){
	t_sep = util.GetSftCfg().ServFile.OrgSep

	if len(t_sep)<=0 {t_sep=">"}
	return
}