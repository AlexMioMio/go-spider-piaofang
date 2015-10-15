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

func main() {
	//	saveTodayShowTime()
	updateAllShowTimePeople()
	//	GetCinema(getCitysFromDB())
	//	getCinemaSingleCity()
	//	GetShowTime()
	//	getShowTimeSingleCinema(1156, "2015-10-16")
	//	ts, _ := time.Parse("2006-01-02 15:04:05", "2015-10-15 14:15:00")
	//fmt.Println(ts.Unix())
	//	getNowTimeTsString()
}




