package alfred

import (
	"fmt"

	aw "github.com/deanishe/awgo"
)

const ongoingTimelog = "timelog.json"

type Timelog struct {
	DriveFolderID   string `json:"drive_folder_id"`
	SpreadsheetID   string `json:"spreadsheet_id"`
	SpreadsheetName string `json:"spreadsheet_name"`
	WeekStartUnix   int64  `json:"week_start_unix"`
	WeekEndUnix     int64  `json:"week_end_unix"`
}

func LoadOngoingTimelog(wf *aw.Workflow) (Timelog, error) {
	// fallback load function doing nothing
	nop := func() (interface{}, error) {
		return Timelog{}, nil
	}

	var timelog Timelog
	if err := wf.Data.LoadOrStoreJSON(ongoingTimelog, 0, nop, &timelog); err != nil {
		return Timelog{}, fmt.Errorf("error loading the ongoing timelog: %w", err)
	}

	return timelog, nil
}

func StoreOngoingTimelog(wf *aw.Workflow, timelog Timelog) error {
	if err := wf.Data.StoreJSON(ongoingTimelog, timelog); err != nil {
		return fmt.Errorf("error storing the ongoing timelog: %w", err)
	}

	return nil
}
