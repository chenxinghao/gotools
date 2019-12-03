package LogUtils

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"gotools/Util/FileUtils"
	"io/ioutil"
)

type LogHtml struct {
	Html *goquery.Document
}

func (this *LogHtml) Create(filePath string) error {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	html, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return err
	}
	this.Html = html
	return nil
}

func (this *LogHtml) AddOneLog(content string) {
	LabelStrStart := "<div"
	//if extraParam!=nil{
	//	for _,s:=range extraParam{
	//		LabelStrStart+=" "+s
	//	}
	//}
	LabelStrStart += ">"
	LabelStrEnd := "</div>"
	this.Html.Find("body").AppendHtml(LabelStrStart + content + LabelStrEnd)
}

func (this *LogHtml) AddLogs(contents []string) {
	LabelStrStart := "<div"
	//if extraParam!=nil{
	//	for _,s:=range extraParam{
	//		LabelStrStart+=" "+s
	//	}
	//}
	LabelStrStart += ">"
	LabelStrEnd := "</div>"
	SelectionP := this.Html.Find("body")
	for _, content := range contents {
		SelectionP.AppendHtml(LabelStrStart + content + LabelStrEnd)
	}
}

func (this *LogHtml) GetString() (string, error) {
	res, err := this.Html.Html()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (this *LogHtml) GetFile(filePath string) error {
	fileIO := FileUtils.FileIO{}
	htmlStr, err := this.GetString()
	if err != nil {
		return err
	}
	_, err = fileIO.WriteFile(filePath, []byte(htmlStr), false)
	if err != nil {
		return err
	}
	return nil
}
