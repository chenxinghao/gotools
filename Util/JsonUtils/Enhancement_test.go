package JsonUtils

import (
	"encoding/json"
	"fmt"
	"testing"
)

type JsonInterfaceFormatTest struct {
	Name string
	Age  int64
}

func TestJsonInterfaceFormat(t *testing.T) {
	enhancement := &Enhancement{}
	str := ""
	enhancement.JsonInterfaceFormat("test", &str)

	var number int64
	number = 0
	enhancement.JsonInterfaceFormat(int64(10), &number)

	var fnumber float64
	fnumber = 0
	enhancement.JsonInterfaceFormat(10.9, &fnumber)

	jsonInterfaceFormatTest := &JsonInterfaceFormatTest{Name: "cxh", Age: 120}
	jsonInterfaceFormatTest1 := &JsonInterfaceFormatTest{}
	temp, _ := json.Marshal(jsonInterfaceFormatTest)
	tempmap := make(map[string]interface{})
	json.Unmarshal(temp, &tempmap)
	enhancement.JsonInterfaceFormat(tempmap, jsonInterfaceFormatTest1)
}

func TestSliceHandler(t *testing.T) {
	enhancement := &Enhancement{}
	b := []interface{}{float64(1), float64(2)}
	enhancement.SliceHandler(&b)
	fmt.Println(b[0])

	b = []interface{}{map[string]interface{}{"name": "cxh"}, map[string]interface{}{"age": 26}}
	enhancement.SliceHandler(&b)
	fmt.Println(b[0])

	b = []interface{}{[]interface{}{1, 2, 3}, []interface{}{3, 4, 5}}
	enhancement.SliceHandler(&b)
	fmt.Println(b[0])

	q := enhancement.queue
	for e := q.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func TestMapHandler(t *testing.T) {
	enhancement := &Enhancement{}
	testMap := make(map[string]interface{})
	testMap["float"] = float64(1)
	testMap["ptr1"] = &map[string]interface{}{"cxh": "123"}
	testMap["ptr2"] = &[]interface{}{1, 2, 3, 4}
	testMap["map"] = map[string]interface{}{"cxh": 345}
	testMap["slice"] = []interface{}{1, 2, 3, 4}
	enhancement.MapHandler(&testMap)
	fmt.Println(testMap)

	q := enhancement.queue
	for e := q.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func TestHandler(t *testing.T) {
	enhancement := &Enhancement{}
	jsonStr := `{"host": "http://localhost:9090","port": 9090,"analytics_file": "","static_file_version": 1,"static_dir": "E:/Project/goTest/src/","templates_dir": "E:/Project/goTest/src/templates/","serTcpSocketHost": ":12340","serTcpSocketPort": 12340,"fruits": ["apple", "peach"],"objects":[{"cxh":"name"},{"age":99}],"lists":[[1,2,3],[4,5,6]],"lists2":[[1,2,3],7],"testObject":{},"testObject2":{"zzz":{"1":"2","3":"4"}}}`

	mP := enhancement.Json2TypeMap(jsonStr)

	enhancement.TypeMap2StructStr("start", mP)

}
