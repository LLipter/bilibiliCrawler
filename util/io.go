package util

import (
	"encoding/json"
	"fmt"
)

func PrintJson(data interface{}) {
	buf, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buf))
}
