package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

var (
	confFile = flag.String("c", "", "load config from file, for example: /etc/"+prog+".conf")
	conf     = make(map[string]string)
)

func loadConfig() {
	flag.VisitAll(func(f *flag.Flag) {
		//fmt.Println("setting", f.Name, f.DefValue)
		conf[f.Name] = f.DefValue
	})
	if *confFile != "" {
		file, err := os.Open(*confFile)
		if err != nil {
			log.Fatal("Open config file error:", err)
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		i := 0
		for {
			i++
			lineBytes, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal("Read config file error:", err)
			}
			line := strings.TrimSpace(string(lineBytes))
			if len(line) == 0 || line[0] == '#' {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				log.Fatal("Unrecognized config option:", string(line), "on line", i)
			}
			val := strings.TrimSpace(string(parts[1]))
			if len(val) > 2 && val[0] == '"' && val[len(val)-1] == '"' {
				val = val[1 : len(val)-1]
			}
			if _, ok := conf[strings.TrimSpace(string(parts[0]))]; !ok {
				if prog == "nsca-tls-client" && (string(line) == "command_file" || string(line) == "keep_command_file") {
					continue
				}
				log.Println("Unused config option:", string(line), "on line", i)
			} else {
				//fmt.Println("conf-setting", strings.TrimSpace(string(parts[0])), val)
				conf[strings.TrimSpace(string(parts[0]))] = val
			}
		}
	}
	//fmt.Println("walking args")
	flag.Visit(func(f *flag.Flag) {
		//fmt.Println("arg-setting", f.Name, f.Value)
		conf[f.Name] = f.Value.String()
	})
}
