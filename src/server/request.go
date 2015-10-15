package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

const filename string = "gopid-synconlinestatus"

func getroute(c *mgo.Collection) []string {
	var macs []string
	_ = c.Find(nil).Distinct("mac", &macs)
	return macs
}

func updateversion(c *mgo.Collection, mac string, ch chan bool) {
	resp, err := http.Get("http://router-api.xiaoyun.com/api/v1/router/" + mac)
	defer resp.Body.Close()
	if err != nil {
	}
	body, err := ioutil.ReadAll(resp.Body)
	jso, err := simplejson.NewJson(body)
	version, err := jso.Get("item").Get("version").String()
	if err == nil {
		fmt.Println(version)
		//		s := string(body[0:len(body)])
		//		index := strings.Index(s, "version")
		//		//用json比较麻烦，需要建一个struct才能转换，直接用字符串匹配
		//		version := string(s[index+10 : index+21])
		//		in := version[0:1]
		//		if in != "\"" {
		c.Update(bson.M{"mac": mac}, bson.M{"$set": bson.M{"version": version}})
		//		}
	}
	ch <- true
}

func syncApOnlineStatus() {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("shenji").C("ap")
	macs := getroute(c)
	fmt.Println(len(macs))
	chs := make([]chan bool, len(macs))
	for i := 0; i < len(macs); i++ {
		chs[i] = make(chan bool)
		go updateversion(c, macs[i], chs[i])
		//这样写会有同时连接数过多的问题
	}
	for _, ch := range chs {
		<-ch
	}
	fmt.Println("all ap status sync success end")
}

func syncApOnlineStatus2() {
	pinum := 4
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("shenji").C("ap")
	macs := getroute(c)
	maccount := len(macs)
	picount := maccount / pinum //21 5
	yu := maccount % pinum
	if yu != 0 {
		picount++ //6
	}
	chs := make([]chan bool, pinum)
	for i := 0; i < picount; i++ {
		for j := 0; j < pinum; j++ {
			chs[j] = make(chan bool)
			if (i != (picount - 1)) || j < yu {
				go updateversion(c, macs[pinum*i+j], chs[j])
			} else {
				go updateversion(c, "", chs[j])
			}
		}
		for _, ch := range chs {
			<-ch
		}
	}
	fmt.Println("all ap status sync success end")
}

func main() {
	syncApOnlineStatus2()
	//		ticker := time.NewTicker(5 * 60 * time.Second)
	//		for {
	//			<-ticker.C
	//			fmt.Println("begin")
	//			syncApOnlineStatus()
	//		}
}
