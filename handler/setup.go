package handler

import (
	"context"
	"fmt"

	aw "github.com/deanishe/awgo"
	"golang.org/x/oauth2"

	"github.com/cate1016/alfred-timelog/alfred"
	"github.com/cate1016/alfred-timelog/api"
)

func DoSetup(wf *aw.Workflow, _ []string) (string, error) {
	token, err := alfred.GetToken(wf)
	if err != nil {
		return "", fmt.Errorf("please authorize Timelog with `tl authorize` first 👀 (%w)", err)
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

	err = alfred.StoreOngoingTimelog(wf, alfred.Timelog{
		DriveFolderID:   did,
		SpreadsheetName: sheetname,
		SpreadsheetID:   csid,
		WeekStartUnix:   start,
		WeekEndUnix:     end,
	})
	if err != nil {
		return "", fmt.Errorf("could not save Timelog data, please try again later 🙏 (%w)", err)
	}

	return "Timelog initialize successfully ⌛", nil
}
