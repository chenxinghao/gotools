package JsonUtils

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

const UNKNOW = "unknow"

type Enhancement struct {
	queue list.List
}

//根据指针转换成对应对象
func (this *Enhancement) JsonInterfaceFormat(data interface{}, objPointer interface{}) (bool, error) {
	var kind string
	TypeOfObj := reflect.TypeOf(objPointer)
	TypeOfData := reflect.TypeOf(data)
	if TypeOfObj.Kind().String() != "ptr" {
		return false, errors.New("must pass in pointer")
	}
	if TypeOfData.Kind().String() == "ptr" {
		return false, errors.New("must pass in struct")
	}
	kind = reflect.ValueOf(objPointer).Elem().Kind().String()
	switch kind {
	case "struct":
		if content1, ok1 := data.(map[string]interface{}); ok1 {
			jstr, err := json.Marshal(content1)
			if err != nil {
				return false, err
			}
			err = json.Unmarshal(jstr, objPointer)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	case "int64":
		if content1, ok1 := data.(int64); ok1 {
			Value := reflect.ValueOf(objPointer)
			Value = reflect.Indirect(Value)
			Value.SetInt(content1)
		}
		return true, nil
	case "float64":
		if content1, ok1 := data.(float64); ok1 {
			Value := reflect.ValueOf(objPointer)
			Value = reflect.Indirect(Value)
			Value.SetFloat(content1)
		}
		return true, nil
	case "string":
		if content1, ok1 := data.(string); ok1 {
			Value := reflect.ValueOf(objPointer)
			Value = reflect.Indirect(Value)
			Value.SetString(content1)
		}
		return true, nil
	default:
		return false, errors.New("Can't Find Supported Type")
	}
	return false, nil

}

func (this *Enhancement) SliceHandler(listP *[]interface{}) {
	list := *listP
	tempList := []interface{}{}
	if getSameType(list) {
		firstType := reflect.TypeOf(list[0]).Kind().String()
		switch firstType {
		case "bool", "string", "float64":
			*listP = (*listP)[0:0] //清空list
			*listP = append(*listP, firstType)
		case "map":
			newMap := make(map[string]interface{})
			for _, i := range list {
				m := i.(map[string]interface{})
				for k, v := range m {
					if _, ok := newMap[k]; !ok {
						newMap[k] = v
					}
				}
			}
			this.queue.PushBack(&newMap)
			*listP = (*listP)[0:0]
			*listP = append(*listP, &newMap)
		case "slice":
			for _, i := range list {
				l := i.([]interface{})
				tempList = append(tempList, l...)
			}
			this.queue.PushBack(&tempList)
			*listP = (*listP)[0:0]
			*listP = append(*listP, &tempList)
		default:
			fmt.Println("err")
		}
	} else {
		*listP = (*listP)[0:0]
		*listP = append(*listP, UNKNOW)
	}

}
func (this *Enhancement) MapHandler(objectMap *map[string]interface{}) {
	for k, v := range *objectMap {
		typeStr := reflect.TypeOf(v).Kind().String()
		switch typeStr {
		case "bool", "string", "float64":
			(*objectMap)[k] = typeStr
		case "ptr":
		case "map":
			tmpV := v.(map[string]interface{})
			(*objectMap)[k] = &tmpV
			this.queue.PushBack(&tmpV)
		case "slice":
			tmpV := v.([]interface{})
			(*objectMap)[k] = &tmpV
			this.queue.PushBack(&tmpV)
		default:
			panic(errors.New("find unkonw object!!!!!"))
		}
	}
}

func getSameType(list []interface{}) bool {
	if len(list) < 1 {
		return false
	}
	firstType := reflect.TypeOf(list[0]).Kind()
	for _, sublist := range list {
		if firstType != reflect.TypeOf(sublist).Kind() {
			return false
		}
	}
	return true
}

func (this *Enhancement) GetStructStr(name string, param interface{}) (string, string) {
	value := reflect.ValueOf(param)
	typeStr := value.Kind().String()
	if typeStr != "ptr" {
		if data, ok := param.(string); ok {
			if data == UNKNOW {
				return "interface{}", "Param"
			}
			return data, "Param"
		}
	}
	if value.Elem().Kind().String() == "map" {
		if data, ok := param.(*map[string]interface{}); ok {
			if len(*data) <= 0 {
				return "interface{}", "Obj"
			}
		}
		return name, "Obj"
	}
	if value.Elem().Kind().String() == "slice" {
		if data, ok := param.(*[]interface{}); ok {
			str, _ := this.GetStructStr(name, (*data)[0])
			return "[]" + str, "List"
		}
	}
	panic(errors.New("param is not type of *[]interface{},ptr,*map[string]interface{}"))
}

func (this *Enhancement) FormatStructStr(name, typeStr, str string) string {
	return name + typeStr + "    " + str + "    `json:\"" + name + "\"`"
}

func (this *Enhancement) Json2TypeMap(jsonStr string) *map[string]interface{} {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &m)
	this.queue.PushBack(&m)
	for {
		if this.queue.Len() <= 0 {
			break
		}
		element := this.queue.Front()
		this.queue.Remove(element)
		typeStr := reflect.ValueOf(element.Value).Elem().Kind().String()
		if typeStr == "map" {
			tmpP := element.Value.(*map[string]interface{})
			this.MapHandler(tmpP)
		} else if typeStr == "slice" {
			tmpP := element.Value.(*[]interface{})
			this.SliceHandler(tmpP)
		}
	}

	bt, _ := json.Marshal(m)
	fmt.Println(string(bt))

	return &m
}

func (this *Enhancement) TypeMap2OneStructStr(name string, mP *map[string]interface{}) string {
	var pairParam [2]interface{}

	fmt.Println("type " + name + " struct{")
	for k, v := range *mP {

		str, typeStr := this.GetStructStr(k, v)
		if typeStr == "Obj" {
			pairParam = [2]interface{}{k, v}
			this.queue.PushBack(pairParam)
		}
		str = this.FormatStructStr(k, typeStr, str)
		fmt.Println("    " + str)

	}
	fmt.Println("}")

	return ""
}

func (this *Enhancement) TypeMap2StructStr(startName string, mP *map[string]interface{}) {
	pairParam := [2]interface{}{startName, mP}
	this.queue.PushBack(pairParam)
	for {
		if this.queue.Len() <= 0 {
			break
		}
		element := this.queue.Front()
		this.queue.Remove(element)
		tempParam := element.Value.([2]interface{})
		this.TypeMap2OneStructStr(tempParam[0].(string), tempParam[1].(*map[string]interface{}))
	}

}
