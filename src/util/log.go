package util

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type sxLog struct {
	path    string
	pathBK  string
	sizeCur int
	sizeMax int
	nlevel  int
	nOutObj int
	lckx    sync.Mutex
	file    *os.File
}

var lgx sxLog

const (
	cst_depth int = 2

	cst_lgL1T = iota
	cst_lgL2D
	cst_lgL3I
	cst_lgL4E
	cst_lgL5F
)

func InitLog(t_pathx, t_pathy string, t_level, t_obj, t_maxsize int) {
	lgx.path, lgx.pathBK, lgx.nlevel, lgx.sizeMax, lgx.nOutObj = t_pathx, t_pathy, t_level, t_maxsize, t_obj
	lgx.createFile()
}

func (p *sxLog) createFile() {
	p.sizeCur = GetFileSize(p.path)
	var err error
	p.file, err = os.OpenFile(p.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		log.Printf("%s\r\n", err.Error())
	}
}

func getFileStack(t_depth, t_garde int) (r_msg, r_func string) {
	for ix := 0; ix < t_garde; ix++ {
		pc, file, line, _ := runtime.Caller(t_depth - ix)
		fname := GetFileName(file)
		fx := runtime.FuncForPC(pc)
		funcname := fx.Name()
		lst := strings.Split(funcname, ".")

		r_msg += fmt.Sprintf("/%s %s:%d", funcname, fname, line)
		r_func = lst[len(lst)-1]
	}

	return
}

func (p *sxLog) printx(t_tab string, format string, v ...interface{}) {
	p.lckx.Lock()
	defer p.lckx.Unlock()

	pc, file, line, _ := runtime.Caller(cst_depth)
	fname := GetFileName(file)
	fx := runtime.FuncForPC(pc)
	funname := fx.Name()
	lst := strings.Split(funname, ".")

	var strMsg string
	strtmx := time.Now().Format("2006-01-02 15:04:05")
	strMsg = fmt.Sprintf("["+t_tab+"] "+lst[len(lst)-1]+" "+format, v...) + fmt.Sprintf("  %s:%d", fname, line)
	//strMsg = fmt.Sprintf("["+t_tab+"] "+format, v...) + fmt.Sprintf("  %s",getFileStack(4,4))
	strF := strtmx + "  " + strMsg + "\r\n"
	if (p.nOutObj & 0x01) > 0 {
		log.Println(strMsg)
	}
	if (p.nOutObj & 0x02) > 0 {
		lgx.file.Write([]byte(strF))
	}
	lgx.sizeCur += len(strF)

	if lgx.sizeCur > lgx.sizeMax {
		lgx.rotate()
	}
}

func (p *sxLog) printEx(t_dep, t_grad int, t_tab string, format string, v ...interface{}) {
	p.lckx.Lock()
	defer p.lckx.Unlock()

	// _, file, line, _ := runtime.Caller(cst_depth)
	// fname := GetFileName(file)

	var strMsg string
	strtmx := time.Now().Format("2006-01-02 15:04:05")
	//strMsg = fmt.Sprintf("["+t_tab+"] "+format, v...) + fmt.Sprintf("  %s:%d", fname, line)
	strmsg, strFN := getFileStack(t_dep, t_grad)
	strMsg = fmt.Sprintf("["+t_tab+"] "+strFN+" "+format, v...) + fmt.Sprintf("  %s", strmsg)
	strF := strtmx + "  " + strMsg + "\r\n"
	if (p.nOutObj & 0x01) > 0 {
		log.Println(strMsg)
	}
	if (p.nOutObj & 0x02) > 0 {
		lgx.file.Write([]byte(strF))
	}
	lgx.sizeCur += len(strF)

	if lgx.sizeCur > lgx.sizeMax {
		lgx.rotate()
	}
}

func (p *sxLog) rotate() {
	p.file.Close()
	os.Remove(lgx.pathBK)
	os.Rename(lgx.path, lgx.pathBK)
	p.sizeCur = 0
	p.createFile()
}

func L1T(format string, v ...interface{}) {
	if lgx.nlevel > cst_lgL1T {
		return
	}
	lgx.printx("T", format, v...)
}

func L2D(format string, v ...interface{}) {
	if lgx.nlevel > cst_lgL2D {
		return
	}
	lgx.printx("D", format, v...)
}

func L3I(format string, v ...interface{}) {
	if lgx.nlevel > cst_lgL3I {
		return
	}
	lgx.printx("I", format, v...)
}

func L4E(format string, v ...interface{}) {
	if lgx.nlevel > cst_lgL4E {
		return
	}
	lgx.printx("E", format, v...)
}
func L4Ex(t_dep, t_grad int, format string, v ...interface{}) {
	if lgx.nlevel > cst_lgL4E {
		return
	}
	lgx.printEx(t_dep, t_grad, "E", format, v...)
}

func L5F(format string, v ...interface{}) {
	if lgx.nlevel > cst_lgL5F {
		return
	}
	lgx.printx("F", format, v...)
}

func L5Fx(t_dep, t_grad int, format string, v ...interface{}) {
	if lgx.nlevel > cst_lgL5F {
		return
	}
	lgx.printEx(t_dep, t_grad, "F", format, v...)
}
