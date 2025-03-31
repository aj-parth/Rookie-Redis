package main

import (
	"fmt"
	"log/slog"
	"net"
	"sync/atomic"
	"time"
)

type Config struct {
	listenAddr string
}

type Server struct {
	config      Config
	ln          net.Listener
	peers       map[*Peer]bool
	addPeerChan chan *Peer
	quitChan    chan struct{}
}

func NewServer(cfg Config) *Server {
	if cfg.listenAddr == "" {
		cfg.listenAddr = ":8082"
	}
	return &Server{
		config:      cfg,
		peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		quitChan:    make(chan struct{}),
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

	go s.loop()
	s.acceptConn()
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case <-s.quitChan:
			return
		case p := <-s.addPeerChan:
			s.peers[p] = true
		}
	}
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
		go s.handleConn(conn, localConnNum)
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

func (s *Server) handleConn(conn net.Conn, localConnNum int64) {

	peer := NewPeer(conn)
	s.addPeerChan <- peer
	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())
	err := peer.readLoop(localConnNum)
	if err != nil {
		slog.Error("readLoop", "err", err.Error())
		return
	}
	return
}

func main() {
	InitCommandRegexObjMap()
	InitCommandFuncMap()
	err := NewServer(Config{listenAddr: ":8082"}).StartServer()
	slog.Info("server stopped")
	if err != nil {
		panic(err)
	}
}
