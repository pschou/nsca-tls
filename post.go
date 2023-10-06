package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	_       = flag.String("server", "my.server", "endpoint host to send messages")
	_       = flag.Int("port", 5668, "endpoint port to send messages")
	_       = flag.Duration("delay", time.Second*5, "heartbeat interval")
	file    *os.File
	conn    net.Conn
	version string
	prog    = "nsca-tls-client"
	delay   time.Duration
)

func main() {
	fmt.Println("NSCA-TLS Post, Version", version, "(https://github.com/pschou/nsca-tls)")
	flag.Parse()
	loadConfig()
	loadTLS()

	var err error
	delay, err = time.ParseDuration(conf["delay"])
	if err != nil || delay < time.Second {
		log.Fatal("Bad delay", err)
	}

	dial()

	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadBytes('\n')
		//log.Printf("read line: %q\n", string(line))
		if err != nil {
			//log.Println("reading err", err)
			break
		}
		for len(line) > 0 && line[len(line)-1] == '\n' {
			if c := conn; c != nil {
				//log.Printf("write line to conn: %q\n", string(line))
				_, err = c.Write(line)
				if err == nil {
					break
				}
			}
			//log.Println("retrying")
			time.Sleep(time.Second)
		}
	}
}

func dial() {
	if conn != nil {
		conn.Close()
		conn = nil
	}
	newConn, err := tls.Dial("tcp", net.JoinHostPort(conf["server"], conf["port"]), tlsConfig)
	if err != nil {
		log.Println("client: dial error:", err)
	} else {
		go func() {
			// Keep heartbeating (by sending an empty line) until the connection is closed
			for err == nil {
				_, err = newConn.Write([]byte("\n"))
				time.Sleep(delay)
			}
			newConn.Close()
		}()
		conn = newConn
		log.Println("client: connected to:", conn.RemoteAddr())
	}
}
