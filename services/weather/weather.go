package weather

import (
  "encoding/json"
  "errors"
  "fmt"
  "net/http"
  "net/url"
  "strings"
  "strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// GeoLocator holds token and has methods for weather operations
type GeoLocator struct {
  Token string
}

type geoCoordinatesResponse struct {
  Results []struct{
		FormattedAddr string `json:"formatted_address"` 
    Geometry struct{
      Location struct{
        Lat float32 `json:"lat"`
        Lng float32 `json:"lng"`
      } `json:"location"`
    } `json:"geometry"`
  } `json:"results"`
  Status string `json:"status"`
}

type weatherURLData struct {
  Properties struct{
    ForecastURL string `json:"forecastHourly"`
  } `json:"properties"`
  Status int `json:"status"`
  Title string `json:"title"`
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

type location struct {
  Latitude  float32 `json:"latitude"`
  Longitude float32 `json:"longitude"`
	FormattedAddr string
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func (wb* GeoLocator) GetWeather(locationTokens []string) (*discordgo.MessageSend, error) {
	msgsend := discordgo.MessageSend{}

  joinedLoc := strings.Join(locationTokens, " ")
  coords, err := wb.getGeoCoordinates(joinedLoc)
  if err != nil {
    return &msgsend, err
  }

  forecastURL, err := getForecastURL(coords)
  if err != nil {
    return &msgsend, err
  }

  forecastData, err := getForecastData(forecastURL)
  if err != nil {
    return &msgsend, err
	}
	
	forecast := forecastData.Properties.Periods[0]

	imgRes, err := httpClient.Get(forecast.IconURL)
  if err != nil {
    return &msgsend, err
	}

	msgsend.Files = []*discordgo.File{
		&discordgo.File {
			Name: "weather_icon.png",
			ContentType: "image/png",
			Reader: imgRes.Body,
		},
	}

	msg := fmt.Sprintf("**%s**:\nTemp: %d°%s\nWind speed: %s %s\nDescription: %s", coords.FormattedAddr, forecast.Temperature, forecast.TemperatureUnit, forecast.WindSpeed, forecast.WindDirection, forecast.ShortForecast)
	msgsend.Content = msg
  return &msgsend, nil
}

func (wb* GeoLocator) getGeoCoordinates(locStr string) (location, error) {
  values := url.Values{}
  values.Add("address", locStr)
  values.Add("key", wb.Token)

  googleMapsBaseResponse, err := httpClient.Get("https://maps.googleapis.com/maps/api/geocode/json?" + values.Encode())
  if err != nil {
    return location{}, err
  }
  defer googleMapsBaseResponse.Body.Close()

  var gcr geoCoordinatesResponse
  json.NewDecoder(googleMapsBaseResponse.Body).Decode(&gcr)

  if gcr.Status != "OK" {
    return location{}, errors.New("Google maps error: " + gcr.Status)
  }

  loc := location{
    Latitude: gcr.Results[0].Geometry.Location.Lat,
		Longitude: gcr.Results[0].Geometry.Location.Lng,
		FormattedAddr: gcr.Results[0].FormattedAddr,
  }

  return loc, nil
}

func getForecastURL(loc location) (string, error) {
  locStr := fmt.Sprintf("%f,%f", loc.Latitude, loc.Longitude)
  weatherResp, err := httpClient.Get("https://api.weather.gov/points/" + locStr)
  if err != nil {
    return "", err
  }
  defer weatherResp.Body.Close()

  var weatherData weatherURLData
  json.NewDecoder(weatherResp.Body).Decode(&weatherData)

  weatherDataStatus := weatherData.Status
  if weatherDataStatus != 0 {
    return "", errors.New("Forecast URL request responded with status " + strconv.Itoa(weatherDataStatus) + ". Error: " + weatherData.Title)
  }

  forecastURL := weatherData.Properties.ForecastURL
  if forecastURL == "" {
    return forecastURL, errors.New("Forecast URL was not mapped properly")
  }

  return forecastURL, nil
}

func getForecastData(forecastURL string) (forecastData, error) {
  forecastResp, err := httpClient.Get(forecastURL)
  if err != nil {
    return forecastData{}, err
  }
  defer forecastResp.Body.Close()

  var forecast forecastData
  json.NewDecoder(forecastResp.Body).Decode(&forecast)

  forecastDataStatus := forecast.Status
  if forecastDataStatus != 0 {
    return forecastData{}, errors.New("Forecast data request responded with status " + strconv.Itoa(forecastDataStatus))
  }

  return forecast, nil
}