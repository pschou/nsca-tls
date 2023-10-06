package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"unicode"

	_ "unsafe"
)

var (
	_ = flag.String("cert", "/etc/pki/server.pem", "pem encoded certificate file")
	_ = flag.String("key", "/etc/pki/server.pem", "pem encoded unencrypted key file")
	_ = flag.String("ca", "/etc/pki/tls/certs/ca-bundle.crt", "pem encoded certificate authority chains")
	_ = flag.String("tls_ciphers", "AES_128_GCM_SHA256:AES_256_GCM_SHA384", cipher_list)

	tlsConfig *tls.Config
)

//go:linkname defaultCipherSuitesTLS13  crypto/tls.defaultCipherSuitesTLS13
var defaultCipherSuitesTLS13 []uint16

//go:linkname defaultCipherSuitesTLS13NoAES crypto/tls.defaultCipherSuitesTLS13NoAES
var defaultCipherSuitesTLS13NoAES []uint16

func loadTLS() {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(conf["cert"], conf["key"])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Loaded certs from", conf["cert"], conf["key"])

	// Load CA cert
	caCert, err := ioutil.ReadFile(conf["ca"])
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cipherList, minVer, maxVer := buildCipherList()
	//fmt.Println("cipher list:", cipherList)

	// Strip out ciphers which are not requested
	defaultCipherSuitesTLS13 = intersect(defaultCipherSuitesTLS13, cipherList)
	defaultCipherSuitesTLS13NoAES = intersect(defaultCipherSuitesTLS13NoAES, cipherList)

	// Setup TLS Config
	tlsConfig = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		ClientCAs:          caCertPool,
		InsecureSkipVerify: false,
		ClientAuth:         tls.RequireAndVerifyClientCert,

		MinVersion:               minVer,
		MaxVersion:               maxVer,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites:             cipherList,
	}
	tlsConfig.BuildNameToCertificate()
}

func intersect(in, match []uint16) (out []uint16) {
	for _, a := range in {
		for _, b := range match {
			if a == b {
				out = append(out, a)
				break
			}
		}
	}
	return
}

// Make the map and list for building allowed ciphers
var cipher_list = func() string {
	list := "Available ciphers to pick from:\n"
	for _, avail := range tls.CipherSuites() {
		shortName := strings.TrimPrefix(avail.Name, "TLS_")
		list = list + "- " + shortName + "\n"
	}
	return list
}()

func buildCipherList() (cipherList []uint16, minVer, maxVer uint16) {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_'
	}
	minVer = 0xffff

	for _, testCipher := range strings.FieldsFunc(conf["tls_ciphers"], f) {
		testCipher = strings.TrimSpace(testCipher)
		var found bool
		for _, c := range tls.CipherSuites() {
			shortName := strings.TrimPrefix(c.Name, "TLS_")
			if testCipher == shortName {
				found = true
				cipherList = append(cipherList, c.ID)
				if first := c.SupportedVersions[0]; first < minVer {
					minVer = first
				}
				if last := c.SupportedVersions[len(c.SupportedVersions)-1]; last > maxVer {
					maxVer = last
				}
				break
			}
		}
		if !found {
			log.Fatal("Unknown cipher: ", testCipher)
		}
	}
	return
}
