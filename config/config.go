package config

import (
	"errors"
	"nav_service/hilog"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//var flavor = "dev"

var config Config

func GetConfig() Config {
	return config
}

func init() {
	hilog.Info("环境是：" + flavor)
	_ = Unmarshal()
}

func Unmarshal() error {
	data, err := os.ReadFile("config/" + flavor + ".ini")
	if err != nil {
		return err
	}
	config = Config{}
	return unmarshal(&config, data)
}

func unmarshal(result interface{}, data []byte) error {
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Ptr {
		return errors.New("result type kind must be ptr")
	}
	if resultType.Elem().Kind() != reflect.Struct {
		return errors.New("result must be struct")
	}
	lines := strings.Split(string(data), "\n")
	return unmarshalLine(result, lines, 0)
}

func unmarshalLine(result interface{}, lines []string, index int) (err error) {
	if !validLines(lines, index) {
		return
	}
	line := strings.TrimSpace(lines[index])
	if len(line) == 0 || isComments(line) {
		return unmarshalLine(result, lines, index+1)
	}
	if isBeginStruct(line) {
		iniName := line[1 : len(line)-1]
		elem := reflect.TypeOf(result).Elem()
		for i := 0; i < elem.NumField(); i++ {
			field := elem.Field(i)
			tagValue := field.Tag.Get("ini")
			if iniName == tagValue || iniName == field.Name {
				if err = unmarshalStruct(result, valueFieldsByTag(reflect.ValueOf(result).Elem().FieldByName(field.Name)), lines, index+1); err != nil {
					return err
				}
				break
			}
		}
	} else {
		if err = unmarshalStruct(result, valueFieldsByTag(reflect.ValueOf(result).Elem()), lines, index); err != nil {
			return err
		}
	}
	return
}

func unmarshalStruct(result interface{}, fieldsMap map[string]reflect.Value, lines []string, index int) (err error) {
	if !validLines(lines, index) {
		return
	}
	line := strings.TrimSpace(lines[index])
	if len(line) == 0 || isComments(line) {
		return unmarshalStruct(result, fieldsMap, lines, index+1)
	}
	if isBeginStruct(line) {
		return unmarshalLine(result, lines, index)
	}
	key := strings.TrimSpace(line[0:strings.Index(line, "=")])
	val := strings.TrimSpace(line[strings.Index(line, "=")+1:])
	field := fieldsMap[key]
	switch field.Type().Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Int:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint:
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Float32:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	default:
		hilog.Error("unsupported")
	}
	return unmarshalStruct(result, fieldsMap, lines, index+1)
}

func valueFieldsByTag(value reflect.Value) map[string]reflect.Value {
	fieldsMap := make(map[string]reflect.Value)
	valueType := value.Type()
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		tagName := field.Tag.Get("ini")
		fieldsMap[tagName] = value.FieldByName(field.Name)
	}
	return fieldsMap
}

func validLines(lines []string, index int) bool {
	return len(lines) >= 0 && len(lines) > index
}

func isBeginStruct(str string) bool {
	str = strings.TrimSpace(str)
	return str[0] == '[' && str[len(str)-1] == ']'
}

func isComments(str string) bool {
	str = strings.TrimSpace(str)
	return str[0] == ';' || str[0] == '#'
}
