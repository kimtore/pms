package log

import (
	"github.com/ambientsound/pms/list"
)

var logLineList map[Level]list.List

func List(level Level) list.List {
	return logLineList[level]
}
