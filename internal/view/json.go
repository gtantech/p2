package view

import (
	jsonv1 "encoding/json"
	json "encoding/json/v2"
	"fmt"
)

func marshalParams(in any) string {
	out, err := json.Marshal(in, jsonv1.FormatDurationAsNano(true))
	if err != nil {
		panic(fmt.Sprintf("failed to marshal json from params: %v", in))
	}
	return string(out)
}

func marshalParamsToJsonField(in any) string {
	js := fmt.Sprintf("{\"json\":%s}", marshalParams(in))
	return js
}
