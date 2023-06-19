package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	apiSequenceUrl      = "https://api.openstreetcam.org/2.0/sequence/?region=%d&sequenceStatus=public&sequenceType=photo&orderBy=dateAdded&orderDirection=desc&itemsPerPage=1&page=%d"
	apiSequencePhotoUrl = "https://api.openstreetcam.org/2.0/sequence/%s/photos"
)

type RegionsResponse struct {
	Result struct {
		Data []struct {
			Id          string `json:"id"`
			CountryCode string `json:"countryCode"`
		} `json:"data"`
	} `json:"result"`
}

type SequenceResponse struct {
	Result struct {
		Data []struct {
			FileUrl string `json:"fileurlLTh"`
		} `json:"data"`
	} `json:"result"`
	CountryCode string
}

type Region struct {
	code    int
	maxPage int
}

func getRandomRegion() Region {
	regions := []Region{
		{code: 38, maxPage: 2500},  // Central America
		{code: 34, maxPage: 5500},  // South America
		{code: 32, maxPage: 60000}, // Europe
		{code: 33, maxPage: 2500},  // Africa
		{code: 35, maxPage: 13000}, // Asia
	}

	rand.Seed(time.Now().Unix())
	randomRegionCode := regions[rand.Intn(len(regions))]

	return randomRegionCode
}

func getRandomRegionSequences() RegionsResponse {
	return getRegionSequences(getRandomRegion())
}

func getRegionSequences(region Region) RegionsResponse {
	randomPage := rand.Intn(region.maxPage)
	regionSequencesUrl := fmt.Sprintf(apiSequenceUrl, region.code, randomPage)
	var regionSequences RegionsResponse

	err := fetchApi(regionSequencesUrl, &regionSequences)
	if err != nil {
		log.Fatal(err)
	}

	if len(regionSequences.Result.Data) < 1 {
		log.Fatal("Region sequences data length is < 1")
	}

	return regionSequences
}

func getSequence(regionSequences RegionsResponse) SequenceResponse {
	var sequence SequenceResponse
	var regSequence = regionSequences.Result.Data[rand.Intn(len(regionSequences.Result.Data))]
	err := fetchApi(fmt.Sprintf(apiSequencePhotoUrl, regSequence.Id), &sequence)
	if err != nil {
		log.Fatal(err)
	}
	sequence.CountryCode = regSequence.CountryCode
	return sequence
}

func fetchApi(url string, obj interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}

	return nil
}

func GetRandomPhoto() SequenceResponse {
	regionSequences := getRandomRegionSequences()
	return getSequence(regionSequences)
}
