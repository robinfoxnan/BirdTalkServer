package utils

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"reflect"
	"strconv"
	"strings"
)

// 辅助函数：检查切片中是否包含某个元素
func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func convertToMap(params interface{}) (map[string]interface{}, bool) {
	// 使用类型断言将 params 转换为 map[string]interface{}
	if paramMap, ok := params.(map[string]interface{}); ok {
		// 如果转换成功，返回 map 及 true
		return paramMap, true
	} else {
		// 如果转换失败，返回 nil 及 false
		return nil, false
	}
}

// 在AnyToMap内部使用，当字段为map,加上前缀，直接保存到父级别
func addMapToMap(params interface{}, data map[string]interface{}, prefix string) error {
	paramMap, ok := params.(map[string]string)
	if !ok {
		//fmt.Println("convert error")
		return errors.New("")
	}
	for k, v := range paramMap {
		name := prefix + "." + k
		data[name] = v
	}
	return nil
}

// 对于字符串数字，使用"|"拼接到一起
//func convertArrayToStr(lst []string) string {
//	return strings.Join(lst, "|")
//}
//
//func convertStrToArray(str string) []string {
//	return strings.Split(str, "|")
//}

// 将切片转换为使用 "|" 符号拼接的字符串
//func convertSliceToString(slice interface{}) string {
//	v := reflect.ValueOf(slice)
//	var strSlice = make([]string, v.Len())
//	for i := 0; i < v.Len(); i++ {
//
//		if str, ok := v.Index(i).Interface().(string); ok {
//			// 这里的 str 就是切片中的元素，且类型是字符串
//			strSlice[i] = str
//		} else {
//			// 切片中的元素不是字符串类型，处理错误或者做其他操作
//		}
//	}
//	return strings.Join(strSlice, "|")
//}

func convertSliceToJsonString(slice interface{}) string {
	// 将切片序列化为 JSON 字符串
	bytes, err := json.Marshal(slice)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// 使用 JSON 反序列化将字符串转换为切片
func convertStringToSlice(str string) ([]string, error) {
	// 将 JSON 字符串反序列化为切片
	var slice []string
	err := json.Unmarshal([]byte(str), &slice)
	if err != nil {
		return nil, err
	}
	return slice, nil
}

func printFieldInfo(field reflect.StructField) {
	fName := field.Name
	fType := field.Type.Kind()
	//fValue := v.FieldByName(t.Field(i).Name)
	//fTypeK = t.Field(i).Type.Kind()
	//此处可得到所有字段的基本信息
	fmt.Printf("%v		%v	  \n", fName, fType)
}

// 从指针或者结构体转为map类型
func AnyToMap(item interface{}, excludedFields []string) (map[string]interface{}, error) {
	t := reflect.TypeOf(item)
	v := reflect.ValueOf(item)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		//fmt.Println("不是struct类型")
		return nil, errors.New("item must be struct or pointer")
	}
	//fmt.Printf("name:'%v' kind:'%v'\n", t.Name(), t.Kind())

	var data = make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		// 检查字段的可导出性
		field := t.Field(i)
		// 私有成员变量，跳过
		if field.PkgPath != "" {
			continue
		}

		// 检查字段是否在排除列表中
		if contains(excludedFields, field.Name) {
			continue
		}

		fieldValue := v.Field(i).Interface()
		if field.Type.Kind() == reflect.Map {
			addMapToMap(fieldValue, data, t.Field(i).Name)
		} else if field.Type.Kind() == reflect.Array || field.Type.Kind() == reflect.Slice {

			key := t.Field(i).Name
			data[key] = convertSliceToJsonString(fieldValue)

		} else {
			//key := strings.ToLower(t.Field(i).Name)
			key := t.Field(i).Name
			data[key] = fieldValue
		}

	}
	return data, nil
}

// 提取data中所有前缀为prefix的key； key类似于"Params.tilte", Params就是前缀, key去掉前缀和点作为newkey，
// 把newkey 和value 放到新的map中返回
func abstractMapByPrefix(prefix string, data map[string]string) map[string]string {
	//fmt.Println("prefix=", prefix)
	abstractMap := make(map[string]string)
	// 遍历 data 中的键值对
	for key, value := range data {
		// 检查键是否以 prefix 开头
		if strings.HasPrefix(key, prefix) {
			// 提取新的键名
			newKey := strings.TrimPrefix(key, prefix+".")
			// 将新键名和值添加到 abstractMap 中
			abstractMap[newKey] = value
		}
	}

	return abstractMap
}

// 从map转为希望的类型
func FromMapString(data map[string]string, item interface{}) error {
	v := reflect.ValueOf(item).Elem()

	var fName string
	var keyName string
	var fType reflect.Kind
	var id int
	var err error
	var id64 int64
	var id32 int32

	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		fName = fieldInfo.Name
		fType = fieldInfo.Type.Kind()
		keyName = fName //strings.ToLower(fName)
		if strValue, ok := data[keyName]; ok {
			if fType == reflect.String {
				v.FieldByName(fieldInfo.Name).SetString(strValue)
			} else if fType == reflect.Int64 {
				if id, err = strconv.Atoi(strValue); err == nil {
					id64 = int64(id)
					v.FieldByName(fName).SetInt(id64)
				} else {
					fmt.Println(id64, err)
				}
			} else if fType == reflect.Int32 {
				if id, err = strconv.Atoi(strValue); err == nil {
					id32 = int32(id)
					v.FieldByName(fName).SetInt(int64(id32))
				} else {
					fmt.Println(err)
				}
			} else if fType == reflect.Slice {

				// 解析切片类型
				sliceType := fieldInfo.Type
				//fmt.Println(sliceType)
				tempSlice, err1 := convertStringToSlice(strValue)
				if err1 != nil {
					continue
				}
				sliceValue := reflect.MakeSlice(sliceType, len(tempSlice), len(tempSlice))
				// 将值逐个赋给切片
				for sliceIndex, sliceItem := range tempSlice {
					sliceValue.Index(sliceIndex).SetString(sliceItem)
					//sliceValue.Index(1).SetString("聊天")
				}

				//fmt.Println(strValue, sliceValue)
				v.FieldByName(fName).Set(sliceValue)
			}
		} else {
			if fType == reflect.Map {

				//fmt.Println("find map \n")
				// 获取字段的实际类型
				fieldType := fieldInfo.Type
				// 创建一个空的 map，类型为字段的类型
				subMap := reflect.MakeMap(fieldType)
				// 获取字段的键类型和值类型
				keyType := fieldType.Key()
				valueType := fieldType.Elem()
				// 获取 map 的反射值
				mapValue := reflect.ValueOf(subMap.Interface())
				// 获取数据中以字段名为前缀的键值对
				namePrefix := fieldInfo.Name
				subData := abstractMapByPrefix(namePrefix, data)
				// 遍历数据中的键值对，将其添加到 subMap 中
				for k, v1 := range subData {
					// 将键和值转换为对应的类型
					keyValue := reflect.New(keyType).Elem()
					keyValue.SetString(k)
					valueValue := reflect.New(valueType).Elem()
					valueValue.SetString(v1)
					// 将键值对添加到 subMap 中
					mapValue.SetMapIndex(keyValue, valueValue)
				}
				// 将填充好的 map 设置到结构体字段中
				v.FieldByName(fName).Set(mapValue)
			}
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
