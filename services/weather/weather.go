package weather

import (
  "encoding/json"
  "errors"
  "fmt"
  "net/http"
  "net/url"
  "strings"
  "time"
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

// weatherURLData is a struct which maps to the weather.gov request for weather URL based on coordinates
type weatherURLData struct {
  Properties struct{
    ForecastURL string `json:"forecastHourly"`
  } `json:"properties"`
  Status int `json:"status"`
}

type forecastData struct {
  Properties struct{
    Periods []struct{
      Temperature     int `json:"temperature"`
      TemperatureUnit string `json:"temperatureUnit"`
      WindSpeed       string `json:"windSpeed"`
      WindDirection   string `json:"windDirection"`
      IconURL         string `json:"icon"`
      ShortForecast   string `json:"shortForecast"`
    } `json:"periods"`
  }   `json:"properties"`
  Status int `json:"status"`
}

type Location struct {
  Latitude  float32 `json:"latitude"`
  Longitude float32 `json:"longitude"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}


func (wb* GeoLocator) GetWeather(locationTokens []string) (string, error) {
  joinedLoc := strings.Join(locationTokens, " ")
  coords, err := wb.GetGeoCoordinates(joinedLoc)
  if err != nil {
    return "", err
  }

  forecastURL, err := GetForecastURL(coords)
  if err != nil {
    return "", err
  }

  forecastData, err := GetForecastData(forecastURL)
  if err != nil {
    return "", err
  }

  forecast := forecastData.Properties.Periods[0]

  msg := fmt.Sprintf("Temp: %d°%s\nWind speed: %s %s\nDescription: %s", forecast.Temperature, forecast.TemperatureUnit, forecast.WindSpeed, forecast.WindDirection, forecast.ShortForecast)
  return msg, nil
}

func (wb* GeoLocator) GetGeoCoordinates(location string) (Location, error) {
  values := url.Values{}
  values.Add("address", location)
  values.Add("key", wb.Token)

  googleMapsBaseResponse, err := httpClient.Get("https://maps.googleapis.com/maps/api/geocode/json?" + values.Encode())
  if err != nil {
    return Location{}, err
  }
  defer googleMapsBaseResponse.Body.Close()

  var gcr geoCoordinatesResponse
  json.NewDecoder(googleMapsBaseResponse.Body).Decode(&gcr)

  if gcr.Status != "OK" {
    return Location{}, errors.New("Google maps error: " + gcr.Status)
  }

  loc := Location{
    Latitude: gcr.Results[0].Geometry.Location.Lat,
    Longitude: gcr.Results[0].Geometry.Location.Lng,
  }

  return loc, nil
}

/*
func (wb* GeoLocator) TestGet(location []string) (string, error) {
  res, err := wb.GetGeoCoordinates(location[0])
  return fmt.Sprintf("%f,%f", res.Latitude, res.Longitude), err
}
*/

func GetForecastURL(location Location) (string, error) {
  locStr := fmt.Sprintf("%f,%f", location.Latitude, location.Longitude)
  weatherResp, err := httpClient.Get("https://api.weather.gov/points/" + locStr)
  if err != nil {
    return "", err
  }
  defer weatherResp.Body.Close()

  var weatherData weatherURLData
  json.NewDecoder(weatherResp.Body).Decode(&weatherData)

  weatherDataStatus := weatherData.Status
  if weatherDataStatus != 0 {
    return "", errors.New("Forecast URL request responded with status " + string(weatherDataStatus))
  }

  forecastURL := weatherData.Properties.ForecastURL
  if forecastURL == "" {
    return forecastURL, errors.New("Forecast URL was not mapped properly")
  }

  return forecastURL, nil
}

func GetForecastData(forecastURL string) (forecastData, error) {
  forecastResp, err := httpClient.Get(forecastURL)
  if err != nil {
    return forecastData{}, err
  }
  defer forecastResp.Body.Close()

  var forecast forecastData
  json.NewDecoder(forecastResp.Body).Decode(&forecast)

  forecastDataStatus := forecast.Status
  if forecastDataStatus != 0 {
    return forecastData{}, errors.New("Forecast data request responded with status " + string(forecastDataStatus))
  }

  return forecast, nil
}