package config

//
//import (
//	"fmt"
//	"io/ioutil"
//	"reflect"
//	"strconv"
//	"strings"
//)
//
////MysqlConfig  mysql配置文件结构体
//type MysqlConfig struct {
//	Address  string `ini:"address"`
//	Port     int64  `ini:"port"`
//	Username string `ini:"username"`
//	Password string `ini:"password"`
//}
//
////RedisConfig  redis配置文件结构体
//type RedisConfig struct {
//	Host     string `ini:"host"`
//	Port     int32  `ini:"port"`
//	Password string `ini:"password"`
//	Database string `ini:"database"`
//	Test     bool   `ini:"test"`
//}
//
////IniConfig 所有配置文件的结构体集合
//type IniConfig struct {
//	MysqlConfig `ini:"mysql"`
//	RedisConfig `ini:"redis"`
//}
//
////将ini配置文件中的每个小节给写到对应的结构体
////1.验证参数
////1.1是否是指针类型（因为需要在函数中对其赋值）
////1.2是否是结构体指针
////2.读取配置文件
////2.1读取配置文件得到字节类型数据并赋值给一个变量，
////2.2转换为字符串，并以空格来分割字符串
////2.3跳过注释（注释是以;和#开头的）以及空行
////2.4以[来找到小节的头；以此进入每个小节中循环处理
////2.5如果不是以[开头的话，就进入到小节中的内容中拿到键值对
//func loadIni(fileName string, v interface{}) (err error) {
//	//1
//	t := reflect.TypeOf(v)
//	//1.1
//	if t.Kind() != reflect.Ptr {
//		err = fmt.Errorf("params error;must be a pointer")
//		return
//	}
//	//1.2
//	if t.Elem().Kind() != reflect.Struct {
//		err = fmt.Errorf("params error;must be a struct")
//		return
//	}
//	//2.1
//	b, err := ioutil.ReadFile(fileName) //不要使用os.Open等方法；直接使用ioutil读取文件到变量；避免了一致开启着配置文件
//	if err != nil {
//		return
//	}
//	//2.2
//	sli := strings.Split(string(b), "\r\n")
//	// fmt.Printf("%#v\n", sli)
//	var configName string
//	for lineIndex, sliceEle := range sli {
//		//2.3
//		sliceEle = strings.TrimSpace(sliceEle) //去除每个元素的首尾空格
//		if strings.HasPrefix(sliceEle, ";") || strings.HasPrefix(sliceEle, "#") {
//			continue
//		}
//		if len(sliceEle) == 0 {
//			continue
//		}
//
//		//2.4
//		if strings.HasPrefix(sliceEle, "[") {
//			if sliceEle[0] != '[' || sliceEle[len(sliceEle)-1] != ']' {
//				err = fmt.Errorf("line:%d;syntax error", lineIndex+1)
//				return
//			}
//			//把这一行首尾的[]去掉，取消空格看剩下的长度并拿到内容
//			sectionName := strings.TrimSpace(sliceEle[1 : len(sliceEle)-1])
//			if len(sectionName) == 0 {
//				err = fmt.Errorf("line:%d;syntax error", lineIndex+1)
//				return
//			}
//			for i := 0; i < reflect.TypeOf(v).Elem().NumField(); i++ { //循环IniConfig结构体取出每个field名
//				if reflect.TypeOf(v).Elem().Field(i).Tag.Get("ini") == sectionName {
//					configName = reflect.TypeOf(v).Elem().Field(i).Name
//				}
//			}
//		} else {
//			//2.5
//			if strings.Index(sliceEle, "=") == -1 {
//				err = fmt.Errorf("line:%d;syntax error", lineIndex+1)
//				return
//			}
//			if strings.HasPrefix(sliceEle, "=") {
//				err = fmt.Errorf("line:%d;syntax error", lineIndex+1)
//				return
//			}
//
//			setionField := strings.Split(sliceEle, "=")
//			eleKey := strings.TrimSpace(setionField[0])
//			eleValue := strings.TrimSpace(setionField[1])
//			var configName2 string
//			var configType2 reflect.StructField
//			configValue := reflect.ValueOf(v).Elem().FieldByName(configName) //根据field名拿到对应结构体的Value和Type
//			configType := configValue.Type()
//			for j := 0; j < configValue.NumField(); j++ {
//				field := configType.Field(j)
//				configType2 = field
//				if field.Tag.Get("ini") == eleKey {
//					configName2 = field.Name
//					break
//				}
//			}
//			if len(configName2) == 0 { //跳过没有匹配的key
//				continue
//			}
//			fileValue := configValue.FieldByName(configName2)
//			switch configType2.Type.Kind() {
//			case reflect.String:
//				fileValue.SetString(eleValue)
//			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//				var eleInt int64
//				eleInt, err = strconv.ParseInt(eleValue, 10, 64)
//				if err != nil {
//					err = fmt.Errorf("line：%d;syntax error", lineIndex+1)
//					return
//				}
//				fileValue.SetInt(eleInt)
//			case reflect.Bool:
//				var eleBool bool
//				eleBool, err = strconv.ParseBool(eleValue)
//				if err != nil {
//					err = fmt.Errorf("line：%d;syntax error", lineIndex+1)
//					return
//				}
//				fileValue.SetBool(eleBool)
//			case reflect.Float32, reflect.Float64:
//				var eleFlo float64
//				eleFlo, err = strconv.ParseFloat(eleValue, 64)
//				if err != nil {
//					err = fmt.Errorf("line：%d;syntax error", lineIndex+1)
//					return
//				}
//				fileValue.SetFloat(eleFlo)
//			}
//		}
//	}
//	return
//}

//
//func main() {
//	var iniC IniConfig
//	err := loadIni("./demo.ini", &iniC)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(iniC)
//}
