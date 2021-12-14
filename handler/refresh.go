package handler

import (
	"encoding/json"
	"fmt"

	aw "github.com/deanishe/awgo"

	"github.com/cate1016/alfred-timelog/alfred"
	"github.com/cate1016/alfred-timelog/api"
)

func DoRefresh(wf *aw.Workflow, _ []string) (string, error) {
	token, err := alfred.GetToken(wf)
	if err != nil {
		return "", fmt.Errorf("please authorize Timelog with `tl authorize` first 👀 (%v)", err)
	}

	token, err = api.RefreshToken(alfred.GetClientID(wf), alfred.GetClientSecret(wf), token.RefreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to renew token: %v", err)
	}

	b, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("cannot serialize the token to JSON 😢 (%w)", err)
	}

	if err := alfred.SetToken(wf, string(b)); err != nil {
		return "", fmt.Errorf("cannot store the token in the keychain 😢 (%w)", err)
	}

	return "Do refresh token successfully 👍", nil
}
