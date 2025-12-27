package version

import (
	"fmt"
)

const (
	AppName   = "cliplab"
	HumanName = "cliplab"
)

var (
	GitVersion = "Unknown"
	AppVersion = "0.0.1"
)

func Version() string {
	return fmt.Sprintf("git=%s , app=%s", GitVersion, AppVersion)
}
