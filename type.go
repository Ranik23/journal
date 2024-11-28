package main

import (
	"fmt"
	"github.com/fatih/color"
)


type LogItem struct {
	Service  string
	Time     string
	Message  string
	Priority string
}


func (l LogItem) String() string {
	return fmt.Sprintf(
		"%s Priority=%s %s %s",
		color.GreenString("%s", l.Time),  
		color.RedString("%s", l.Priority),
		color.CyanString("%s", l.Service), 
		l.Message,
	)
}
