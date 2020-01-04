package log

import (
	"fmt"
	"github.com/ambientsound/pms/list"
	"io"
	"time"
)

// Message is a message passed from anywhere inside PMS, relayed to the user
// through the statusbar.
type Message struct {
	Timestamp time.Time
	Level     Level
	Text      string
}

// Slice of all messages logged since starting up.
// The level referred to contains all messages with AT LEAST that severy.
// i.e., messages[INFO] will also contain messages with ERROR level.
var messages map[Level][]Message

func (msg Message) Write(w io.Writer) (int, error) {
	prefix := fmt.Sprintf("[%010.3f] [%s] ", msg.Timestamp.Sub(since).Seconds(), strLevel[msg.Level])
	return w.Write([]byte(prefix + msg.Text + "\n"))
}

func Messages(level Level) []Message {
	return messages[level]
}

func Clear() {
	messages = make(map[Level][]Message)
	logLineList = make(map[Level]list.List)

	for level := range strLevel {
		messages[level] = make([]Message, 0)
		logLineList[level] = list.New()
		logLineList[level].SetID(level.String() + " logs")
		logLineList[level].SetName("Log console")
		logLineList[level].SetVisibleColumns([]string{"timestamp", "logLevel", "logMessage"})
	}
}

func Last(level Level) *Message {
	n := len(messages[level])
	if n == 0 {
		return nil
	}
	return &messages[level][n-1]
}

func init() {
	Clear()
}

func appendMessage(msg Message) {
	for level := range strLevel {
		if msg.Level > level {
			continue
		}
		messages[level] = append(messages[level], msg)
		logLineList[level].Add(list.Row{
			list.RowIDKey: string(logLineList[level].Len()),
			"logLevel":    msg.Level.String(),
			"logMessage":  msg.Text,
			"timestamp":   msg.Timestamp.Format(time.RFC822),
		})
	}
}
