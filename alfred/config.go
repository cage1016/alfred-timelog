package alfred

import (
	aw "github.com/deanishe/awgo"
)

const (
	clientID        = "client_id"
	clientSecret    = "client_secret"
	DriveFolderName = "drive_folder_name"
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
