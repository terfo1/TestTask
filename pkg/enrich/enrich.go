package enrich

import (
	"TestTask/internal/config"
	"encoding/json"
	"fmt"
	"net/http"
)

var cfg *config.Config = config.LoadYaml("config.yaml")

type Enriched struct {
	Age         int
	Gender      string
	Nationality string
}

func EnrichData(name string) (*Enriched, error) {
	type ageResp struct{ Age int }
	type genderResp struct{ Gender string }
	type natResp struct {
		Country []struct {
			CountryID   string `json:"country_id"`
			Probability float64
		}
	}

	var (
		ageData    ageResp
		genderData genderResp
		natData    natResp
	)

	if err := fetchJSON(fmt.Sprintf(cfg.URL.Age, name), &ageData); err != nil {
		return nil, err
	}
	if err := fetchJSON(fmt.Sprintf(cfg.URL.Gender, name), &genderData); err != nil {
		return nil, err
	}
	if err := fetchJSON(fmt.Sprintf(cfg.URL.Nationality, name), &natData); err != nil {
		return nil, err
	}

	nationality := ""
	if len(natData.Country) > 0 {
		nationality = natData.Country[0].CountryID
	}

	return &Enriched{
		Age:         ageData.Age,
		Gender:      genderData.Gender,
		Nationality: nationality,
	}, nil
}

func fetchJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
