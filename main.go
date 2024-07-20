package main

import (
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync/atomic"
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
	var connectionNum int64 = 0
	for {
		slog.Info("Accepting connection")
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("Error accepting connection: " + err.Error())
		}

		atomic.AddInt64(&connectionNum, 1)
		localConnNum := connectionNum

		slog.Info(fmt.Sprintf("Starting connection #%d", localConnNum))
		go handleConn(conn, localConnNum)
		/*select {
		case <-ctx.Done():
			fmt.Printf("For conn no #%d Context cancelled: %v\n", localConnNum, ctx.Err())
		case result := <-handleConnChan:
			fmt.Printf("For conn no #%d Received: %s\n", localConnNum, result)
		}*/
	}
}

func timeOutSignal(timeout time.Duration, timeOutSignalChan chan bool) {
	time.Sleep(timeout)
	slog.Info("timed out")
	timeOutSignalChan <- true
}

func handleConn(conn net.Conn, localConnNum int64) {

	reqBuffer := make([]byte, 1024)
	for {
		reqLen, err := conn.Read(reqBuffer)
		text := strings.TrimSpace(string(reqBuffer[:reqLen]))
		if err != nil {
			slog.Error("Error reading request: " + err.Error())
		}

		slog.Info("Request received: " + string(reqBuffer[:reqLen]))

		if text == "Good Bye" {
			slog.Info(fmt.Sprintf("Closing connection #%d gracefully", localConnNum))
			conn.Close()
			return
		}
		conn.Write([]byte(fmt.Sprintf("Rookie-Redis Connection #%d > ", localConnNum)))
	}
}

func main() {
	err := NewServer(Config{listenAddr: ":8080"}).StartServer()
	slog.Info("server stopped")
	if err != nil {
		panic(err)
	}
}
