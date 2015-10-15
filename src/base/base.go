package base

import (
	"container/list"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"strconv"
	"time"
)

type BaseMovie struct {
	Id   int
	Name string
	Jpg  string
}
type BaseCity struct {
	Id   int
	Name string
	//        Type int
	//        TypeName string
	TypeIndex int
	//        TypeIndexS string
}

type BaseShowTime struct {
	Id   int
	Type int
	TypeName string
	TypeShowIndex int
	TypeCinemaIndex int
	TypeMovieIndex int
	TypeMovieName string
	TypeHallID int
	TypeSaleEndTime int64
	TypeSaleEndTimeS time.Time
	Price int
	SeatCount int
	TypeCityIndex int
	
	//        Type int
	//        TypeName string
	
	//        TypeIndexS string
}

type BaseCinema struct {
	Id   int
	Name string
	TypeIndex int
}

func getNowTimeTsString() string {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	//	fmt.Println(ts)
	return ts
}


func StringListToArray(l *list.List) []string {
	strs := make([]string, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		strs[i] = e.Value.(string)
		i++
	}
	return strs
}

func IntListToArray(l *list.List) []int {
	strs := make([]int, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		strs[i] = e.Value.(int)
		i++
	}
	return strs
}

func InsertList(l *list.List, execpre string, f func(interface{}) string) int64 {
	db, err := sql.Open("mysql", "root:1CUI@/piaofang")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	var insertstrs = make([]string, l.Len())
	o := 0
	if l.Len() != 0 {
		for e := l.Front(); e != nil; e = e.Next() {
			f(e.Value)
			insertstrs[o] = f(e.Value)
			o++
		}
		//		fmt.Println(execpre + strings.Join(insertstrs, ","))
		res, err := db.Exec(execpre + strings.Join(insertstrs, ","))
		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}
		count, _ := res.RowsAffected()
		return count
	}
	return 0
}


