package hubspot

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Config is the global configuration object that holds global configuration settings
var Config *ConfigStruct

//ConfigStruct holds the various configuration options
type ConfigStruct struct {
	Environment   string
	RootURL       string
	HubspotAPIKey string
	Logging       bool
	// Certain API endpoints require an OAuth application; those are specified here
	HubspotApplicationID     string
	HubspotUserID            string
	HubspotClientID          string
	HubspotClientSecret      string
	HubSpotOAuthRefreshToken string
}

// ConfigSetup sets up the config struct with data from the environment
func ConfigSetup() *ConfigStruct {
	c := new(ConfigStruct)

	c.Environment = strings.ToLower(os.Getenv("HUBSPOT_SDK_ENV"))
	if c.Environment == "prod" || c.Environment == "production" {
		c.Environment = "production"
	} else if c.Environment == "" || c.Environment == "dev" || c.Environment == "development" {
		c.Environment = "dev"
	}

	c.RootURL = strings.ToLower(os.Getenv("HUBSPOT_SDK_ROOT_URL"))
	if c.RootURL == "" {
		c.RootURL = "https://api.hubapi.com"
	}
	if !strings.HasSuffix(c.RootURL, "/") {
		c.RootURL += "/"
	}

	c.HubspotAPIKey = os.Getenv("HUBSPOT_SDK_API_KEY")
	if c.HubspotAPIKey == "" {
		c.HubspotAPIKey = "demo"
	}

	c.HubspotApplicationID = os.Getenv("HUBSPOT_SDK_APPLICATION_ID")
	if c.HubspotApplicationID == "" {
		c.HubspotApplicationID = "test"
	}

	c.HubspotUserID = os.Getenv("HUBSPOT_SDK_USER_ID")
	if c.HubspotUserID == "" {
		c.HubspotUserID = "test"
	}
	c.HubspotClientID = os.Getenv("HUBSPOT_SDK_CLIENT_ID")
	c.HubspotClientSecret = os.Getenv("HUBSPOT_SDK_CLIENT_SECRET")
	c.HubSpotOAuthRefreshToken = os.Getenv("HUBSPOT_SDK_OAUTH_REFRESH_TOKEN")

	shouldLog := strings.ToLower(os.Getenv("HUBSPOT_SDK_LOGGING"))
	if shouldLog == "off" || shouldLog == "false" || shouldLog == "no" {
		c.Logging = false
	} else {
		c.Logging = true
	}

	Config = c

	return c
}

// init is called when the host application starts up and sets the
// Configuration and logging settings
func init() {
	ConfigSetup()
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

// log provides structured logging through logrus. We support info, warning, and error
func log(level, key, message string, data interface{}) string {
	if Config.Logging {
		level = strings.ToLower(level)

		fields := logrus.Fields{
			"key":  key,
			"data": data,
		}

		switch level {
		case "info":
			logrus.WithFields(fields).Info(message)
		case "warning":
			logrus.WithFields(fields).Warning(message)
		case "error":
			logrus.WithFields(fields).Error(message)
		}
		return fmt.Sprintf("%s: %s", strings.ToUpper(level), key)
	}
	return ""
}
