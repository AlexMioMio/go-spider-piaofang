package wangpiao

import (
	"testing"
)

func Test_GetMovie(t *testing.T) {
	GetMovie()
}

func Test_saveTodayShowTime(t *testing.T) {
	saveTodayShowTime()

}
func Test_updateAllShowTimePeople(t *testing.T) {
	updateAllShowTimePeople()
}

func Test_GetCinema(t *testing.T) {
	GetCinema(getCitysFromDB())
}
func Test_getCinemaSingleCity(t *testing.T) {
	getCinemaSingleCity(1)
}

func Test_getShowTimeSingleCinema(t *testing.T) {
	getShowTimeSingleCinema(1156, "2015-10-16")
}