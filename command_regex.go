package main

import "regexp"

var CommandRegexStringMap = map[string]string{
	RookieSet:    `^ROOKIE-SET (\S+) (\S+)$`,
	RookieGet:    `^ROOKIE-GET (\S+)$`,
	RookieDelete: `^ROOKIE-DELETE (\S+)$`,
	RookieExit:   `^ROOKIE-EXIT$`,
}

func InitCommandRegexObjMap() {
	CommandRegexObjMap = make(map[string]*regexp.Regexp)
	for k, v := range CommandRegexStringMap {
		CommandRegexObjMap[k] = regexp.MustCompile(v)
	}
}
