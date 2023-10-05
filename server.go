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
	"strconv"
	"strings"
	"time"
)

var (
	_                = flag.String("listen", ":5668", "endpoint to listen for messages")
	_                = flag.String("command_file", "/usr/local/var/nagios/rw/nagios.cmd", "target to send updates")
	_                = flag.String("allow", "allowList.txt", "file with allowed certificate DNs to accept")
	_                = flag.Int("max_command_size", 16384, "accept commands of length")
	_                = flag.Int("max_queue_size", 1024, "queue upto this specified number of megabytes")
	allowMap         = make(map[string]struct{})
	_                = flag.Duration("delay", time.Second*5, "time between heartbeats (should match client)")
	delay            time.Duration
	max_command_size int64
	outFile          *os.File
	buf              = Buffer{c: make(chan int, 10)}
	version          string
	prog             = "nsca-tls-server"
)

func main() {
	fmt.Println("NSCA-TLS Server, Version", version, "(https://github.com/pschou/nsca-tls)")
	flag.Parse()
	loadConfig()
	loadTLS()

	var err error
	delay, err = time.ParseDuration(conf["delay"])
	if err != nil {
		log.Fatal("Bad delay", err)
	}
	max_queue_size, err = strconv.ParseInt(conf["max_queue_size"], 10, 64)
	if err != nil {
		log.Fatal("Bad max_queue_size", err)
	}
	max_command_size, err = strconv.ParseInt(conf["max_command_size"], 10, 64)
	if err != nil {
		log.Fatal("Bad max_command_size", err)
	}

	go func() {
	}()

	loadAllows()
	go func() {
		for {
			time.Sleep(time.Minute)
			loadAllows()
		}
	}()

	go handlePipe()

	log.Println("Listening on", conf["listen"], "...")
	listener, err := tls.Listen("tcp", conf["listen"], tlsConfig)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}

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

var lastAllow time.Time

func loadAllows() {
	fileinfo, err := os.Stat(conf["allow"])
	if err != nil {
		log.Fatal(err)
	}
	if fileinfo.ModTime().Sub(lastAllow) == 0 {
		return
	}
	log.Println("Loading allow list")
	newAllowMap := make(map[string]struct{})
	file, err := os.Open(conf["allow"])
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

func handlePipe() {
	log.Println("Opening pipe", conf["command_file"], "if a process is not listening then will wait for a process")
	var err error
	outFile, err = os.OpenFile(conf["command_file"], os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to pipe", conf["command_file"])

	for {
		line, err := buf.ReadString('\n')
		//log.Printf("got line from buffer: %q\n", string(line))
		for {
			_, err = outFile.Write([]byte(line))
			if err == nil {
				break
			}
			//log.Println("error writing to pipe", err)
			time.Sleep(time.Second * 3)
		}
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	//reader := bufio.NewReader(conn)
	log.Println("server: conn: opened", conn.RemoteAddr())

	deadline := time.Now()
	go func() {
		for {
			if time.Now().Sub(deadline) > time.Duration(2)*(delay) {
				conn.Close()
				break
			}
			time.Sleep(delay)
		}
	}()

	var eof bool
	line_buf := make([]byte, max_command_size)
	for !eof {
		line, err := sliceLine(line_buf, conn)
		eof = err == io.EOF
		if err != nil && err != io.EOF {
			//log.Printf("server: conn: read: %s", err)
			return
		}
		deadline = time.Now()
		if len(line) > 1 && line[len(line)-1] == '\n' {
			//log.Printf("write line to buffer: %q\n", string(line))
			_, err = buf.Write(line)
		}
	}
	log.Println("server: conn: closed", conn.RemoteAddr())
}
