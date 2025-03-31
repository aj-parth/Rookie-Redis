package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"strings"
)

type Peer struct {
	conn net.Conn
	memo map[string]string
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		conn: conn,
		memo: make(map[string]string),
	}
}

func (p *Peer) readLoop(localConnNum int64) error {
	reqBuffer := make([]byte, 1024)
	for {
		p.conn.Write([]byte(fmt.Sprintf("Rookie-Redis Connection #%d > ", localConnNum)))
		reqLen, err := p.conn.Read(reqBuffer)
		text := strings.TrimSpace(string(reqBuffer[:reqLen]))
		if err != nil {
			log.Println("readLoop err:", err)
			continue
		}

		slog.Info("Request received: " + string(reqBuffer[:reqLen]))

		commandSep := strings.Split(text, " ")
		if len(commandSep) < 1 {
			log.Println("readLoop err: command is empty")
			p.conn.Write([]byte(fmt.Sprintf("Command is empty \n")))
			continue
		}
		command := commandSep[0]
		regexObj, ok := CommandRegexObjMap[command]
		if !ok {
			log.Println("readLoop err: command not found in CommandRegexObjMap")
			p.conn.Write([]byte(fmt.Sprintf("Invalid command \n")))
			continue
		}
		matchedCommand := regexObj.FindStringSubmatch(text)
		commandFunc, ok := CommandFuncMap[command]
		if !ok {
			log.Println("readLoop err: command not found in CommandFuncMap")
			p.conn.Write([]byte(fmt.Sprintf("Invalid command \n")))
			continue
		}
		res, err := commandFunc(p.memo, matchedCommand...)
		p.conn.Write([]byte(fmt.Sprintf("%s\n", res)))
		if err != nil {
			log.Println("readLoop err:", err)
			continue
		}

		if res == "exit" {
			slog.Info(fmt.Sprintf("Closing connection #%d gracefully", localConnNum))
			p.conn.Close()
			return nil
		}
	}
}
