package reminders

import (
	"github.com/araddon/dateparse"
//  "encoding/json"
//  "errors"
//  "fmt"
//  "net/http"
//  "net/url"
//  "strings"
//  "strconv"
	"time"
//
//	"github.com/bwmarrin/discordgo"
)

// GeoLocator holds token and has methods for weather operations
type GeoLocator struct {
  Token string
}

func SetReminder(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (error) {
	if len(args) != 2 {
		return errors.New("Syntax: !jeff remindme \"[date]\" \"[message]\"")
	}

	date := args[0]
	msg := args[1]

	t, err := dateparse.ParseLocal(date)

	msgsend := discordgo.MessageSend{}
}

func (wb* GeoLocator) GetWeather(s *discordgo.Session, m *discordgo.MessageCreate, locationTokens []string) (error) {
	if len(locationTokens) == 0 {
		return errors.New("Provide a location, yo")
	}
	msgsend := discordgo.MessageSend{}

  joinedLoc := strings.Join(locationTokens, " ")
  coords, err := wb.getGeoCoordinates(joinedLoc)
  if err != nil {
    return err
  }

  forecastURL, err := getForecastURL(coords)
  if err != nil {
    return err
  }

  forecastData, err := getForecastData(forecastURL)
  if err != nil {
    return err
	}
	
	forecast := forecastData.Properties.Periods[0]

	imgRes, err := httpClient.Get(forecast.IconURL)
  if err != nil {
    return err
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
	s.ChannelMessageSendComplex(m.ChannelID, &msgsend)
  return nil
}

