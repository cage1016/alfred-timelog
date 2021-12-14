package handler

import (
	"context"
	"fmt"
	"time"

	aw "github.com/deanishe/awgo"
	"golang.org/x/oauth2"

	"github.com/cate1016/alfred-timetrack/alfred"
	"github.com/cate1016/alfred-timetrack/api"
)

func Add(wf *aw.Workflow, args []string) (string, error) {
	token, err := alfred.GetToken(wf)
	if err != nil {
		return "", fmt.Errorf("please authorize TimeTracker with `tt authorize` first 👀 (%v)", err)
	}

	// update token if expired
	if token.Expiry.Before(time.Now()) {
		if _, err = DoRefresh(wf, []string{}); err != nil {
			return "", fmt.Errorf("Refresh token failed: 👀 (%v)", err)
		}

		token, err = alfred.GetToken(wf)
		if err != nil {
			return "", fmt.Errorf("please authorize TimeTracker with `tt authorize` first 👀 (%v)", err)
		}
	}

	ctx := context.Background()
	client, err := api.NewClient(
		oauth2.NewClient(ctx, api.NewConfig(alfred.GetClientID(wf), alfred.GetClientSecret(wf)).TokenSource(ctx, token)),
	)
	if err != nil {
		return "", fmt.Errorf("something wrong happened, please try again later 🙏 (%v)", err)
	}

	tt, err := alfred.LoadOngoingTimetrack(wf)
	if err != nil {
		return "", fmt.Errorf("something wrong happened, please try again later 🙏 (%v)", err)
	}

	// create new spreadsheet if necessary
	if time.Now().Unix() > int64(tt.WeekEndUnix) {
		if _, err := DoSetup(wf, []string{""}); err != nil {
			return "", err
		}
	}

	// append message to cell
	if err = client.AppendField(
		tt.WeekStartUnix,
		tt.SpreadsheetID,
		args[0],
	); err != nil {
		return "", fmt.Errorf("could not Add action description to Spreadsheet, please try again later 🙏 (%v)", err)
	}

	return fmt.Sprintf("Timetrack add '%s' to SpreadSheet '%s' successfully 😎", args[0], tt.SpreadsheetName), nil
}
