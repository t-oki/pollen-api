package hanako

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sclevine/agouti"
	"github.com/t-oki/pollen-api/internal/domain/entity"
)

type PollenRepositoryImpl struct{}

func NewPollenRepositoryImpl() entity.PollenRepository {
	return &PollenRepositoryImpl{}
}

func (r *PollenRepositoryImpl) FetchPollen(areaName string, observatoryID int64, from, to time.Time) ([]entity.Pollen, error) {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions(
			"args", []string{
				"--headless",
			}),
	)
	if err := driver.Start(); err != nil {
		return nil, errors.WithStack(err)
	}
	defer driver.Stop()
	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result interface{}
	session := page.Session()
	if err := session.Send("POST", "chromium/send_command", map[string]interface{}{
		"cmd": "Page.setDownloadBehavior",
		"params": map[string]string{
			"behavior":     "allow",
			"downloadPath": ".",
		},
	}, &result); err != nil {
		log.Printf("Failed to Send: %v", err)
	}
	if err := page.Navigate("http://kafun.taiki.go.jp/DownLoad1.aspx"); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("ddlArea").Select(areaName); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FirstByName("CheckBoxMstList$1").Click(); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("download").Click(); err != nil {
		return nil, errors.WithStack(err)
	}
	// ダウンロードするまでセッションを確保する
	time.Sleep(time.Millisecond * 100)

	file, err := os.Open("Data.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	// res := make([]entity.Pollen, 0)
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// shiftJIS対応
		fmt.Println(line)
	}

	return nil, nil
}
