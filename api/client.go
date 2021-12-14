package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/snabb/isoweek"
)

const (
	PassContextValueKey = "Values"
)

func GetWeekDay(t time.Time, loc *time.Location) []string {
	wyear, week := isoweek.FromDate(t.Year(), t.Month(), t.Day())
	c := isoweek.StartTime(wyear, week, loc)
	buf := make([]string, 7)
	for i := 0; i < 7; i++ {
		buf[i] = c.AddDate(0, 0, i).Format("01/02, Monday")
	}
	return buf
}

func GetTimeRange() []string {
	buf := []string{
		"5",
		"5:30",
		"6",
		"6:30",
		"7",
		"7:30",
		"8",
		"8:30",
		"9",
		"9:30",
		"10",
		"10:30",
		"11",
		"11:30",
		"12PM",
		"12:30",
		"1",
		"1:30",
		"2",
		"2:30",
		"3",
		"3:30",
		"4",
		"4:30",
		"5",
		"5:30",
		"6",
		"6:30",
		"7",
		"7:30",
		"8",
		"8:30",
		"9",
		"9:30",
		"10",
		"10:30",
		"11",
		"11:30",
		"12AM",
		"12:30",
		"1",
		"1:30",
		"2",
		"2:30",
		"3",
		"3:30",
		"4",
		"4:30",
	}
	return buf
}

type Client struct {
	drive DriveService
	sheet SheetService
}

func (s *Client) Setup(drivefolderName string) (did, sheetname, csid string, start, end int64, err error) {
	// check drive folder
	did, err = s.drive.GetOrCreate(drivefolderName)
	if err != nil {
		return "", "", "", 0, 0, err
	}

	// get sheetname
	now := time.Now()
	wyear, week := isoweek.FromDate(now.Year(), now.Month(), now.Day())

	sheetname = fmt.Sprintf("%dw%d", wyear, week)
	wds := GetWeekDay(now, now.Location())
	tr := GetTimeRange()

	c := isoweek.StartTime(wyear, week, now.Location())
	a := time.Date(c.Year(), c.Month(), c.Day(), 5, 0, 0, 0, now.Location())
	b := a.AddDate(0, 0, 7)
	start = a.Unix()
	end = b.Unix()

	// create spreadsheet
	csid, err = s.drive.GetOrCreateSheet(sheetname, did)
	if err != nil {
		return "", "", "", 0, 0, err
	}

	// update spreadsheet data and format
	err = s.sheet.TimelogSheetInitialize(csid, wds, tr)
	if err != nil {
		return "", "", "", 0, 0, err
	}
	return
}

func (s *Client) AppendField(WeekStartUnix int64, ssid, content string) (err error) {
	now := time.Now()
	d := (now.Unix() - WeekStartUnix) / 60 / 30
	return s.sheet.AppendField(ssid, fmt.Sprintf("%s%d", IdOf(int((d/48)+1)), (d%48)+2), content)
}

func NewClient(client *http.Client) (ks *Client, err error) {
	drive, err := NewDrive(client)
	if err != nil {
		return nil, err
	}
	sheet, err := NewSheet(client)
	if err != nil {
		return nil, err
	}

	ks = &Client{
		drive: drive,
		sheet: sheet,
	}
	return
}
