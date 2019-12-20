package stations

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

//站点信息

var (
	stationMap map[string]string
	m          sync.Mutex
)

func GetStationValueByKey(key string) (value string) {
	if stationMap == nil {
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
			first := strings.Index(s, "|")
			second := strings.Index(s[first+1:], "|")
			third := strings.Index(s[first+1+second+1:], "|")
			key := s[(first + 1):(first + 1 + second)]
			value := s[(first + 1 + second + 1):(first + 1 + second + 1 + third)]
			stationMap[key] = value
			stationMap[value] = key
		}

	}
	return stationMap[key]

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
		first := strings.Index(s, "|")
		second := strings.Index(s[first+1:], "|")
		third := strings.Index(s[first+1+second+1:], "|")
		key := s[(first + 1):(first + 1 + second)]
		value := s[(first + 1 + second + 1):(first + 1 + second + 1 + third)]
		stationMap[key] = value
		stationMap[value] = key
	}
}
