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

type TCinema struct {
	base.BaseCinema
}

func saveCinema() {
	InsertCityList(GetCity())
}


func getCinemasFromDB() *list.List {
	lcitys := list.New().Init()
	db, err := sql.Open("mysql", "root:1CUI@/piaofang")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	stmtOut, err := db.Prepare("SELECT id,name,typeindex FROM cinema")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	rows, err := stmtOut.Query()
	for rows.Next() {
		cinema := new(TCinema)
		err := rows.Scan(&cinema.Id, &cinema.Name, &cinema.TypeIndex)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		//		fmt.Println(city.Name)
		lcitys.PushBack(cinema)
	}
	return lcitys
}
//
//func GetCinema() *list.List {
//	lcitys := list.New().Init()
//	resp, err := http.Get("http://dataservices.wangpiao.com/Portal/ajaxcms/ajaxjson_cinemalist.aspx?CityID=);)
//	if err != nil {
//		panic(err)
//	}
//	doc, err := goquery.NewDocumentFromResponse(resp)
//	if err != nil {
//		panic(err)
//	}
//	cityatags := doc.Find("body").Find("div[class=tab-content]").Find("li").Find("a")
//	cityatags.Each(func(num int, s *goquery.Selection) {
//		city := new(TCity)
//		city.Name = s.Text()
//		id, _ := strconv.Atoi(s.AttrOr("cityid", ""))
//		city.TypeIndex = id
//		lcitys.PushBack(city)
//	})
//	return lcitys
//}


func InsertCinemaList(l *list.List) {
	InsertList(l, "insert into cinema(id, name, typeindex, typeindexs)  values ", func(any interface{}) string {
		cinema := any.(*TCinema)
		//		fmt.Println(cinema.Name)
		fmt.Println("(0,'" + cinema.Name + "'," + strconv.Itoa(cinema.TypeIndex) + ",'" + strconv.Itoa(cinema.TypeIndex) + "')")
		return "(0,'" + cinema.Name + "'," + strconv.Itoa(cinema.TypeIndex) + ",'" + strconv.Itoa(cinema.TypeIndex) + "')"
	})
}


func getCinemaSingleCity(cityid int) []cinema {
	resp, err := http.Get("http://dataservices.wangpiao.com/Portal/ajaxcms/ajaxjson_cinemalist.aspx?CityID" + strconv.Itoa(cityid))
	defer resp.Body.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	body, err := ioutil.ReadAll(resp.Body)
	body = body[1 : len(body)-1]
	//	jso, err := simplejson.NewJson(body)
	//	cinemas, err := jso.Array()
	//	fmt.Println(len(cinemas))
	//	for _, v := range cinemas {
	//		jmap := v.(map[string]interface{})
	//		for key, value := range jmap {
	//			fmt.Println(key)
	//		}
	//	}
	var cinemas []cinema
	err = json.Unmarshal(body, &cinemas)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", cinemas)
	return cinemas
}

func GetCinema(lcitys *list.List) {
	for e := lcitys.Front(); e != nil; e = e.Next() {
		lcinemas := list.New().Init()
		city := e.Value.(*TCity)
		cinemas := getCinemaSingleCity(city.TypeIndex)
		for i := 0; i < len(cinemas); i++ {
			cinema := new(TCinema)
			cinema.Name = cinemas[i].ICinemaName
			cinema.TypeIndex = cinemas[i].CinemaIndex
			lcinemas.PushBack(cinema)
		}
		InsertCinemaList(lcinemas)
	}
}
type cinema struct {
	CinemaIndex int
	ICinemaName string
}


