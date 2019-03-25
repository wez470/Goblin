package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
)

type Conf struct {
	Backends []string `yaml:"backends"`
	ListenerPort string `yaml:"listenerPort"`
}

func getConfig() Conf {
	confYaml, err := ioutil.ReadFile("Conf.yaml")
	if err != nil {
		log.Printf("confYaml.Get err: #%v ", err)
	}
	conf := Conf{}
	err = yaml.Unmarshal(confYaml, &conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return conf
}

type Server struct {
	Config Conf
	currBackend uint32
	mu sync.Mutex
}

func NewServer() *Server {
	conf := getConfig()
	return &Server{Config: conf, currBackend: 0}
}

func (s *Server) Run() {
	log.Printf("Initial Server Config: %+v", s.Config)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", s.Config.ListenerPort))
	if err != nil {
		log.Fatalf("Failed to setup TCP listener. error: %v", err)
	}
	log.Printf("Running server on localhost port: %+v", s.Config.ListenerPort)

	for {
		conn, err := ln.Accept()
		log.Printf("Connection accepted")
		if err != nil {
			log.Printf("Failed to handle connection. error: %v", err)
			break
		}
		currBackend := s.currBackend
		s.currBackend = (s.currBackend + 1) % uint32(len(s.Config.Backends))
		go s.handleConnection(conn, s.Config.Backends[currBackend])
	}
	ln.Close()
}

func (s *Server) handleConnection(upConn net.Conn, backendAddr string) {
	log.Printf("Handling connection: %+v", upConn.RemoteAddr())
	defer upConn.Close()

	downConn, err := net.Dial("tcp", backendAddr)
	if err != nil {
		log.Printf("Failed to establish downstream connection.")
		return
	}
	defer downConn.Close()
	log.Printf("Connection to downstream server established: %v", backendAddr)

	pipeTraffic(upConn, downConn)
	log.Printf("Done pipeing traffic. Closing connections")
}

func pipeTraffic(upConn, downConn net.Conn) {
	errChan := make(chan error)

	copy := func(src, dst net.Conn) {
		_, err := io.Copy(dst, src)
		errChan <- err
	}

	go copy(downConn, upConn)
	go copy(upConn, downConn)
	<-errChan
	<-errChan
	return
}

func (s Server) Shutdown() {
	log.Printf("Shutting down server.")
	os.Exit(0)
}

