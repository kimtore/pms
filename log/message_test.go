package log_test

import (
	"github.com/ambientsound/pms/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test that logs with similar og higher error levels are returned when requesting a list of log lines.
func TestLevels(t *testing.T) {
	log.Clear()

	// Test that requesting a high log level (more verbose) returns also lower levels.
	for i := 0; i < 10; i++ {
		log.Debugf("foo")
		log.Infof("bar")
		log.Errorf("baz")
	}
	assert.Len(t, log.Messages(log.DebugLevel), 30)
	assert.Len(t, log.Messages(log.InfoLevel), 20)
	assert.Len(t, log.Messages(log.ErrorLevel), 10)

	assert.Equal(t, 30, log.List(log.DebugLevel).Len())
	assert.Equal(t, 20, log.List(log.InfoLevel).Len())
	assert.Equal(t, 10, log.List(log.ErrorLevel).Len())

	// Test that ten of each log level were logged.
	lens := make(map[log.Level]int)
	allMessages := log.Messages(log.DebugLevel)
	for _, msg := range allMessages {
		lens[msg.Level]++
	}
	assert.Equal(t, 10, lens[log.DebugLevel])
	assert.Equal(t, 10, lens[log.InfoLevel])
	assert.Equal(t, 10, lens[log.ErrorLevel])
}

// Test that Last() returns the last logged message.
func TestLast(t *testing.T) {
	log.Clear()
	last := log.Last(log.DebugLevel)
	assert.Nil(t, last)

	log.Errorf("error")
	last = log.Last(log.ErrorLevel)

	assert.NotNil(t, last)
	assert.Equal(t, "error", last.Text)

	log.Infof("info")
	last = log.Last(log.ErrorLevel)
	assert.NotNil(t, last)
	assert.Equal(t, "error", last.Text)
}
