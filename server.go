package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	service   = flag.String("listen", ":5568", "endpoint to listen for messages")
	namedPipe = flag.String("pipe", "/usr/local/var/nagios/rw/nagios.cmd", "target to send updates")
	allowList = flag.String("allow", "allowList.txt", "file with allowed certificate DNs to accept")
	allowMap  = make(map[string]struct{})
	delay     = flag.Duration("delay", time.Second*5, "time between heartbeats (should match client)")
	file      *os.File
	version   string
)

func main() {
	fmt.Println("NSCA-TLS Server, Version", version, "(https://github.com/pschou/nsca-tls)")
	flag.Parse()
	loadTLS()

	var err error
	file, err = os.OpenFile(*namedPipe, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}

	loadAllows()
	go func() {
		for {
			time.Sleep(time.Minute)
			loadAllows()
		}
	}()

	listener, err := tls.Listen("tcp", *service, tlsConfig)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Println("Listening", *service, "...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: failed to accept: %s", err)
			conn.Close()
			continue
		}

		var allowed bool
		tlscon, ok := conn.(*tls.Conn)
		var cn string
		if ok {
			//log.Print("ok=true")
			err = tlscon.Handshake()
			if err != nil {
				log.Println("Handshake error: %s", err)
				conn.Close()
				continue
			}
			//fmt.Printf("tlscon=%#v\n", tlscon)
			state := tlscon.ConnectionState()
			//fmt.Printf("state=%#v\n", state)
			for _, v := range state.PeerCertificates {
				_, inMap := allowMap[v.Subject.CommonName]
				allowed = allowed || inMap
				cn = v.Subject.CommonName
			}
		}
		if allowed {
			log.Printf("server: connection from %q at %s", cn, conn.RemoteAddr())
			go handleClient(conn)
		} else {
			log.Printf("server: connection denied from %q at %s", cn, conn.RemoteAddr())
			conn.Close()
		}
	}
}

func loadAllows() {
	newAllowMap := make(map[string]struct{})
	file, err := os.Open(*allowList)
	if err != nil {
		log.Println("Error reading allow file", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		newAllowMap[strings.TrimSpace(scanner.Text())] = struct{}{}
	}
	allowMap = newAllowMap
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	log.Println("server: conn: opened", conn.RemoteAddr())

	deadline := time.Now()
	go func() {
		for {
			if time.Now().Sub(deadline) > time.Duration(2)*(*delay) {
				conn.Close()
				break
			}
			time.Sleep(*delay)
		}
	}()

	var eof bool
	for !eof {
		line, err := reader.ReadBytes('\n')
		eof = err == io.EOF
		if err != nil && err != io.EOF {
			log.Printf("server: conn: read: %s", err)
			break
		}
		deadline = time.Now()
		if len(line) > 1 && line[len(line)-1] == '\n' {
			_, err = file.Write(line)
		}
		if err != nil {
			if file != nil {
				file.Close()
			}
			file, err = os.OpenFile(*namedPipe, os.O_WRONLY|os.O_APPEND, 0777)
			if err != nil {
				log.Println("Could not re-open", *namedPipe, err)
				log.Println("Warning: metrics dropped")
			}
		}
	}
	log.Println("server: conn: closed", conn.RemoteAddr())
}
