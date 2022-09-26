package alfred

import (
	aw "github.com/deanishe/awgo"
	"github.com/joho/godotenv"
)

const (
	timeZone   = "TIME_ZONE"
	BaseFolder = "TIMELOG_BASE_FOLDER"
)

func GetTimeZone(wf *aw.Workflow) string {
	myEnv, _ := godotenv.Read()
	if ok := myEnv["alfred_workflow_timezone"]; ok != "" {
		return myEnv["alfred_workflow_timezone"]
	}
	return wf.Config.Get(timeZone, "")
}

func GetBaseFolder(wf *aw.Workflow) string {
	return wf.Config.Get(BaseFolder, "")
}
