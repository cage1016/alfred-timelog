package handler

import (
	"fmt"

	aw "github.com/deanishe/awgo"

	"github.com/cate1016/timetrack-alfred-workflows/alfred"
	"github.com/cate1016/timetrack-alfred-workflows/api"
)

func DoDeAuthorize(wf *aw.Workflow, _ []string) (string, error) {
	token, err := alfred.GetToken(wf)
	if err != nil {
		return "", fmt.Errorf("already deauthorized 👀 (%w)", err)
	}

	if err := api.RevokeToken(token); err != nil {
		return "", fmt.Errorf("error during deauthorization, please try again later 🙏 (%w)", err)
	}

	// nolint:errcheck // removing the token from the keychain is best effort
	_ = alfred.RemoveToken(wf)

	return "Timetrack deauthorized successfully 😎", nil
}
