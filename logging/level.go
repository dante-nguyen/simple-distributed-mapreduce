package logging

import "fmt"

type Level int
type levelInfo struct {
	name     string
	severity int
}

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var infoByLevel = map[Level]levelInfo{
	DEBUG:   {name: "DEBUG", severity: 0},
	INFO:    {name: "INFO", severity: 1},
	WARNING: {name: "WARNING", severity: 2},
	ERROR:   {name: "ERROR", severity: 3},
	FATAL:   {name: "FATAL", severity: 4},
}

func (l Level) Name() string {
	return mustGetLevelInfo(l).name
}

func (l Level) Allow(other Level) bool {
	thisSev := mustGetLevelInfo(l).severity
	otherSev := mustGetLevelInfo(other).severity
	return thisSev <= otherSev
}

func assertValidLevel(l Level) {
	mustGetLevelInfo(l)
}

func mustGetLevelInfo(l Level) levelInfo {
	if i, err := getLevelInfo(l); err != nil {
		panic(err)
	} else {
		return i
	}
}

func getLevelInfo(l Level) (levelInfo, error) {
	if i, ok := infoByLevel[l]; ok {
		return i, nil
	} else {
		return i, invalidLevel(l)
	}
}

func invalidLevel(l Level) error {
	return fmt.Errorf("invalid log level %v", l)
}
