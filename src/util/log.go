package util

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	"sync"
)

type sxLog struct {
	path    string
	pathBK    string
	sizeCur int
	sizeMax int
	nlevel  int
	nOutObj int
	lckx sync.Mutex
	file *os.File
}

var lgx sxLog

const (
	cst_depth int = 2

	cst_lgL1 = iota
	cst_lgL2
	cst_lgL3
	cst_lgL4
)


func InitLog(t_pathx,t_pathy string,t_level ,t_obj,t_maxsize int){
	lgx.path,lgx.pathBK,lgx.nlevel,lgx.sizeMax,lgx.nOutObj = t_pathx,t_pathy,t_level,t_maxsize,t_obj
	lgx.createFile()
}

func (p *sxLog)createFile() {
	p.sizeCur = GetFileSize(p.path)
	var err error
	p.file, err = os.OpenFile(p.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		log.Printf("%s\r\n", err.Error())
	}
}

func (p *sxLog) printx(t_tab string, format string, v ...interface{}) {
	p.lckx.Lock()
	defer p.lckx.Unlock()

	_, file, line, _ := runtime.Caller(cst_depth)
	_, _, fname := GetPathEle(file)

	var strMsg string
	strtmx := time.Now().Format("2006-01-02 15:04:05")
	strMsg = fmt.Sprintf("["+t_tab+"] "+format, v...) + fmt.Sprintf("  %s:%d", fname, line)
	strF := strtmx + "  " + strMsg +"\r\n"
	if (p.nOutObj & 0x01)>0 { log.Println(strMsg) }
	if (p.nOutObj & 0x02)>0 { 
		lgx.file.Write([]byte(strF)) 
	}
	lgx.sizeCur += len(strF)

	if lgx.sizeCur > lgx.sizeMax { lgx.rotate() }

	
}

func (p* sxLog)rotate(){
	p.file.Close()
	os.Remove(lgx.pathBK)
	os.Rename(lgx.path,lgx.pathBK)
	p.sizeCur = 0
	p.createFile()
}

func L1T(format string, v ...interface{}) {
	if lgx.nlevel > cst_lgL1 {
		return
	}
	lgx.printx("T", format, v...)
}

func L2I(format string, v ...interface{}) {
	if lgx.nlevel > cst_lgL2 {
		return
	}
	lgx.printx("I", format, v...)
}

func L3E(format string, v ...interface{}) {
	if lgx.nlevel > cst_lgL3 {
		return
	}
	lgx.printx("E", format, v...)
}

func L4F(format string, v ...interface{}) {
	if lgx.nlevel > cst_lgL4 {
		return
	}
	lgx.printx("F", format, v...)
}
