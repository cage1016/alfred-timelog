/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"

	"github.com/cage1016/alfred-timelog/alfred"
	"github.com/cage1016/alfred-timelog/lib"
)

var (
	av = aw.NewArgVars()
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add time log",
	Run:   runAddCmd,
}

func ErrorHandle(err error) {
	av.Var("error", err.Error())
	if err := av.Send(); err != nil {
		wf.Fatalf("failed to send args to Alfred: %v", err)
	}
}

func runAddCmd(c *cobra.Command, args []string) {
	CheckForUpdate()

	tt, _ := alfred.LoadOngoingTimelog(wf)
	loc, err := time.LoadLocation(alfred.GetTimeZone(wf))
	if err != nil {
		ErrorHandle(fmt.Errorf("failed parse Time zone: %v", err))
		return
	}
	now := time.Now().In(loc)

	dir := alfred.GetBaseFolder(wf)
	if dir == "" {
		dir = wf.DataDir()
	}

	if tt.Dir != dir || tt.FileName == "" || now.Unix() > tt.WeekEndUnix {
		var err error
		tt, err = lib.InitOrUpdate(now, dir)
		if err != nil {
			ErrorHandle(fmt.Errorf("failed init Excel: %v", err))
			return
		}
		err = alfred.StoreOngoingTimelog(wf, tt)
		if err != nil {
			ErrorHandle(fmt.Errorf("failed store ongoing timelog: %v", err))
			return
		}
	}

	fullName := filepath.Join(tt.Dir, tt.FileName)
	f, err := excelize.OpenFile(fullName)
	if err != nil {
		ErrorHandle(fmt.Errorf("failed open Excel: %v", err))
		return
	}
	defer func() {
		if err := f.Save(); err != nil {
			ErrorHandle(fmt.Errorf("failed save Excel: %v", err))
		}
	}()

	d := (now.Unix() - tt.WeekStartUnix) / 60 / 30
	ra := fmt.Sprintf("%s%d", lib.IdOf(int((d/48)+1)), (d%48)+2)
	v, _ := f.GetCellValue("Sheet1", ra)
	cs, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
	})
	f.SetCellStyle("Sheet1", ra, ra, cs)

	if v != "" {
		f.SetCellValue("Sheet1", ra, v+"\n"+args[0])
	} else {
		f.SetCellValue("Sheet1", ra, args[0])
	}

	av.Var("file", filepath.Base(tt.FileName))
	av.Var("msg", args[0])
	if err := av.Send(); err != nil {
		wf.Fatalf("failed to send args to Alfred: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
}
