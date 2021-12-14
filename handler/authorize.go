package handler

import (
	"encoding/json"
	"fmt"

	aw "github.com/deanishe/awgo"

	"github.com/cate1016/alfred-timelog/alfred"
	"github.com/cate1016/alfred-timelog/api"
)

func DoAuthorize(wf *aw.Workflow, _ []string) (string, error) {
	clientID := alfred.GetClientID(wf)
	clientSecret := alfred.GetClientSecret(wf)

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("client id or client secret is empty 👀")
	}

	token, err := api.GetToken(api.NewConfig(clientID, clientSecret))
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
