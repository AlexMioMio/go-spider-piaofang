package wangpiao

import (
	"../base"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"container/list"
	"io/ioutil"
	"net/http"
	"strconv"
	"github.com/PuerkitoBio/goquery"
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
		lcitys.PushBack(cinema)
	}
	return lcitys
}

func InsertCinemaList(l *list.List) {
	base.InsertList(l, "insert into cinema(id, name, typeindex, typeindexs)  values ", func(any interface{}) string {
		cinema := any.(*TCinema)
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
	body = body[1 : len(body) - 1]
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


