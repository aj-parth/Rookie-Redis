package main

import (
	"log/slog"
	"net"
	"time"
)

type Config struct {
	listenAddr string
}

type Server struct {
	config Config
	ln     net.Listener
}

func NewServer(cfg Config) *Server {
	if cfg.listenAddr == "" {
		cfg.listenAddr = ":8080"
	}
	return &Server{
		config: cfg,
	}
}

func (s *Server) StartServer() error {
	ln, err := net.Listen("tcp", s.config.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	slog.Info("Listening on " + s.config.listenAddr)
	s.ln = ln

	s.acceptConn()
	return nil
}

func (s *Server) acceptConn() {
	timeOutSignalChan := make(chan bool, 1)
	slog.Info("Accepting connection")
	conn, err := s.ln.Accept()
	if err != nil {
		slog.Error("Error accepting connection: " + err.Error())
	}
	slog.Info("calling handler")
	go handleConn(conn)
	go timeOutSignal(5*time.Second, timeOutSignalChan)
	isTimeOut := <-timeOutSignalChan
	if isTimeOut {
		slog.Info("closing connection")
		return
	}
}

func timeOutSignal(timeout time.Duration, timeOutSignalChan chan bool) {
	time.Sleep(timeout)
	slog.Info("timed out")
	timeOutSignalChan <- true
}

func handleConn(conn net.Conn) {

	reqBuffer := make([]byte, 1024)
	for {
		reqLen, err := conn.Read(reqBuffer)
		if err != nil {
			slog.Error("Error reading request: " + err.Error())
		}
		slog.Info("Request received: " + string(reqBuffer[:reqLen]))
		conn.Write([]byte("Rookie-Redis > "))
	}
}

func main() {
	err := NewServer(Config{listenAddr: ":8080"}).StartServer()
	slog.Info("server stopped")
	if err != nil {
		panic(err)
	}
}
