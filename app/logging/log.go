package logging

import (
	"fmt"
	"time"

	"github.com/okira-e/go-as-your-backend/app/utils"
)

type Severity int

const (
	SeverityDebug Severity = iota
	SeverityInfo
	SeverityWarn
	SeverityError
)

func Log(severity Severity, msg string, logContext map[string]any) {
	env := utils.RequireEnv("ENV")

	// suppress debug logs unless explicitly in debug mode
	if severity == SeverityDebug && env != "debug" {
		return
	}

	var tag string
	switch severity {
	case SeverityDebug:
		tag = "[DEBUG]"
	case SeverityInfo:
		tag = "[INFO]"
	case SeverityWarn:
		tag = "[WARN]"
	case SeverityError:
		tag = "[ERROR]"
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	contextStr := ""
	if len(logContext) > 0 {
		contextStr = fmt.Sprintf(" %v", logContext)
	}

	println(now, tag, msg+contextStr)
}
