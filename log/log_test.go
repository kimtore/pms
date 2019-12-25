package log_test

import (
	"github.com/ambientsound/pms/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClear(t *testing.T) {
	log.Clear()
	assert.Len(t, log.Messages(log.InfoLevel), 0)
	for i := 0; i < 10; i++ {
		log.Infof("foo")
	}
	assert.Len(t, log.Messages(log.InfoLevel), 10)
	log.Clear()
	assert.Len(t, log.Messages(log.InfoLevel), 0)
	assert.Equal(t, 0, log.List().Len())
}

func TestErrorf(t *testing.T) {
	log.Clear()

	fmt := "this %s a string with %d"
	_, err := log.Errorf(fmt, "is", 645)

	assert.NoError(t, err)

	linebuffer := log.Messages(log.ErrorLevel)

	assert.Len(t, linebuffer, 1)
	assert.Len(t, linebuffer[0].Text, 25)
	assert.Equal(t, "this is a string with 645", linebuffer[0].Text)
	assert.Equal(t, log.ErrorLevel, linebuffer[0].Level)
}
