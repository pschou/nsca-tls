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
	_       = flag.String("server", "my.server", "endpoint host to send messages")
	_       = flag.Int("port", 5668, "endpoint port to send messages")
	_       = flag.String("command_file", "/dev/shm/nagios.cmd", "create a listening file here")
	_       = flag.Duration("delay", time.Second*5, "heartbeat interval")
	file    *os.File
	conn    net.Conn
	version string
	prog    = "nsca-tls-client"
	delay   time.Duration
)

func main() {
	fmt.Println("NSCA-TLS Client, Version", version, "(https://github.com/pschou/nsca-tls)")
	flag.Parse()
	loadConfig()
	loadTLS()

	var err error
	delay, err = time.ParseDuration(conf["delay"])
	if err != nil || delay < time.Second {
		log.Fatal("Bad delay", err)
	}

	log.Println("Starting up...")
	err = unix.Mkfifo(conf["command_file"], 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}
	keepFIFO := err != nil
	// log.Fatal("Make named pipe file error:", err, " ", conf["command_file"])

	// Handle signals to make sure the fifo file is removed
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	shutdown := func() {
		log.Println("Shutting down")
		if file != nil {
			file.Close()
		}
		if c := conn; c != nil {
			c.Close()
		}
		if keepFIFO {
			os.Remove(conf["command_file"])
		}
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
				log.Println("Connecting")
				dial()
			} else {
				// Keep heartbeating (by sending an empty line) until the connection is closed
				_, err := conn.Write([]byte{'\n'})
				if err != nil {
					log.Println("Error writing to conn,", err)
					dial()
				}
			}
			time.Sleep(delay)
		}
	}()

	log.Println("Waiting for input on fifo socket:", conf["command_file"])
	for {
		if file != nil {
			file.Close()
		}
		file, err := os.OpenFile(conf["command_file"], os.O_RDONLY|os.O_CREATE, os.ModeNamedPipe|0666)
		if err != nil {
			log.Println("Open named pipe file error:", err)
		}

		reader := bufio.NewReader(file)

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
		conn = newConn
		log.Println("client: connected to:", conn.RemoteAddr())
	}
}
