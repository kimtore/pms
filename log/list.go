package log

import (
	"github.com/ambientsound/pms/list"
)

var logLineList list.List

func init() {
	logLineList = list.New()
}

func List() list.List {
	return logLineList
}
