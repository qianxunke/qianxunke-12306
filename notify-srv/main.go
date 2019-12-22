package main

import (
	"encoding/json"
	"fmt"
	"github.com/FrontMage/xinge"
	"github.com/FrontMage/xinge/auth"
	"github.com/FrontMage/xinge/req"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"log"
	"net/http"
	"notify-srv/config"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	stationMap map[string]string
	m          sync.Mutex
)

type City struct {
	Id        uint `gorm:"primary_key"`
	CName     string
	CPinyin   string
	CCode     string
	CProvince string
}

//判断已经可以买票，可以提前30天
func isCanQuery(trainDate string) float64 {
	a, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	b, _ := time.Parse("2006-01-02", trainDate)
	d := b.Sub(a)

	return d.Hours() / 24

}

func main() {
	auther := auth.Auther{AppID: "0b532c7673194", SecretKey: "edf0455d5e55fc128203bbd309b1aa91"}
	pushReq, _ := req.NewSingleAndroidAccountPush("1833****2052", "抢票捷报", "尊敬的"+"2052"+"：EsayGo已为您抢票成功，请在30分钟内到12306官方网站或APP完成支付。")

	auther.Auth(pushReq)

	c := &http.Client{}
	rsp, _ := c.Do(pushReq)
	defer rsp.Body.Close()
	body, _ := ioutil.ReadAll(rsp.Body)

	r := &xinge.CommonRsp{}
	_ = json.Unmarshal(body, r)
	fmt.Printf("%+v", r)

}

func Init() {
	m.Lock()
	defer m.Unlock()
	if stationMap != nil {
		return
	}
	stationMap = make(map[string]string)
	//解析数据
	f, err := os.Open("./stations/station_name.js")
	defer f.Close()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	sArr, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	str := string(sArr)
	s1 := sArr[strings.Index(str, "'")+1 : len(str)-2]
	str2 := string(s1)
	strArr := strings.Split(str2, "@")
	for i, s := range strArr {
		if i == 0 {
			continue
		}
		city := &City{}
		first := strings.Index(s, "|")
		second := strings.Index(s[first+1:], "|")
		third := strings.Index(s[first+1+second+1:], "|")
		four := strings.Index(s[first+1+second+1+third+1:], "|")
		city.CName = s[(first + 1):(first + 1 + second)]
		city.CCode = s[(first + 1 + second + 1):(first + 1 + second + 1 + third)]
		city.CPinyin = s[(first + 1 + second + 1 + third + 1):(first + 1 + second + 1 + third + 1 + four)]
		log.Printf("city :%v\n", city)
		config.MasterEngine().Create(&city)
		time.Sleep(time.Millisecond * 10)
	}
}
