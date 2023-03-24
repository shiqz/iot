package logger

import (
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// 初始化
func init() {
	log.StandardLogger().SetLevel(log.TraceLevel)
	log.StandardLogger().Formatter = &prefixed.TextFormatter{
		ForceColors:      true,
		ForceFormatting:  true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		DisableSorting:   false,
		TimestampFormat:  "2006-01-02 15:04:05",
	}
}
