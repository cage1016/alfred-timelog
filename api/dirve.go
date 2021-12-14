package api

import (
	"fmt"
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

//go:generate mockgen -destination ../mocks/driveservice.go -package=automocks . DriveService
type DriveService interface {
	GetOrCreate(fd string) (string, error)
	GetOrCreateSheet(sheetName, parentId string) (string, error)
}

type Drive struct {
	dsvc *drive.Service
}

func (s *Drive) GetOrCreate(fd string) (string, error) {
	f, err := s.dsvc.Files.List().Q(fmt.Sprintf("name='%s'", fd)).Fields(googleapi.Field("nextPageToken, files(id, name)")).Spaces("drive").Do()
	if err != nil {
		return "", fmt.Errorf("query drive API fail: %v", err)
	}

	if len(f.Files) == 0 {
		f, err := s.dsvc.
			Files.Create(&drive.File{
			Name:     fd,
			MimeType: "application/vnd.google-apps.folder"}).Fields(googleapi.Field("id")).Do()
		if err != nil {
			return "", fmt.Errorf("create drive folder fail: %v", err)
		}
		return f.Id, nil
	}

	return f.Files[0].Id, nil
}

func (s *Drive) GetOrCreateSheet(sheetName, parentId string) (string, error) {
	f, err := s.dsvc.Files.List().Q(fmt.Sprintf("name='%s' and parents in '%s'", sheetName, parentId)).Fields(googleapi.Field("nextPageToken, files(id, name)")).Spaces("drive").Do()
	if err != nil {
		return "", fmt.Errorf("query drive API fail: %v", err)
	}

	if len(f.Files) == 0 {
		f, err := s.dsvc.Files.Create(&drive.File{
			Name:     sheetName,
			MimeType: "application/vnd.google-apps.spreadsheet",
			Parents:  []string{parentId},
		}).Do()
		if err != nil {
			return "", fmt.Errorf("create spreadsheet: %v", err)
		}
		return f.Id, nil
	}
	return f.Files[0].Id, nil
}

func NewDrive(client *http.Client) (DriveService, error) {
	dsvc, err := drive.New(client)
	if err != nil {
		return nil, err
	}
	return &Drive{dsvc: dsvc}, nil
}
