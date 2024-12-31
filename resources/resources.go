package resources

import (
	"time"
)

type TerraformConfiguration struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	Timeout      time.Duration
}
