package alfred

import (
	aw "github.com/deanishe/awgo"
)

const (
	timeZone   = "TIME_ZONE"
	BaseFolder = "TIMELOG_BASE_FOLDER"
)

func GetTimeZone(wf *aw.Workflow) string {
	return wf.Config.Get(timeZone, "")
}

func GetBaseFolder(wf *aw.Workflow) string {
	return wf.Config.Get(BaseFolder, "")
}
