package main

import (
	"log"
	"os/exec"
	"time"
)

func main(){	
	
	for {
		cmd := exec.Command("tasklist","/FI","IMAGENAME eq SmartDoorClient.exe")
		if bts,err := cmd.Output(); err!=nil{
			log.Println(err.Error())
			log.Println(string(bts))
		} else {
			log.Println("miss object will be luanched")			
			cmd = exec.Command(".\\SmartDoorClient.exe")
			if _,err = cmd.Output();err!=nil{log.Println(err.Error())}
		}

		time.Sleep(1e9)
	}
	

}