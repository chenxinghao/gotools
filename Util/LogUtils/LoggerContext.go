package LogUtils

import "log"

type LoggerContext struct {
	LogConfig *LogConfig
	Logger    *log.Logger
}
