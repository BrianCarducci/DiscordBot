package weather

import (
	"io/ioutil"
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
type GeoCoordinatesResponse struct {
	Latitude float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Status string `json:"status"`
}

type Location struct {
	Latitude float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}


func (wb* GeoLocator) GetGeoCoordinates(location string) (GeoCoordinatesResponse, error) {
	values := url.Values{}
	values.Add("address", location)
	values.Add("key", wb.Token)

	googleMapsBaseResponse, err := httpClient.Get("https://maps.googleapis.com/maps/api/geocode/json?" + values.Encode())
	if err != nil {
		fmt.Println("err in get req")
		return GeoCoordinatesResponse{}, err
	}
	defer googleMapsBaseResponse.Body.Close()

	var geoCoordinatesResponseMap map[string]interface{}
	json.NewDecoder(googleMapsBaseResponse.Body).Decode(&geoCoordinatesResponseMap)

	if geoCoordinatesResponseMap["status"] != "OK" {
		fmt.Println("not ok")
		return GeoCoordinatesResponse{}, errors.New("Google maps error: " + geoCoordinatesResponseMap["status"].(string))
	}

	geoCoordinatesResponse := GeoCoordinatesResponse{
		Latitude: geoCoordinatesResponseMap["results"].([]interface{})[0]["geometry"]["location"]["lat"].(float32),
		Longitude: geoCoordinatesResponseMap["results"][0]["geometry"]["location"]["lng"].(float32),

	}

	return geoCoordinatesResponse, nil
}

func (wb* GeoLocator) TestGet(location []string) (string, error) {
	res, err := wb.GetGeoCoordinates(location[0])
	fmt.Println("status: " + res.Status)
	return fmt.Sprintf("%f,%f", res.Latitude, res.Longitude), err
}