package main

import (
	"encoding/json"
	"fmt"
)

/**
* type:
 */
type Command struct {
	Type      int               `json:"type"`
	Parameter map[string]string `json:"parameter"`
}

func Parse(commandJson string) *Command {
	var command Command
	err := json.Unmarshal([]byte(commandJson), &command)
	if err != nil {
		fmt.Println("json unmarshal err:", err)
		return nil
	}
	return &command
}

func Build(command *Command) string {
	commandJson, err := json.Marshal(command)
	if err != nil {
		fmt.Println("json marshal err:", err)
		return ""
	}
	return string(commandJson)
}
