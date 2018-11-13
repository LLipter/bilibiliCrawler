package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func PrintJson(data interface{}) {
	buf, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buf))
}

func LoadDBConf(path string) (string, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	var dbname DBConf
	err = json.Unmarshal(buf, &dbname)
	if err != nil {
		return "", errors.New("invalid database configuration file")
	}

	var strBuf bytes.Buffer
	strBuf.WriteString(dbname.User)
	strBuf.WriteString(":")
	strBuf.WriteString(dbname.Passwd)
	strBuf.WriteString("@tcp(")
	strBuf.WriteString(dbname.Host)
	strBuf.WriteString(")/")
	strBuf.WriteString(dbname.DBname)
	return strBuf.String(), nil
}
