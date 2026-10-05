package jsonmodels

import (
	jsonv1 "encoding/json"
	json "encoding/json/v2"
	"fmt"
)

const (
	JsonMarshalField string = "json" // {\"json\":%s}
)

func MarshalParams(in any) string {
	out, err := json.Marshal(in, jsonv1.FormatDurationAsNano(true))
	if err != nil {
		panic(fmt.Sprintf("failed to marshal json from params: %v", in))
	}
	return string(out)
}

func MarshalParamsToJsonField(in any) string {
	js := fmt.Sprintf("{\"%s\":%s}", JsonMarshalField, MarshalParams(in))
	return js
}
