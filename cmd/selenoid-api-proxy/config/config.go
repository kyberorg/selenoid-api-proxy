package config

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"net/url"
)

var (
	port = kingpin.Flag("port", "Port server start at").Envar("PORT").
		Default("4443").Uint16()

	selenoidApiUrl = kingpin.Flag("selenoid.api.url", "Selenoid URL").Envar("SELENOID_API_URL").
			Default("http://localhost:4444").URL()

	token = kingpin.Flag("token", "Token for authorization").Envar("TOKEN").Required().String()
)

//internal vars
var (
	appConfig *applicationConfiguration
)

type applicationConfiguration struct {
	Port           uint16
	SelenoidApiUrl **url.URL
	Token          string
}

func init() {
	kingpin.Parse()

	appConfig = &applicationConfiguration{
		Port:           *port,
		SelenoidApiUrl: selenoidApiUrl,
		Token:          *token,
	}
}

//GetAppConfig returns application configuration object
func GetAppConfig() *applicationConfiguration {
	return appConfig
}
