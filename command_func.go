package main

import (
	"errors"
	"fmt"
)

func InitCommandFuncMap() {
	CommandFuncMap = make(map[string]func(map[string]string, ...string) (string, error))
	CommandFuncMap[RookieSet] = SetFunc
	CommandFuncMap[RookieGet] = GetFunc
	CommandFuncMap[RookieDelete] = DeleteFunc
	CommandFuncMap[RookieExit] = ExitFunc
}

var SetFunc = func(memo map[string]string, args ...string) (string, error) {
	key := args[1]
	val := args[2]
	memo[key] = val
	printThis := fmt.Sprintf("%s : %s", key, val)
	return printThis, nil
}

var GetFunc = func(memo map[string]string, args ...string) (string, error) {
	val, ok := memo[args[1]]
	if !ok {
		printThis := fmt.Sprintf("No key found : %s", args[1])
		return printThis, errors.New("memo is empty")
	}
	printThis := fmt.Sprintf("%s : %s", args[1], val)
	return printThis, nil
}

var DeleteFunc = func(memo map[string]string, args ...string) (string, error) {
	key := args[1]
	delete(memo, key)
	printThis := fmt.Sprintf("Removed %s", args[1])
	return printThis, nil
}

var ExitFunc = func(memo map[string]string, args ...string) (string, error) {
	return "exit", nil
}
