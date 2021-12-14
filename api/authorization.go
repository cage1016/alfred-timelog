package api

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dvsekhvalnov/jose2go/base64url"
	"github.com/google/uuid"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const port = 38146

type response struct {
	values url.Values
	err    error
}

func NewConfig(clientID, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:" + strconv.Itoa(port),
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
			"https://www.googleapis.com/auth/drive",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
		},
	}
}

func GetToken(config *oauth2.Config) (*oauth2.Token, error) {
	state := uuid.New().String()

	challengeRaw, err := randomStringURLSafe(96)
	if err != nil {
		return nil, fmt.Errorf("cannot generate a random string for the challenge: %w", err)
	}

	challengeSha256 := sha256.Sum256([]byte(challengeRaw))
	challengeURLEncoded := base64url.Encode(challengeSha256[:])

	codeChallenge := oauth2.SetAuthURLParam("code_challenge", challengeURLEncoded)
	codeChallengeMethod := oauth2.SetAuthURLParam("code_challenge_method", "S256")

	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline, codeChallenge, codeChallengeMethod)

	log.Println("open the browser and start the authorization server")

	if err := browser.OpenURL(authURL); err != nil {
		return nil, fmt.Errorf("cannot open a browser to handle the authorization flow: %w", err)
	}

	res := <-callback("127.0.0.1:" + strconv.Itoa(port))

	if errorCode := res.values.Get("error"); errorCode != "" {
		return nil, fmt.Errorf("the user did not grant the required permissions")
	}

	actual := res.values.Get("state")
	if state != actual {
		return nil, fmt.Errorf("state does not match the original one, something nasty happened")
	}

	code := res.values.Get("code")
	verifier := oauth2.SetAuthURLParam("code_verifier", challengeRaw)
	token, err := config.Exchange(context.Background(), code, verifier)

	if err != nil {
		return nil, fmt.Errorf("cannot exchange the OAuth 2 code for an access token: %w", err)
	}

	return token, nil
}

// nolint:gocognit // this function is only slightly more complex than the allowed threshold
func callback(address string) chan *response {
	responseCh, shutdownCh, interruptCh := make(chan *response), make(chan bool), make(chan bool)
	server := &http.Server{Addr: address}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		var msg string
		if r.URL.Query().Get("code") != "" {
			msg = "Timelog authenticated correctly, you can now close this window."
		} else {
			msg = "Something went wrong with the authorization workflow. Please try again."
		}

		if _, err := w.Write([]byte(msg)); err != nil {
			log.Printf("http.ResponseWriter write failed: %v", err)
		}

		interruptCh <- true
		responseCh <- &response{values: r.URL.Query()}
		shutdownCh <- true
	})

	// run the server
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error running the authorization server: %s\n", err)
		}
	}()

	// shutdown the server after a timeout
	go func() {
		select {
		case <-interruptCh:
		case <-time.After(10 * time.Minute):
			responseCh <- &response{err: fmt.Errorf("timeout to complete the authorization flow expired")}
			shutdownCh <- false
		}
	}()

	// shutdown the server gracefully
	go func() {
		done := <-shutdownCh

		if done {
			log.Println("authorization flow done, shutting down the authorization server")
		} else {
			log.Println("timeout to done the authorization flow expired, shutting down the HTTP server")
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("authorization server could not shutdown gracefully: %v", err)
		}
	}()

	return responseCh
}

func RevokeToken(token *oauth2.Token) error {
	data := url.Values{"token": {token.RefreshToken}}
	resp, err := http.PostForm("https://accounts.google.com/o/oauth2/revoke", data)

	if err != nil {
		return fmt.Errorf("networking error while trying to revoke the token: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("revoke endpoint returned a %d status code", resp.StatusCode)
	}

	return nil
}

type OauthState struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	IDToken     string `json:"id_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func RefreshToken(clientID, clientSecret, refresToken string) (*oauth2.Token, error) {
	formdata := url.Values{}
	formdata.Add("client_id", clientID)
	formdata.Add("grant_type", "refresh_token")
	formdata.Add("client_secret", clientSecret)
	formdata.Add("refresh_token", refresToken)
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", formdata)
	if err != nil {
		return nil, fmt.Errorf("networking error while trying to renew the token: %w", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			var self = new(OauthState)
			err = json.NewDecoder(resp.Body).Decode(self)
			if err != nil {
				return nil, err
			}

			token := &oauth2.Token{
				AccessToken:  self.AccessToken,
				TokenType:    self.TokenType,
				RefreshToken: refresToken,
				Expiry:       time.Now().Add(time.Duration(self.ExpiresIn) * time.Second),
			}
			return token, nil
		} else {
			return nil, fmt.Errorf("renew token failed, status code: %d", resp.StatusCode)
		}
	}
}
