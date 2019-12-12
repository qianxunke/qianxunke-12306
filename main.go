package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	/*result1 := login.Login(login.User{UserName: "dh17862709691", Pwd: "736567805"})
	result2 := login.Login(login.User{UserName: "laidanchao", Pwd: "lai19920127"})
	sysinit.Init()
	for true {
		query.QueryTrainMessage("dh17862709691", result1.Conversat)
		query.QueryTrainMessage("laidanchao", result2.Conversat)
		time.Sleep(time.Second * 3)
	}


	*/
	tm := time.Now()
	log.Println(fmt.Sprintf("%d-%d-%d", tm.Year(), tm.Month(), tm.Day()))

}
