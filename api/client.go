package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/snabb/isoweek"

	"github.com/cate1016/timetrack-alfred-workflows/utils"
)

const (
	PassContextValueKey = "Values"
)

type Client struct {
	drive *Drive
	sheet *Sheet
	ctx   context.Context
}

func (s *Client) Setup(drivefolderName string) (did, csid string, start, end int64, err error) {
	// check drive folder
	did, err = s.drive.GetOrCreate(drivefolderName)
	if err != nil {
		return "", "", 0, 0, err
	}

	// get sheetname
	now := time.Now()
	sheetname := utils.GetWeekFileName(now)
	wds := utils.GetWeekDay(now, now.Location())
	tr := utils.GetTimeRange()

	wyear, week := isoweek.FromDate(now.Year(), now.Month(), now.Day())

	c := isoweek.StartTime(wyear, week, now.Location())
	f := c.AddDate(0, 0, -1)
	a := time.Date(f.Year(), f.Month(), f.Day(), 5, 0, 0, 0, now.Location())
	b := a.AddDate(0, 0, 7)
	start = a.Unix()
	end = b.Unix()

	// create car model sheet
	csid, err = s.drive.GetOrCreateSheet(sheetname, did)
	if err != nil {
		return "", "", 0, 0, err
	}

	// create sheet
	err = s.sheet.TimetrackSheetInitialize(csid, wds, tr)
	if err != nil {
		return "", "", 0, 0, err
	}
	return
}

func (s *Client) AppendField(WeekStartUnix int64, ssid, content string) (err error) {
	now := time.Now()
	d := (now.Unix() - WeekStartUnix) / 60 / 30
	return s.sheet.AppendField(ssid, fmt.Sprintf("%s%d", utils.IdOf(int((d/48)+1)), (d%48)+2), content)
}

func NewClient(ctx context.Context, client *http.Client) (ks *Client, err error) {
	drive, err := NewDrive(ctx, client)
	if err != nil {
		return nil, err
	}
	sheet, err := NewSheet(ctx, client)
	if err != nil {
		return nil, err
	}

	ks = &Client{
		ctx:   ctx,
		drive: drive,
		sheet: sheet,
	}
	return
}
