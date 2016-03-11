package wangpiao

import (
	"../base"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"container/list"
	"strconv"
	"net/http"
	"io/ioutil"
	"strings"
)

func getShowTimeTypeIndexFromDB() []int {
	lshowindex := list.New().Init()
	db, err := sql.Open("mysql", "root:1CUI@/piaofang")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	stmtOut, err := db.Prepare("SELECT typeshowindex FROM showtime")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	rows, err := stmtOut.Query()
	for rows.Next() {
		var index int
		err := rows.Scan(&index)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		lshowindex.PushBack(index)
	}
	return base.IntListToArray(lshowindex)
}

type ShowTime struct {
	base.BaseShowTime
}

func InsertShowTimeList(l *list.List) {
	base.InsertList(l, "insert into showtime (id,type,typename,typeshowindex,typecinemaindex,typemovieindex,typemoviename,typehallid,typesaleendtime,price,seatcount,typecityindex)  values ", func(any interface{}) string {
		showtime := any.(*ShowTime)
		fmt.Println("(0,0,'网票'," + strconv.Itoa(showtime.TypeShowIndex) + "," + strconv.Itoa(showtime.TypeCinemaIndex) + "," + strconv.Itoa(showtime.TypeMovieIndex) + ",'" + showtime.TypeMovieName + "'," + strconv.Itoa(showtime.TypeHallID) + "," + strconv.FormatInt(showtime.TypeSaleEndTime, 10) + "," + strconv.Itoa(showtime.Price) + "," + strconv.Itoa(showtime.SeatCount) + "," + strconv.Itoa(showtime.TypeCityIndex) + ")")
		return "(0,0,'网票'," + strconv.Itoa(showtime.TypeShowIndex) + "," + strconv.Itoa(showtime.TypeCinemaIndex) + "," + strconv.Itoa(showtime.TypeMovieIndex) + ",'" + showtime.TypeMovieName + "'," + strconv.Itoa(showtime.TypeHallID) + "," + strconv.FormatInt(showtime.TypeSaleEndTime, 10) + "," + strconv.Itoa(showtime.Price) + "," + strconv.Itoa(showtime.SeatCount) + "," + strconv.Itoa(showtime.TypeCityIndex) + ")"
	})
}

func GetShowTime(lcinema *list.List) {
	for e := lcinema.Front(); e != nil; e = e.Next() {
		lshowtimes := list.New().Init()
		cinema := e.Value.(*TCinema)
		showtimes := getShowTimeSingleCinema(cinema.TypeIndex, "2015-10-15")
		for i := 0; i < len(showtimes); i++ {
			jshowtime := showtimes[i]
			showtime := new(ShowTime)
			showtime.SeatCount = jshowtime.SeatCount
			showtime.Type = 0
			showtime.TypeCinemaIndex = jshowtime.CinemaID
			showtime.TypeCityIndex = jshowtime.CityID
			showtime.TypeHallID = jshowtime.HallID
			showtime.TypeMovieIndex = jshowtime.FilmID
			showtime.TypeMovieName = jshowtime.FilmName
			showtime.TypeName = "wangpiao"
			//			showtime.TypeSaleEndTimeS = jshowtime.SaleEndTime
			ts, _ := time.Parse("2006-01-02 15:04:05", jshowtime.SaleEndTime)
			showtime.TypeSaleEndTimeS = ts
			showtime.TypeSaleEndTime = ts.Unix()
			showtime.TypeShowIndex = jshowtime.ShowIndex
			lshowtimes.PushBack(showtime)
		}
		InsertShowTimeList(lshowtimes)
	}
}

type showtime struct {
	ShowIndex   int
	CinemaID    int
	HallID      int
	FilmID      int
	FilmName    string
	//LG: 原版,
	//ShowTime: 2015-10-15 22:00:00,
	SaleEndTime string
	//Status: 1,
	//UPrice: 45,
	//VPrice: 45,
	CityID      int
	//UWPrice: 50,
	//SPType: 1|2|5,
	//SPPrice: 5|0|0,
	//HallName: 12号厅,
	//CPrice: 50,
	//IsImax: false,
	//Dimensional: 3D,
	SeatCount   int
}

type showtimeout struct {
	ErrNo int
	Sign  string
	Msg   string
	Data  []showtime
}

func getShowTimeSingleCinema(cinemaid int, datestr string) []showtime {
	resp, err := http.Get("http://dataservices1.wangpiao.com/API.aspx?Target=Base_FilmShow&Param=CinemaID=" + strconv.Itoa(cinemaid) + "&Date=" + datestr)
	defer resp.Body.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	body, err := ioutil.ReadAll(resp.Body)
	body = body[1 : len(body) - 1]
	var sto showtimeout
	err = json.Unmarshal(body, &sto)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", sto.Data)
	return sto.Data
}

func saveTodayShowTime() {
	lcinemas := getCinemasFromDB()
	GetShowTime(lcinemas)
}

func getSingleShowTimeCurrentPeople(lmap map[int]int, showindex int, ch chan bool) {
	if showindex != 0 {
		greq, _ := http.NewRequest("GET", "http://dataservices.wangpiao.com/Data.aspx?getpageurl=Http%3A//dataservices.wangpiao.com/Portal/ajaxcms/ajax_SeatGrid.aspx&getpageparam=SeqNo%3D" + strconv.Itoa(showindex) + "&format=json&_=" + getNowTimeTsString(), nil)
		greq.Header.Add("Referer", "http://www.wangpiao.com")
		c := &http.Client{}
		resp, err := c.Do(greq)
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		s := string(body[0:len(body)])
		saledcount := strings.Count(s, "#F")
		lmap[showindex] = saledcount
	}
	ch <- true
}

func updateAllShowTimePeople() {
	pinum := 10
	arr := getShowTimeTypeIndexFromDB()
	itemcount := len(arr)
	picount := itemcount / pinum //21 5
	yu := itemcount % pinum
	if yu != 0 {
		picount++ //6
	}
	chs := make([]chan bool, pinum)
	for i := 0; i < picount; i++ {
		pimap := make(map[int]int)
		for j := 0; j < pinum; j++ {
			chs[j] = make(chan bool)
			if (i != (picount - 1)) || j < yu {
				go getSingleShowTimeCurrentPeople(pimap, arr[pinum * i + j], chs[j])
			} else {
				go getSingleShowTimeCurrentPeople(pimap, 0, chs[j])
			}
		}
		for _, ch := range chs {
			<-ch
		}
		UpdateShowTimeSaledByShowIndex(pimap)
	}
	fmt.Println("all show time  update finished")
}

func UpdateShowTimeSaledByShowIndex(mapinfo map[int]int) int64 {
	db, err := sql.Open("mysql", "root:1CUI@/piaofang")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	//UPDATE categories
	//    SET display_order = CASE id
	//        WHEN 1 THEN 3
	//        WHEN 2 THEN 4
	//        WHEN 3 THEN 5
	//    END
	//WHERE id IN (1,2,3)
	exestrpre := "UPDATE piaofang.showtime set salecount= CASE typeshowindex\n"
	exestrsub := ""
	lwhere := list.New().Init()
	for k, v := range mapinfo {
		if k != 0 && v != 0 {
			exestrsub += "WHEN " + strconv.Itoa(k) + " THEN " + strconv.Itoa(v) + "\n"
			lwhere.PushBack(strconv.Itoa(k))
		}
	}
	res, err := db.Exec(exestrpre + exestrsub + "\n End Where typeshowindex In (" + strings.Join(base.StringListToArray(lwhere), ",") + ")")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	count, _ := res.RowsAffected()
	fmt.Print("update people" + strconv.Itoa(count))
	return count
}


