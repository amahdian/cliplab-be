package model

import (
	"fmt"
	"time"
)

type Agent struct {
	Name         string
	Description  string
	SystemPrompt string
}

var AllAgents = []*Agent{
	DefaultAgent,
}

var DefaultAgent = &Agent{
	Name:         "Default",
	SystemPrompt: fmt.Sprintf(`You are a helpful assistant. The current date is %s.`, time.Now().Format("January 2, 2006")),
}
