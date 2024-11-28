package main

import "fmt"


var (
	ErrNoLogsFound = fmt.Errorf("no logs found matching criteria")
	ErrEOF = fmt.Errorf("end of file")
	ErrSeekTail = fmt.Errorf("failed to seek tail")
	ErrRandom = fmt.Errorf("random")
)
