package LogUtils

import (
	"gotools/Util/FileUtils"
	"io"
	"log"
	"os"
)

var Log *Loggers

type Loggers struct {
	LogContext  map[string]*LoggerContext
	CLoseFiles  []*os.File
	PermitLevel int
}

func init() {
	Log = &Loggers{}
	Log.LogContext = make(map[string]*LoggerContext, 0)
}

func (this *Loggers) RegisterLogContext(name string, Context *LoggerContext) {
	if name != "" {
		this.LogContext[name] = Context
	} else if Context.LogConfig.Prefix != "" {
		this.LogContext[Context.LogConfig.Prefix] = Context
	} else {
		this.LogContext["Unkwown"] = Context
	}

}
func (this *Loggers) getLogContextByName(name string) *LoggerContext {
	if context, ok := this.LogContext[name]; ok {
		return context
	} else {
		return nil
	}
}

func (this *Loggers) CreateLoggerContext(config *LogConfig) (*LoggerContext, error) {
	context := &LoggerContext{}
	context.LogConfig = config
	writers := make([]io.Writer, 0)
	for _, logPath := range config.LogFilePaths {
		f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		writers = append(writers, f)
		this.CLoseFiles = append(this.CLoseFiles, f)
	}
	if config.IsUseStdout {
		writers = append(writers, os.Stdout)
	}
	multiWriter := io.MultiWriter(writers...)
	context.Logger = log.New(multiWriter, config.Prefix, 0)
	return context, nil
}
func (this *Loggers) Clear() {
	for _, f := range this.CLoseFiles {
		f.Close()
	}

}
func (this *Loggers) Save(name string, indent int, contents ...interface{}) {
	LoggerP := this.getLogContextByName(name)
	if LoggerP != nil {
		for _, content := range contents {
			LoggerP.Logger.Println(content)
		}
	} else {
		LoggerP := this.getLogContextByName("Error")
		LoggerP.Logger.Println("LoggerContext is not exist,name: " + name)
	}
}

func (this *Loggers) GetHtml(name, baseHtmlPath string) {
	fileIO := FileUtils.FileIO{}
	LoggerP := this.getLogContextByName(name)
	logHtml := &LogHtml{}
	if LoggerP != nil {
		length := len(LoggerP.LogConfig.LogFilePaths)
		if length < 1 {
			LoggerP = this.getLogContextByName("Error")
			LoggerP.Logger.Println("can't find available logfile")
			return
		}
		lines, err := fileIO.ReadStrLines(LoggerP.LogConfig.LogFilePaths[length-1])
		if err != nil {
			LoggerP = this.getLogContextByName("Error")
			LoggerP.Logger.Println("read failed")
			LoggerP.Logger.Println(err)
			return
		}

		logHtml.Create(baseHtmlPath)
		//logHtml.AddOneLog("test")
		logHtml.AddLogs(lines)
		logHtml.GetFile(LoggerP.LogConfig.LogFilePaths[length-1] + ".html")
		if err != nil {
			LoggerP = this.getLogContextByName("Error")
			LoggerP.Logger.Println("read failed")
			LoggerP.Logger.Println(err)
			return
		}

	} else {
		LoggerP := this.getLogContextByName("Error")
		LoggerP.Logger.Println("loggerContext is not exist,name: " + name)
	}
}
