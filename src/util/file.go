package util
import(
	"log"
	"os"
	"fmt"
	"github.com/satori/go.uuid"
)
const (
	Cst_sept = "/"
)

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return os.IsExist(err)
}

func logx(t_msg string){
	log.Println("file  "+t_msg)
}

func GetPathEle(t_path string)(r_file,r_fldID,r_Folder string){
	var lstStr [100]string
	var ix,numx int =0,0

	for ix < len(t_path) {
		if t_path[ix] == '/' {
			numx++
			ix++
			continue
		}
		lstStr[numx] = lstStr[numx] + string(t_path[ix])

		ix++
	}

	// for iy,item:= range lstStr {
	// 	log.Println(iy,"  ",item)
	// }


	var folder,filex,folderID string
	filex = lstStr[numx]
	folderID = lstStr[numx-1]
	ix = 1
	for ix < numx {
		if len(lstStr [ix])>0 {
			if ix==1 {
				folder = folder + lstStr [ix]
			} else {
				folder = folder +"/"+ lstStr [ix]
			}
		}
		log.Println(ix,folder)
		ix++
	}

	log.Println(t_path)
	log.Println(folder,"|",filex)

	return folder,folderID,filex
}

func SaveFileBytes(filename string,buf [] byte) (r_path string, b_ret bool){
	t,err :=os.Create(filename)
	if(err!=nil){
		log.Println("SaveFileBytes fail to create file "+filename+" ",err.Error())
		return "", false
	}
	defer t.Close()

	if _,err := t.Write(buf);err!=nil{
		log.Println("saveFile fail to write data to file "+filename+" ",err.Error())		
		return "" ,false
	}

	return filename,true
}

func Guid()(r_guid string){
	u1 := uuid.NewV4()
	strRet := fmt.Sprintf("%s",u1)
	return strRet
}
