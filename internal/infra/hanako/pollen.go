package hanako

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sclevine/agouti"
	"github.com/t-oki/pollen-api/internal/domain/entity"
)

type PollenRepositoryImpl struct{}

func NewPollenRepositoryImpl() entity.PollenRepository {
	return &PollenRepositoryImpl{}
}

func (r *PollenRepositoryImpl) FetchPollen(observatoryID int64, from, to time.Time) error {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions(
			"args", []string{
				"--headless",
			}),
	)
	defer driver.Stop()
	if err := driver.Start(); err != nil {
		return errors.WithStack(err)
	}
	page, err := driver.NewPage()
	if err != nil {
		return errors.WithStack(err)
	}
	if err := page.Navigate("http://kafun.taiki.go.jp/DownLoad1.aspx"); err != nil {
		return errors.WithStack(err)
	}
	fmt.Println("hi")
	if err := page.FirstByName("CheckBoxMstList$1").Click(); err != nil {
		return errors.WithStack(err)
	}
	if err := page.FindByID("download").Click(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
