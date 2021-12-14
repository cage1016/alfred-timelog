package alfred

import (
	"fmt"

	aw "github.com/deanishe/awgo"
)

const ongoingTimetrack = "timetrack.json"

type Timetrack struct {
	DriveFolderID   string `json:"drive_folder_id"`
	SpreadsheetID   string `json:"spreadsheet_id"`
	SpreadsheetName string `json:"spreadsheet_name"`
	WeekStartUnix   int64  `json:"week_start_unix"`
	WeekEndUnix     int64  `json:"week_end_unix"`
}

func LoadOngoingTimetrack(wf *aw.Workflow) (Timetrack, error) {
	// fallback load function doing nothing
	nop := func() (interface{}, error) {
		return Timetrack{}, nil
	}

	var timetrack Timetrack
	if err := wf.Data.LoadOrStoreJSON(ongoingTimetrack, 0, nop, &timetrack); err != nil {
		return Timetrack{}, fmt.Errorf("error loading the ongoing timetrack: %w", err)
	}

	return timetrack, nil
}

func StoreOngoingTimetrack(wf *aw.Workflow, timetrack Timetrack) error {
	if err := wf.Data.StoreJSON(ongoingTimetrack, timetrack); err != nil {
		return fmt.Errorf("error storing the ongoing timetrack: %w", err)
	}

	return nil
}
