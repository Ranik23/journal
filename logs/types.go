package log

// import "fmt"

type LogItem struct {
	Service  string
	Time     string
	Message  string
	Priority string
}


// func (l *LogItem) Printer() string {
// 	priorityColor := map[string]string{
// 		"INFO":    "\033[1;34m",
// 		"WARNING": "\033[1;33m", 
// 		"ERROR":   "\033[1;31m", 
// 		"DEBUG":   "\033[1;36m",
// 	}

// 	color, exists := priorityColor[l.Priority]
// 	if !exists {
// 		color = "\033[0m" 
// 	}

// 	return fmt.Sprintf(
// 		"%s[%s]\033[0m \033[1;32m%s\033[0m %s: %s",
// 		color,      
// 		l.Priority, 
// 		l.Time,     
// 		l.Service, 
// 		l.Message,  
// 	)
// }

