package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

type Drive struct {
	dsvc *drive.Service
	ctx  context.Context
}

func (s *Drive) GetOrCreate(fd string) (string, error) {
	f, err := s.dsvc.Files.List().Q(fmt.Sprintf("name='%s'", fd)).Fields(googleapi.Field("nextPageToken, files(id, name)")).Spaces("drive").Do()
	if err != nil {
		log.Fatalf("GetOrInsertDriveFolder fail:", err.Error())
		return "", err
	}

	if len(f.Files) == 0 {
		f, err := s.dsvc.
			Files.Create(&drive.File{
			Name:     fd,
			MimeType: "application/vnd.google-apps.folder"}).Fields(googleapi.Field("id")).Do()
		if err != nil {
			log.Fatalf("GetOrInsertDriveFolder fail:", err.Error())
			return "", err
		}
		return f.Id, nil
	}

	return f.Files[0].Id, nil
}

func (s *Drive) GetOrCreateSheet(sheetName, parentId string) (string, error) {
	f, err := s.dsvc.Files.List().Q(fmt.Sprintf("name='%s' and parents in '%s'", sheetName, parentId)).Fields(googleapi.Field("nextPageToken, files(id, name)")).Spaces("drive").Do()
	if err != nil {
		log.Fatalf("GetOrCreateCarModelSpreadsheet fail:", err.Error())
		return "", err
	}

	if len(f.Files) == 0 {
		f, err := s.dsvc.Files.Create(&drive.File{
			Name:     sheetName,
			MimeType: "application/vnd.google-apps.spreadsheet",
			Parents:  []string{parentId},
		}).Do()
		if err != nil {
			log.Fatalf("GetOrInsertDriveFolder fail:", err.Error())
			return "", err
		}
		return f.Id, nil
	}
	return f.Files[0].Id, nil
}

func NewDrive(ctx context.Context, client *http.Client) (ks *Drive, err error) {
	dsvc, err := drive.New(client)
	if err != nil {
		return nil, err
	}
	ks = &Drive{
		dsvc: dsvc,
	}
	return
}
