package handler

import (
	"encoding/json"
	"fmt"

	aw "github.com/deanishe/awgo"

	"github.com/cate1016/timetrack-alfred-workflows/alfred"
	"github.com/cate1016/timetrack-alfred-workflows/api"
)

func DoAuthorize(wf *aw.Workflow, _ []string) (string, error) {
	config := api.NewConfig(alfred.GetClientID(wf))

	token, err := api.GetToken(config)
	if err != nil {
		return "", fmt.Errorf("cannot get an access token 😢 (%w)", err)
	}

	b, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("cannot serialize the token to JSON 😢 (%w)", err)
	}

	if err := alfred.SetToken(wf, string(b)); err != nil {
		return "", fmt.Errorf("cannot store the token in the keychain 😢 (%w)", err)
	}

	return "Token stored successfully 😎", nil
}
