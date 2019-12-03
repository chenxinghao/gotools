package LogUtils

import (
	"fmt"
	"testing"
)

func TestLoggers_Save(t *testing.T) {
	LConfig := &LogConfig{}
	filepath := make([]string, 0)
	LConfig.IsUseStdout = true
	LConfig.Prefix = "[test]"
	LConfig.LogFilePaths = append(filepath, "C:\\workspace\\GoProject\\gotools\\src\\gotools\\Util\\test.log")
	LConfig.AutoIndentStr = ">>>>"
	LConfig.Level = 10
	LCP, err := Log.CreateLoggerContext(LConfig)
	if err != nil {
		fmt.Println(err)
	}
	Log.RegisterLogContext("test", LCP)

	Log.Save("test", 0, "testing", LConfig)

	Log.GetHtml("test", "C:\\workspace\\GoProject\\gotools\\src\\gotools\\Util\\LogUtils\\View\\base.html")

	Log.Clear()
}
