package wangpiao

import (
	"../base"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"container/list"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
)

type TMovie struct {
	base.BaseMovie
	tid string
}

func GetMovie() *list.List {
	var lmovies = list.New().Init()
	resp, err := http.Get("http://www.wangpiao.com/Movie/movies/")
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		panic(err)
	}
	htmlbody := doc.Find("body")
	maindiv := htmlbody.Find("div[id=filmon]")
	divs := maindiv.Find("div[class=mt20\\ pb20\\ bodbCCC]")
	divs.Each(func(num int, s *goquery.Selection) {
		img := s.Find("div[class=movie_bg\\ movie_pic\\ pr]").Children().First().Children()
		movie := new(TMovie)
		movie.Name = img.AttrOr("title", "")
		movie.Jpg = img.AttrOr("src", "")
		lmovies.PushBack(movie)
	})
	return lmovies
}

func InsertMovieList(l *list.List) {
	base.InsertList(l, "insert into movie values ", func(any interface{}) string {
		movie := any.(*TMovie)
		return "(0,'" + movie.Name + "','" + movie.Jpg + "')"
	})
}

func getHasNames(rows *sql.Rows) []string {
	lhasnames := list.New().Init()
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		lhasnames.PushBack(name)
	}
	return base.StringListToArray(lhasnames)
}

func getNoHasNames(oldl *list.List, hasnames []string) []string {
	var linsertmovienames = list.New().Init()
	for e := oldl.Front(); e != nil; e = e.Next() {
		movie := e.Value.(*TMovie)
		isinsert := true
		for j := 0; j < len(hasnames); j++ {
			if movie.Name == hasnames[j] {
				isinsert = false
				break
			}
		}
		if isinsert == true {
			linsertmovienames.PushBack(movie.Name)
		}
	}
	return base.StringListToArray(linsertmovienames)
}

func getMovieByNames(oldl *list.List, names []string) *list.List {
	var linsertmovie = list.New().Init()
	for i := 0; i < len(names); i++ {
		movie := getMovieByName(oldl, names[i])
		if movie != nil {
			linsertmovie.PushBack(movie)
		}
	}
	fmt.Println(linsertmovie.Len())
	return linsertmovie
}
func getMovieByName(oldl *list.List, name string) *TMovie {
	for e := oldl.Front(); e != nil; e = e.Next() {
		movie := e.Value.(*TMovie)
		if movie.Name == name {
			return movie
		}
	}
	return nil
}

func saveMovie(l *list.List) {
	db, err := sql.Open("mysql", "root:1CUI@/piaofang")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	names := make([]string, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(*TMovie).Name)
		names[i] = "\"" + e.Value.(*TMovie).Name + "\""
		i++
	}
	namestrs := strings.Join(names, ",")
	stmtOut, err := db.Prepare("SELECT name FROM movie WHERE name in (" + namestrs + ")")
	fmt.Print("SELECT name FROM movie WHERE name in (" + namestrs + ")")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	rows, err := stmtOut.Query()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	hasnames := getHasNames(rows)
	fmt.Println("has count" + strconv.Itoa(len(hasnames)))
	nohasnames := getNoHasNames(l, hasnames)
	fmt.Println("no has count" + strconv.Itoa(len(nohasnames)))
	insertmovies := getMovieByNames(l, nohasnames)
	fmt.Printf(strconv.Itoa(insertmovies.Len()))
	InsertMovieList(insertmovies)
}
