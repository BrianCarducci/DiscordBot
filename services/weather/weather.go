package weather

import (
	"fmt"
	"encoding/json"
	"net/url"
	"net/http"
	"time"
	"errors"
)

// GeoLocator holds token and has methods for weather operations
type GeoLocator struct {
	Token string
}

// GeoCoordinatesResponse is a struct which maps to the google maps request for geocode
type geoCoordinatesResponse struct {
	Results []struct{
		Geometry struct{
			Location struct{
				Lat float32 `json:"lat"`
				Lng float32 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}

type Location struct {
	Latitude float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}


func (wb* GeoLocator) GetGeoCoordinates(location string) (Location, error) {
	values := url.Values{}
	values.Add("address", location)
	values.Add("key", wb.Token)

	googleMapsBaseResponse, err := httpClient.Get("https://maps.googleapis.com/maps/api/geocode/json?" + values.Encode())
	if err != nil {
		fmt.Println("err in get req")
		return Location{}, err
	}
	defer googleMapsBaseResponse.Body.Close()

	var gcr geoCoordinatesResponse
	json.NewDecoder(googleMapsBaseResponse.Body).Decode(&gcr)

	if gcr.Status != "OK" {
		fmt.Printf("%+v\n",gcr)
		//fmt.Println("status: " + gcr.status)
		return Location{}, errors.New("Google maps error: " + gcr.Status)
	}

	loc := Location{
		Latitude: gcr.Results[0].Geometry.Location.Lat,
		Longitude: gcr.Results[0].Geometry.Location.Lng,
	}

	return loc, nil
}

func (wb* GeoLocator) TestGet(location []string) (string, error) {
	res, err := wb.GetGeoCoordinates(location[0])
	return fmt.Sprintf("%f,%f", res.Latitude, res.Longitude), err
}