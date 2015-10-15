package main

import (
	"base"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	//	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	//	"github.com/bitly/go-simplejson"
	//	"go.marzhillstudios.com/pkg/go-html-transform/css/selector"
	//	"go.marzhillstudios.com/pkg/go-html-transform/h5"
	//	"go.marzhillstudios.com/pkg/go-html-transform/html/transform"
	//	"gopkg.in/mgo.v2"
	//	"gopkg.in/mgo.v2/bson"
	"container/list"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type TCity struct {
	base.BaseCity
}
func saveCity() {
	InsertCityList(GetCity())
}

func GetCity() *list.List {
	lcitys := list.New().Init()
	resp, err := http.Get("http://www.wangpiao.com/")
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		panic(err)
	}
	cityatags := doc.Find("body").Find("div[class=tab-content]").Find("li").Find("a")
	cityatags.Each(func(num int, s *goquery.Selection) {
		city := new(TCity)
		city.Name = s.Text()
		id, _ := strconv.Atoi(s.AttrOr("cityid", ""))
		city.TypeIndex = id
		lcitys.PushBack(city)
	})
	return lcitys
}
func getCitysFromDB() *list.List {
	db, err := sql.Open("mysql", "root:1CUI@/piaofang")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	stmtOut, err := db.Prepare("SELECT id,name,typeindex FROM city")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	rows, err := stmtOut.Query()

	for rows.Next() {
		
//		city := new(TCity)
//		var name string
//		err := rows.Scan()
//		if err != nil {
//			panic(err.Error()) // proper error handling instead of panic in your app
//		}
//		lhasnames.PushBack(name)
	}
	return StringListToArray(lhasnames)
}

func InsertCityList(l *list.List) {
	InsertList(l, "insert into city(id, name, typeindex, typename)  values ", func(any interface{}) string {
		city := any.(*TCity)
		return "(0,'" + city.Name + "','" + strconv.Itoa(city.TypeIndex) + "','')"
	})
}


func getCitysFromDB() *list.List {
	lcitys := list.New().Init()
	db, err := sql.Open("mysql", "root:1CUI@/piaofang")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	stmtOut, err := db.Prepare("SELECT id,name,typeindex FROM city")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	rows, err := stmtOut.Query()
	for rows.Next() {
		city := new(TCity)
		err := rows.Scan(&city.Id, &city.Name, &city.TypeIndex)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		fmt.Println(city.Name)
		lcitys.PushBack(city)
	}
	return lcitys
}


