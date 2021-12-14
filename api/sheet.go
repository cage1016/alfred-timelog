package api

import (
	"fmt"
	"net/http"

	"google.golang.org/api/sheets/v4"
)

func IdOf(i int) string {
	var b string
	if i >= 26 {
		b = IdOf((i / 26 >> 0) - 1)
	} else {
		b = ""
	}
	return b + string("ABCDEFGHIJKLMNOPQRSTUVWXYZ"[i%26>>0])
}

//go:generate mockgen -destination ../mocks/sheetservice.go -package=automocks . SheetService
type SheetService interface {
	TimelogSheetInitialize(sheetId string, wds, tr []string) error
	AppendField(spreadsheetId, ra, content string) error
}

type Sheet struct {
	ssvc *sheets.Service
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

func (s *Sheet) TimelogSheetInitialize(sheetId string, wds, tr []string) error {
	wdsp := append([]string{""}, wds...)
	values := make([][]interface{}, len(wdsp))
	var r, c int
	for i := 0; i < len(wdsp); i++ {
		if i == 0 {
			// header row
			trp := append([]string{""}, tr...)
			values[i] = make([]interface{}, len(trp))
			for j := 0; j < len(trp); j++ {
				values[i][j] = trp[j]
			}
		} else {
			// time range rows
			values[i] = make([]interface{}, 1)
			values[i][0] = wdsp[i]
		}

		if b := len(values[i]); b > r {
			r = b
		}
	}
	c = len(values)

	_, err := s.ssvc.Spreadsheets.Values.Update(sheetId, fmt.Sprintf("A1:%s%d", IdOf(c), r), &sheets.ValueRange{
		MajorDimension: "COLUMNS",
		Values:         values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("update spreadsheet value fail: %s", err.Error())
	}

	// update format
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

	_, err = s.ssvc.Spreadsheets.BatchUpdate(sheetId, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: reqs,
	}).Do()
	if err != nil {
		return fmt.Errorf("update spreadsheet format fail: %s", err.Error())
	}

	return nil
}

func (s *Sheet) AppendField(spreadsheetId, ra, content string) error {
	res, err := s.ssvc.Spreadsheets.Values.Get(spreadsheetId, ra).Do()
	if err != nil {
		return err
	}

	row := make([]interface{}, 1)
	if len(res.Values) > 0 {
		row[0] = res.Values[0][0].(string) + "\n" + content
	} else {
		row[0] = content
	}

	data := make([][]interface{}, 1)
	data[0] = row

	val := sheets.ValueRange{}
	val.Values = data
	_, err = s.ssvc.Spreadsheets.Values.Update(spreadsheetId, ra, &val).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}
	return nil
}

func NewSheet(client *http.Client) (SheetService, error) {
	ssvc, err := sheets.New(client)
	if err != nil {
		return nil, err
	}

	return &Sheet{ssvc: ssvc}, nil
}
