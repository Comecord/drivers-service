package tools

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"path/filepath"
	"text/template"
)

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	fmt.Println("Am parsing templates...")

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func GenerateUUID() string {
	return uuid.New().String()
}

func TypeConverter[T any](data any) (*T, error) {
	var result T
	dataJson, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataJson, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func ToString(v interface{}) (string, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ToInt(v interface{}) (int, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return 0, err
	}
	var result int
	err = json.Unmarshal(data, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func ToFloat(v interface{}) (float64, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return 0, err
	}
	var result float64
	err = json.Unmarshal(data, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func ToBool(v interface{}) (bool, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return false, err
	}
	var result bool
	err = json.Unmarshal(data, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func ToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ToSlice(v interface{}) ([]interface{}, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	var result []interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ToSliceMap(v interface{}) ([]map[string]interface{}, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ToSliceString(v interface{}) ([]string, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	var result []string
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ToSliceInt(v interface{}) ([]int, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	var result []int
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ToSliceFloat(v interface{}) ([]float64, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	var result []float64
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ToSliceBool(v interface{}) ([]bool, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	var result []bool
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
