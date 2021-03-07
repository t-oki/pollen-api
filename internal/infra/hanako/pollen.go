package hanako

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/sclevine/agouti"
	"github.com/t-oki/pollen-api/internal/domain/entity"
)

type PollenRepositoryImpl struct{}

func NewPollenRepositoryImpl() entity.PollenRepository {
	return &PollenRepositoryImpl{}
}

func (r *PollenRepositoryImpl) FetchPollen(area entity.Area, observatory entity.Observatory, from, to time.Time) ([]entity.Pollen, error) {
	fromYear, fromMonth, fromDay, fromHour := from.Year(), from.Month(), from.Day(), from.Hour()
	toYear, toMonth, toDay, toHour := to.Year(), to.Month(), to.Day(), to.Hour()
	if from.Hour() == 0 {
		fromMinus1 := from.AddDate(0, 0, -1)
		fromYear, fromMonth, fromDay, fromHour = fromMinus1.Year(), fromMinus1.Month(), fromMinus1.Day(), 24
	}
	if to.Hour() == 0 {
		toMinus1 := to.AddDate(0, 0, -1)
		toYear, toMonth, toDay, toHour = toMinus1.Year(), toMinus1.Month(), toMinus1.Day(), 24
	}
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
		log.Errorf("Failed to Send: %v", err)
	}
	if err := page.Navigate("http://kafun.taiki.go.jp/DownLoad1.aspx"); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("ddlStartYear").Select(strconv.Itoa(fromYear)); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("ddlStartMonth").Select(strconv.Itoa(int(fromMonth))); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("ddlStartDay").Select(strconv.Itoa(fromDay)); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("ddlStartHour").Select(strconv.Itoa(fromHour)); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("ddlEndYear").Select(strconv.Itoa(toYear)); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("ddlEndMonth").Select(strconv.Itoa(int(toMonth))); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("ddlEndDay").Select(strconv.Itoa(toDay)); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("ddlEndHour").Select(strconv.Itoa(toHour)); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("ddlArea").Select(area.Name); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FirstByName(fmt.Sprintf("CheckBoxMstList$%d", observatory.ID-1)).Click(); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := page.FindByID("download").Click(); err != nil {
		return nil, errors.WithStack(err)
	}
	// ダウンロードが完了するまでセッションを確保する
	time.Sleep(time.Millisecond * 100)

	file, err := os.Open("Data.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	res := make([]entity.Pollen, 0)
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		var datetime time.Time
		if line[3] == "24" {
			datetime, err = time.Parse("2006010215", fmt.Sprintf("%02s", line[2]+"23"))
			if err != nil {
				log.Warn(err.Error())
			}
			datetime = datetime.Add(time.Hour)
		} else {
			datetime, err = time.Parse("2006010215", fmt.Sprintf("%02s", line[2]+line[3]))
			if err != nil {
				log.Warn(err.Error())
			}
		}

		pollenCount, err := strconv.ParseInt(line[10], 0, 64)
		if err != nil {
			log.Warn(err.Error())
		}
		if pollenCount < 0 {
			pollenCount = 0
		}
		windSpeed, err := strconv.ParseInt(line[12], 0, 64)
		if err != nil {
			log.Warn(err.Error())
		}
		temperature, err := strconv.ParseFloat(line[13], 64)
		if err != nil {
			log.Warn(err.Error())
		}
		rainfall, err := strconv.ParseInt(line[14], 0, 64)
		if err != nil {
			log.Warn(err.Error())
		}
		res = append(res, entity.Pollen{
			Datetime:      datetime,
			PollenCount:   &pollenCount,
			WindDirection: &line[11],
			WindSpeed:     &windSpeed,
			Temperature:   &temperature,
			Rainfall:      &rainfall,
		})
	}

	return res, nil
}
