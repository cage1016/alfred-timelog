package lib

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cage1016/alfred-timelog/alfred"
	template "github.com/cage1016/alfred-timelog/templates"
	"github.com/xuri/excelize/v2"
)

func IdOf(i int) string {
	var b string
	if i >= 26 {
		b = IdOf((i / 26 >> 0) - 1)
	} else {
		b = ""
	}
	return b + string("ABCDEFGHIJKLMNOPQRSTUVWXYZ"[i%26>>0])
}

func InitOrUpdate(now time.Time, dir string) (res alfred.Timelog, err error) {
	te := template.NewEngine()
	sn, wds, start, end := GetNewStartEndUnix(now)

	var f *excelize.File
	fileName := fmt.Sprintf("%s.xlsx", sn)
	fullName := filepath.Join(dir, fmt.Sprintf("%s.xlsx", sn))
	if _, err := os.Stat(fullName); errors.Is(err, os.ErrNotExist) {
		data := te.MustAssetString("xlsx/template.xlsx")
		f, err = excelize.OpenReader(bytes.NewReader([]byte(data)))
		if err != nil {
			return res, err
		}
	} else {
		f, err = excelize.OpenFile(fullName)
		if err != nil {
			return res, err
		}
	}

	nwds := append([]string{sn}, wds...)
	for i, v := range nwds {
		f.SetCellValue("Sheet1", fmt.Sprintf("%s1", IdOf(i)), v)
	}

	path := filepath.Join(dir, fileName)
	if err := f.SaveAs(path); err != nil {
		return res, err
	}

	res = alfred.Timelog{
		Dir:           dir,
		FileName:      fileName,
		WeekStartUnix: start,
		WeekEndUnix:   end,
	}
	return
}
