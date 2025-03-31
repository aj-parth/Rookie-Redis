package main

import "regexp"

const (
	RookieSet    = "ROOKIE-SET"
	RookieGet    = "ROOKIE-GET"
	RookieDelete = "ROOKIE-DELETE"
	RookieExit   = "ROOKIE-EXIT"
)

var CommandRegexObjMap map[string]*regexp.Regexp
var CommandFuncMap map[string]func(map[string]string, ...string) (string, error)
