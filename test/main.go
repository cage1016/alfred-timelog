package main

import (
	"context"
	"encoding/json"
	"errors"
	"os/user"

	oauth2ns "github.com/nmrshll/oauth2-noserver"
	log "github.com/sirupsen/logrus"
	keyring "github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	serviceName = "test-go-api"
)

var svc *sheets.Service

func Init() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	// ask the user to authenticate on google in the browser
	conf := &oauth2.Config{
		ClientID:     "335558396700-24f5fsevh8rtjpdc8iscdc28596k7nj5.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
			"https://www.googleapis.com/auth/drive.file",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
		},
	}
	client := &oauth2ns.AuthorizedClient{}

	// Try to use existing token
	existToken, err := retrieveToken(user.Name)
	forceToken := false
	service := &sheets.Service{}

	if err != nil || forceToken == true {
		// Token not found
		log.Debug(err)

		// Request a new access token
		client, err = oauth2ns.AuthenticateUser(conf)
		if err != nil {
			log.Debug(err)
		}

		// Store it
		storeToken(user.Name, client.Token)
	} else {
		// Use existing one
		client = &oauth2ns.AuthorizedClient{
			Client: conf.Client(context.Background(), existToken),
			Token:  existToken,
		}
	}
	service, err = sheets.New(client.Client)
	if err != nil {
		log.Fatal(err)
	}

	svc = service
}

func storeToken(googleUserEmail string, token *oauth2.Token) error {
	tokenJSONBytes, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = keyring.Set(serviceName, googleUserEmail, string(tokenJSONBytes))
	if err != nil {
		log.Debugf("failed storing token into keyring: %v", err)
		return err
	}

	return nil
}

func retrieveToken(googleUserEmail string) (*oauth2.Token, error) {
	tokenJSONString, err := keyring.Get(serviceName, googleUserEmail)
	if err != nil {
		if err == keyring.ErrNotFound {
			return nil, err
		}

		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal([]byte(tokenJSONString), &token)
	if err != nil {
		log.Debugf("failed unmarshaling token: %v", err)
		return nil, err
	}

	// validate token
	if !token.Valid() {
		return nil, errors.New("invalid token")
	}

	return &token, nil
}

func setup() {
	sheetId := "1MK8UYpgJcWmax74mxExXtjAURtn-TKmNTpLeXrp_Pmc"

	reqs := []*sheets.Request{
		// frozen column 1
		{
			UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
				Properties: &sheets.SheetProperties{
					SheetId: 0,
					GridProperties: &sheets.GridProperties{
						FrozenColumnCount: int64(1),
					},
				},
				Fields: "gridProperties.frozenColumnCount",
			},
		},
		// frozen row 1
		{
			UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
				Properties: &sheets.SheetProperties{
					SheetId: 0,
					GridProperties: &sheets.GridProperties{
						FrozenRowCount: int64(1),
					},
				},
				Fields: "gridProperties.frozenRowCount",
			},
		},
		// Header size
		{
			UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
				Range: &sheets.DimensionRange{
					SheetId:    0,
					Dimension:  "COLUMNS",
					StartIndex: 0,
					EndIndex:   8,
				},
				Properties: &sheets.DimensionProperties{
					PixelSize: 200,
				},
				Fields: "pixelSize",
			},
		},
		// First Column Row size
		{
			UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
				Range: &sheets.DimensionRange{
					SheetId:    0,
					Dimension:  "ROWS",
					StartIndex: 1,
					EndIndex:   49,
				},
				Properties: &sheets.DimensionProperties{
					PixelSize: 35,
				},
				Fields: "pixelSize",
			},
		},
		// Header size
		{
			UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
				Range: &sheets.DimensionRange{
					SheetId:    0,
					Dimension:  "COLUMNS",
					StartIndex: 0,
					EndIndex:   1,
				},
				Properties: &sheets.DimensionProperties{
					PixelSize: 100,
				},
				Fields: "pixelSize",
			},
		},
		// header size
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          0,
					StartRowIndex:    0,
					EndRowIndex:      1,
					StartColumnIndex: 0,
					EndColumnIndex:   8,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: &sheets.Color{
							Red:   0.0,
							Green: 0.0,
							Blue:  0.0,
						},
						HorizontalAlignment: "CENTER",
						TextFormat: &sheets.TextFormat{
							ForegroundColor: &sheets.Color{
								Red:   1.0,
								Green: 1.0,
								Blue:  1.0,
							},
							FontSize: int64(10),
							Bold:     true,
						},
					},
				},
				Fields: "userEnteredFormat(backgroundColor,textFormat,horizontalAlignment)",
			},
		},
	}

	reqs = append(reqs, FnColor(40, 70.0, 70.0, 70.0, 0.0, 0.0, 0.0))
	reqs = append(reqs, FnColor(16, 70.0, 70.0, 70.0, 0.0, 0.0, 0.0))
	for i := 1; i <= 49; i++ {
		o := i % 2
		if o == 1 || i == 16 || i == 40 {
			reqs = append(reqs, FnAlignment(int64(i), "LEFT"))
		} else {
			reqs = append(reqs, FnAlignment(int64(i), "RIGHT"))
		}
	}

	_, err := svc.Spreadsheets.BatchUpdate(sheetId, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: reqs,
	}).Do()
	if err != nil {
		log.Fatalf("Spreadsheets.BatchUpdate:", err.Error())
	}
}

func FnColor(rowIndex int64, br, bg, bb, fr, fg, fb float64) *sheets.Request {
	return &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Range: &sheets.GridRange{
				SheetId:          0,
				StartRowIndex:    rowIndex - 1,
				EndRowIndex:      rowIndex,
				StartColumnIndex: 0,
				EndColumnIndex:   8,
			},
			Cell: &sheets.CellData{
				UserEnteredFormat: &sheets.CellFormat{
					BackgroundColor: &sheets.Color{
						Red:   br,
						Green: bg,
						Blue:  bb,
					},
					TextFormat: &sheets.TextFormat{
						ForegroundColor: &sheets.Color{
							Red:   fr,
							Green: fg,
							Blue:  fb,
						},
					},
				},
			},
			Fields: "userEnteredFormat(backgroundColor,textFormat)",
		},
	}
}

func FnAlignment(rowIndex int64, Alignment string) *sheets.Request {
	return &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Range: &sheets.GridRange{
				SheetId:          0,
				StartRowIndex:    rowIndex - 1,
				EndRowIndex:      rowIndex,
				StartColumnIndex: 0,
				EndColumnIndex:   1,
			},
			Cell: &sheets.CellData{
				UserEnteredFormat: &sheets.CellFormat{
					HorizontalAlignment: Alignment,
				},
			},
			Fields: "userEnteredFormat(horizontalAlignment)",
		},
	}
}

func main() {
	Init()
	setup()
}
