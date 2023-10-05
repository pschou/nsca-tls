package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

var (
	server    = flag.String("server", "my.server:5568", "endpoint to send messages")
	namedPipe = flag.String("pipe", "/dev/shm/nagios.cmd", "create a listening file here")
	delay     = flag.Duration("delay", time.Second*5, "heartbeat interval")
	file      *os.File
	conn      net.Conn
	version   string
)

func main() {
	fmt.Println("NSCA-TLS Client, Version", version, "(https://github.com/pschou/nsca-tls)")
	flag.Parse()
	loadTLS()
	log.Println("Starting up...")

	err := unix.Mkfifo(*namedPipe, 0666)
	if err != nil {
		log.Fatal("Make named pipe file error:", err, " ", *namedPipe)
	}

	// Handle signals to make sure the fifo file is removed
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	shutdown := func() {
		log.Println("Shutting down")
		if file != nil {
			file.Close()
		}
		if conn != nil {
			conn.Close()
		}
		os.Remove(*namedPipe)
		os.Exit(0)
	}
	go func() {
		<-c
		shutdown()
	}()
	defer shutdown()

	// Loop over testing a write to the connection to ensure the service is available
	go func() {
		for {
			if conn == nil {
				dial()
			} else {
				_, err := conn.Write([]byte{'\n'})
				if err != nil {
					log.Println("Error writing to conn,", err)
					dial()
				}
			}
			if conn == nil {
				log.Println("Retrying connect in", *delay)
			}
			time.Sleep(*delay)
		}
	}()

	log.Println("Waiting for input on fifo socket.")
	for {
		if file != nil {
			file.Close()
		}
		file, err := os.OpenFile(*namedPipe, os.O_RDONLY|os.O_CREATE, os.ModeNamedPipe|0666)
		if err != nil {
			log.Println("Open named pipe file error:", err)
		}

		reader := bufio.NewReader(file)

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				//log.Println("reading err", err)
				break
			}
			for len(line) > 0 && line[len(line)-1] == '\n' {
				if conn != nil {
					_, err = conn.Write(line)
					if err == nil {
						break
					}
				}
				time.Sleep(time.Second)
			}
		}
	}
}

func dial() {
	if conn != nil {
		conn.Close()
	}
	newConn, err := tls.Dial("tcp", *server, tlsConfig)
	if err != nil {
		log.Println("client: dial error: %s", err)
	} else {
		conn = newConn
		log.Println("client: connected to: ", conn.RemoteAddr())
	}
}
