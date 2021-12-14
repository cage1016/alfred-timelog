package api_test

import (
	"context"
	"testing"

	"golang.org/x/oauth2"
	"gotest.tools/assert"

	"github.com/golang/mock/gomock"

	"github.com/cate1016/alfred-timelog/api"
)

func Test_IdOf(t *testing.T) {
	tests := []struct {
		id     int
		expect string
	}{
		{0, "A"},
		{26, "AA"},
		{27, "AB"},
		{500, "SG"},
	}

	for _, tt := range tests {
		id := api.IdOf(tt.id)
		assert.Equal(t, tt.expect, id)
	}
}

func Test_TimelogSheetInitialize(t *testing.T) {
	type fields struct {
	}

	type args struct {
		sheetId string
		wds     []string
		tr      []string
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(err error)
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: "GITHUB_TOKEN"},
			)
			tc := oauth2.NewClient(ctx, ts)
			svc, err := api.NewSheet(tc)
			if err != nil {
				t.Fatal(err)
			}

			if err = svc.TimelogSheetInitialize(tt.args.sheetId, tt.args.wds, tt.args.tr); err != nil {
				t.Errorf("svc.TimelogSheetInitialize error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(err)
				}
			}
		})
	}

}
