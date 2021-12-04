package alfred

import (
	"strconv"

	aw "github.com/deanishe/awgo"
)

const (
	clientID        = "client_id"
	clientSecret    = "client_secret"
	DriveFolderName = "drive_folder_name"
	DriveFolderID   = "drive_folder_id"
	SpreadsheetID   = "spreadsheet_id"
	WeekStartUnix   = "week_start_unix"
	WeekEndUnix     = "week_end_unix"
)

func GetClientID(wf *aw.Workflow) string {
	return wf.Config.GetString(clientID)
}

func GetClientSecret(wf *aw.Workflow) string {
	return wf.Config.Get(clientSecret)
}

func GetDriveFolderName(wf *aw.Workflow) string {
	return wf.Config.GetString(DriveFolderName)
}

// func SetDriveFolderName(wf *aw.Workflow, name string) error {
// 	return wf.Config.Set(DriveFolderName, name, false).Do()
// }

func GetDriveFolderID(wf *aw.Workflow) string {
	return wf.Config.GetString(DriveFolderID)
}

func SetDriveFolderID(wf *aw.Workflow, id string) error {
	return wf.Config.Set(DriveFolderID, id, false).Do()
}

func GetSpreadsheetID(wf *aw.Workflow) string {
	return wf.Config.GetString(SpreadsheetID)
}

func SetSpreadsheetID(wf *aw.Workflow, id string) error {
	return wf.Config.Set(SpreadsheetID, id, false).Do()
}

func GetWeekStartUnix(wf *aw.Workflow) int {
	return wf.Config.GetInt(WeekStartUnix)
}

func SetWeekStartUnix(wf *aw.Workflow, unix int) error {
	return wf.Config.Set(WeekStartUnix, strconv.Itoa(unix), false).Do()
}

func GetWeekEndUnix(wf *aw.Workflow, unix int) int {
	return wf.Config.GetInt(WeekEndUnix)
}

func SetWeekEndUnix(wf *aw.Workflow, unix int) error {
	return wf.Config.Set(WeekEndUnix, strconv.Itoa(unix), false).Do()
}
