package handler

import (
	"context"
	"fmt"

	aw "github.com/deanishe/awgo"
	"golang.org/x/oauth2"

	"github.com/cate1016/timetrack-alfred-workflows/alfred"
	"github.com/cate1016/timetrack-alfred-workflows/api"
)

func DoSetup(wf *aw.Workflow, _ []string) (string, error) {
	token, err := alfred.GetToken(wf)
	if err != nil {
		return "", fmt.Errorf("please authorize TimeTracker with `tt authorize` first 👀 (%w)", err)
	}

	ctx := context.Background()
	clientID := alfred.GetClientID(wf)
	client, err := api.NewClient(ctx, oauth2.NewClient(ctx, api.NewConfig(clientID).TokenSource(ctx, token)))
	if err != nil {
		return "", fmt.Errorf("something wrong happened, please try again later 🙏 (%w)", err)
	}

	folderName := alfred.GetDriveFolderName(wf)
	did, csid, start, end, err := client.Setup(folderName)
	if err != nil {
		return "", fmt.Errorf("could not init setup, please try again later 🙏 (%w)", err)
	}

	if err := alfred.SetDriveFolderID(wf, did); err != nil {
		return "", fmt.Errorf("cannot save the configuration in Alfred, please try again later 🙏 (%w)", err)
	}

	if err := alfred.SetSpreadsheetID(wf, csid); err != nil {
		return "", fmt.Errorf("cannot save the configuration in Alfred, please try again later 🙏 (%w)", err)
	}

	if err := alfred.SetWeekStartUnix(wf, int(start)); err != nil {
		return "", fmt.Errorf("cannot save the configuration in Alfred, please try again later 🙏 (%w)", err)
	}

	if err := alfred.SetWeekEndUnix(wf, int(end)); err != nil {
		return "", fmt.Errorf("cannot save the configuration in Alfred, please try again later 🙏 (%w)", err)
	}

	return "Timetrack created successfully ⌛", nil
}
