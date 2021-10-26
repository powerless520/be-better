package reflectUtil

import (
	"reflect"
)

func BsonObjectToMap(obj interface{}) (map[string]interface{}, error) {
	vt := reflect.TypeOf(obj)
	vv := reflect.ValueOf(obj)

	result := map[string]interface{}{}
	for i:=0;i<vt.NumField();i++{
		f := vt.Field(i)
		name := f.Tag.Get("bson")
		if name == "" {
			name = f.Name
		}
		result[name] = vv.FieldByName(f.Name).Interface()
	}
	return result, nil
}

func MapToBsonObject(properties map[string]interface{}, obj interface{}) error {
	vt := reflect.TypeOf(obj).Elem()

	fieldNames := map[string]string{}
	for i:=0;i<vt.NumField();i++{
		f := vt.Field(i)
		name := f.Tag.Get("bson")
		if name == "" {
			name = f.Name
		}
		fieldNames[name] = f.Name
	}
	//vv := reflect.ValueOf(obj)

	ve:= reflect.ValueOf(obj).Elem()
	for k, v := range properties{
		if _, ok := fieldNames[k]; ok {
			//field := vv.FieldByName(fieldNames[k])
			field := ve.FieldByName(fieldNames[k])
			if field.CanSet(){
				field.Set(reflect.ValueOf(v))
			}
		}
	}

	return nil
}

func BsonConvert(src interface{}, obj interface{}) error {
	properties, err := BsonObjectToMap(src)
	if err != nil {
		return err
	}
	return MapToBsonObject(properties, obj)
}
