package handler

import (
	"context"
	"fmt"

	aw "github.com/deanishe/awgo"
	"golang.org/x/oauth2"

	"github.com/cate1016/alfred-timetrack/alfred"
	"github.com/cate1016/alfred-timetrack/api"
)

func DoSetup(wf *aw.Workflow, _ []string) (string, error) {
	token, err := alfred.GetToken(wf)
	if err != nil {
		return "", fmt.Errorf("please authorize TimeTracker with `tt authorize` first 👀 (%w)", err)
	}

	ctx := context.Background()
	client, err := api.NewClient(
		oauth2.NewClient(ctx, api.NewConfig(alfred.GetClientID(wf), alfred.GetClientSecret(wf)).TokenSource(ctx, token)),
	)
	if err != nil {
		return "", fmt.Errorf("something wrong happened, please try again later 🙏 (%w)", err)
	}

	folderName := alfred.GetDriveFolderName(wf)
	did, sheetname, csid, start, end, err := client.Setup(folderName)
	if err != nil {
		return "", fmt.Errorf("could not init setup, please try again later 🙏 (%w)", err)
	}

	err = alfred.StoreOngoingTimetrack(wf, alfred.Timetrack{
		DriveFolderID:   did,
		SpreadsheetName: sheetname,
		SpreadsheetID:   csid,
		WeekStartUnix:   start,
		WeekEndUnix:     end,
	})
	if err != nil {
		return "", fmt.Errorf("could not save timetrack data, please try again later 🙏 (%w)", err)
	}

	return "Timetrack initialize successfully ⌛", nil
}
