package log_test

import (
	"github.com/ambientsound/pms/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClear(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Printf("foo")
	}
	assert.Len(t, log.Lines(), 10)
	log.Clear()
	assert.Len(t, log.Lines(), 0)
}

func TestPrintf(t *testing.T) {
	log.Clear()

	fmt := "this %s a string with %d"
	n, err := log.Printf(fmt, "is", 645)

	assert.NoError(t, err)

	linebuffer := log.Lines()

	assert.Len(t, linebuffer, 1)
	assert.Len(t, linebuffer[0], n)
	assert.Equal(t, "this is a string with 645", linebuffer[0])
}

func TestErrorf(t *testing.T) {
	log.Clear()

	fmt := "this %s a string with %d"
	n, err := log.Errorf(fmt, "is", 645)

	assert.NoError(t, err)

	linebuffer := log.Lines()

	assert.Len(t, linebuffer, 1)
	assert.Len(t, linebuffer[0], n)
	assert.Equal(t, "[000000.000] [ERROR] this is a string with 645", linebuffer[0])
}
