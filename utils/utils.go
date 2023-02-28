package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type JSONData struct {
	ID    string
	Topic string
	Data  string
	Time  time.Time
}

func EncodeJSONBody(resp http.ResponseWriter, statusCode int, data interface{}) {
	//marshData, err := json.Marshal(data)
	//if err != nil {
	//	logrus.Errorf("EncodeJSONBody : Error marshing response data interface %v", err)
	//}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	err := json.NewEncoder(resp).Encode(data)
	if err != nil {
		logrus.Errorf("EncodeJSONBody : Error encoding response %v", err)
	}
}

func EncodeJSON200Body(resp http.ResponseWriter, data interface{}) {
	var newJSON = jsoniter.ConfigCompatibleWithStandardLibrary
	err := newJSON.NewEncoder(resp).Encode(data)
	if err != nil {
		logrus.Errorf("EncodeJSON200Body : Error encoding response %v", err)
	}
}

func StorePahoDataToJsonFile(data map[string]interface{}) {

	var jsonText = []byte(`[
        {"ID": "", "Topic": 0, "Data": "", "Time": ""}
    ]`)
	var I JSONData
	err := json.Unmarshal([]byte(jsonText), &I)
	if err != nil {
		fmt.Println(err)
	}

	//I = append(I, JSONData{
	//	ID:    uuid.New().String(),
	//	Topic: data["topic"].(string),
	//	Data:  data["data"].(string),
	//})

	dataString := string(data["data"].([]uint8)[:])

	I = JSONData{
		ID:    uuid.New().String(),
		Topic: data["topic"].(string),
		Data:  dataString,
		Time:  time.Now(),
	}

	result, error := json.Marshal(I)
	if error != nil {
		fmt.Println(error)
	}

	f, erro := os.OpenFile("test.json", os.O_APPEND|os.O_WRONLY, 0666)
	if erro != nil {
		fmt.Println(erro)
	}

	n, err := io.WriteString(f, string(result))
	if err != nil {
		fmt.Println(n, err)
	}
}
