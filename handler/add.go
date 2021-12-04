package handler

import (
	"context"
	"fmt"
	"log"

	aw "github.com/deanishe/awgo"
	"golang.org/x/oauth2"

	"github.com/cate1016/timetrack-alfred-workflows/alfred"
	"github.com/cate1016/timetrack-alfred-workflows/api"
)

func Add(wf *aw.Workflow, args []string) (string, error) {
	token, err := alfred.GetToken(wf)
	if err != nil {
		return "", fmt.Errorf("please authorize TimeTracker with `tt authorize` first 👀 (%w)", err)
	}

	log.Print("[DEBUG] token: ", token.AccessToken)
	log.Print("[DEBUG] token: ", token.RefreshToken)
	log.Print("[DEBUG] token: ", token.Expiry)
	log.Print("[DEBUG] token: ", token.TokenType)

	ctx := context.Background()
	clientID := alfred.GetClientID(wf)
	client, err := api.NewClient(ctx, oauth2.NewClient(ctx, api.NewConfig(clientID).TokenSource(ctx, token)))
	if err != nil {
		return "", fmt.Errorf("something wrong happened, please try again later 🙏 (%w)", err)
	}

	weekStartUnix := alfred.GetWeekStartUnix(wf)
	ssid := alfred.GetSpreadsheetID(wf)
	if err != client.AppendField(int64(weekStartUnix), ssid, args[0]) {
		return "", fmt.Errorf("could not Add action description to Spreadsheet, please try again later 🙏 (%w)", err)
	}

	return fmt.Sprintf("Timetrack add '%s'to SpreadSheet successfully 😎", args[0]), nil
}
