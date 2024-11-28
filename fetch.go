package main

import (
	"fmt"
	"time"
	"github.com/coreos/go-systemd/sdjournal"
)

func ConvertUnixTimestampToTime(timestamp uint64) string {
	sec := timestamp / 1_000_000
	usec := timestamp % 1_000_000
	t := time.Unix(int64(sec), int64(usec)*1000)
	return t.Format("2006-01-02 15:04:05")
}

func FetchLogs(service string, priority string, limit int) (*[]LogItem, error) {
	journal, err := sdjournal.NewJournal()
	if err != nil {
		return nil, fmt.Errorf("failed to open journal: %w", err)
	}
	defer journal.Close()

	if service != "" {
		match := fmt.Sprintf("_SYSTEMD_UNIT=%s", service)
		if err := journal.AddMatch(match); err != nil {
			return nil, fmt.Errorf("failed to add service match: %w", err)
		}
	}

	if priority != "" {
		p := fmt.Sprintf("PRIORITY=%s", priority)
		if err := journal.AddMatch(p); err != nil {
			return nil, fmt.Errorf("failed to add priority match: %w", err)
		}
	}

	if err := journal.SeekTail(); err != nil {
		return nil, ErrSeekTail
	}

	count := 0

	for i := 0; i < limit; i++ {
		count += 1
		_, err := journal.Previous()
		if err != nil {
			count -= 1
			break
		}
	}

	limit = count

	i, err := journal.Next()
	if err != nil || i == 0 {
		return nil, ErrNoLogsFound
	}
	
	var Logs []LogItem
	for i := 0; i < limit; i++ {

		entry, err := journal.GetEntry()

		if err != nil {
			if err.Error() == "EOF" {
				return nil, ErrEOF
			}
			return nil, err
		}

		serviceField, ok := entry.Fields["_SYSTEMD_UNIT"]
		if !ok {
			serviceField = "unknown"
		}
		message, ok := entry.Fields["MESSAGE"]
		if !ok {
			message = "no message"
		}
		priority, ok := entry.Fields["PRIORITY"]
		if !ok {
			priority = "unknown"
		}

		Logs = append(Logs, LogItem{
			Service:  serviceField,
			Time:     ConvertUnixTimestampToTime(entry.RealtimeTimestamp),
			Message:  message,
			Priority: priority,
		})

		if _, err := journal.Next(); err != nil {
			break
		}
	}

	if len(Logs) == 0 {
		return nil, ErrNoLogsFound
	}

	return &Logs, nil
}
