# NSCA-TLS - Nagios Service Check Acceptor TLS Edition

NSCA is a Linux/Unix daemon that allows you to integrate passive alerts and
checks from remote machines and applications with Nagios.  This TLS edition
uses certificates [RFC5280](https://datatracker.ietf.org/doc/html/rfc5280) to
authenticate peers and supports TLS 1.2
[RFC5246](https://datatracker.ietf.org/doc/html/rfc5246) and TLS 1.3
[RFC8446](https://datatracker.ietf.org/doc/html/rfc8446).

Written in GOLang, this service is intended to be set up in a client -> server
TCP & TLS handshake method in which the nsca-tls-server sits on the Nagios side
and accepts incoming connections from the clients on a TCP port. Every client
will either have the nsca-tls-client running and will expose a local FIFO file
(for example: /dev/shm/nagios.cmd which will act the same way the
/usr/local/var/nagios/rw/nagios.cmd file works on the Nagios server) or use the
nsca-tls-post utility which takes input from the standard input for can sends a
one off connection with one or more lines of metrics/commands.

Note: Every command being sent must end with a newline character, if the
message is terminated before the final newline charater is sent the last line
will not be forwarded.

Two-way SSL signed certificates do authentication and authorization. The server
reloads the credentials (certificate common name list) from the allowlist file
specified on the command line every minute if a change is detected.

Attempts are made to reconnect if a connection is lost; however back pressure
is sent back to the client when the connection cannot be made to the
nsca-tls-server.  One will need to account for this when interacting with the
fifo as cronjobs can become "stuck" in a waiting-to-write state if the server
endpoint is not available.  Some options may be to call the metric utility
through a `timeout (command)` on the command line to terminate the write if a
certain amount of time is not made, or to put the metrics collection in a loop.


## Example command line usage

Server startup example:
```
./nsca-tls-server -cert ../pki/tests/npe1_cert_DONOTUSE.pem -key ../pki/tests/npe1_key_DONOTUSE.pem -ca ../pki/tests/ca_cert_DONOTUSE.pem -command_file /dev/shm/nagios.cmd -c nsca-tls-server.cfg
```

Client startup example:
```
./nsca-tls-client -server 127.0.0.1 -port 5668 -cert ../pki/tests/npe2_cert_DONOTUSE.pem -key ../pki/tests/npe2_key_DONOTUSE.pem -ca ../pki/tests/ca_cert_DONOTUSE.pem -command_file /dev/shm/client.cmd
```

## Configuration options for the server

```
$ ./nsca-tls-server -h
NSCA-TLS Server
Usage of ./nsca-tls-server:
  -allow string
        file with allowed certificate DNs to accept (default "/etc/nsca-tls-allow.txt")
  -c string
        load config from file, for example: /etc/nsca-tls-server.conf
  -ca string
        pem encoded certificate authority chains (default "/etc/pki/tls/certs/ca-bundle.crt")
  -cert string
        pem encoded certificate file (default "/etc/pki/server.pem")
  -command_file string
        target to send updates (default "/usr/local/var/nagios/rw/nagios.cmd")
  -delay duration
        time between heartbeats (should match client) (default 5s)
  -key string
        pem encoded unencrypted key file (default "/etc/pki/server.pem")
  -listen string
        endpoint to listen for messages (default ":5668")
  -max_command_size int
        accept commands of length (default 16384)
  -max_queue_size string
        queue up to this specified number of bytes (default "100MB")
  -tls_ciphers string
        Available ciphers to pick from:
        - RSA_WITH_AES_128_CBC_SHA
        - RSA_WITH_AES_256_CBC_SHA
        - RSA_WITH_AES_128_GCM_SHA256
        - RSA_WITH_AES_256_GCM_SHA384
        - AES_128_GCM_SHA256
        - AES_256_GCM_SHA384
        - CHACHA20_POLY1305_SHA256
        - ECDHE_ECDSA_WITH_AES_128_CBC_SHA
        - ECDHE_ECDSA_WITH_AES_256_CBC_SHA
        - ECDHE_RSA_WITH_AES_128_CBC_SHA
        - ECDHE_RSA_WITH_AES_256_CBC_SHA
        - ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
        - ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
        - ECDHE_RSA_WITH_AES_128_GCM_SHA256
        - ECDHE_RSA_WITH_AES_256_GCM_SHA384
        - ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256
        - ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256
         (default "AES_128_GCM_SHA256:AES_256_GCM_SHA384")
```

## Configuration options for the client
```
$ ./nsca-tls-client -h
NSCA-TLS Client
Usage of ./nsca-tls-client:
  -c string
        load config from file, for example: /etc/nsca-tls-client.conf
  -ca string
        pem encoded certificate authority chains (default "/etc/pki/tls/certs/ca-bundle.crt")
  -cert string
        pem encoded certificate file (default "/etc/pki/server.pem")
  -command_file string
        create a listening file here (default "/dev/shm/nagios.cmd")
  -delay duration
        heartbeat interval (default 5s)
  -key string
        pem encoded unencrypted key file (default "/etc/pki/server.pem")
  -port int
        endpoint port to send messages (default 5668)
  -server string
        endpoint host to send messages (default "my.server")
  -tls_ciphers string
        Available ciphers to pick from:
        - RSA_WITH_AES_128_CBC_SHA
        - RSA_WITH_AES_256_CBC_SHA
        - RSA_WITH_AES_128_GCM_SHA256
        - RSA_WITH_AES_256_GCM_SHA384
        - AES_128_GCM_SHA256
        - AES_256_GCM_SHA384
        - CHACHA20_POLY1305_SHA256
        - ECDHE_ECDSA_WITH_AES_128_CBC_SHA
        - ECDHE_ECDSA_WITH_AES_256_CBC_SHA
        - ECDHE_RSA_WITH_AES_128_CBC_SHA
        - ECDHE_RSA_WITH_AES_256_CBC_SHA
        - ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
        - ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
        - ECDHE_RSA_WITH_AES_128_GCM_SHA256
        - ECDHE_RSA_WITH_AES_256_GCM_SHA384
        - ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256
        - ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256
         (default "AES_128_GCM_SHA256:AES_256_GCM_SHA384")
```

## Configuration options for the posting util
```
$ ./nsca-tls-post -h
NSCA-TLS Post
Usage of ./nsca-tls-post:
  -c string
        load config from file, for example: /etc/nsca-tls-client.conf
  -ca string
        pem encoded certificate authority chains (default "/etc/pki/tls/certs/ca-bundle.crt")
  -cert string
        pem encoded certificate file (default "/etc/pki/server.pem")
  -delay duration
        heartbeat interval (default 5s)
  -key string
        pem encoded unencrypted key file (default "/etc/pki/server.pem")
  -port int
        endpoint port to send messages (default 5668)
  -server string
        endpoint host to send messages (default "my.server")
  -tls_ciphers string
        Available ciphers to pick from:
        - RSA_WITH_AES_128_CBC_SHA
        - RSA_WITH_AES_256_CBC_SHA
        - RSA_WITH_AES_128_GCM_SHA256
        - RSA_WITH_AES_256_GCM_SHA384
        - AES_128_GCM_SHA256
        - AES_256_GCM_SHA384
        - CHACHA20_POLY1305_SHA256
        - ECDHE_ECDSA_WITH_AES_128_CBC_SHA
        - ECDHE_ECDSA_WITH_AES_256_CBC_SHA
        - ECDHE_RSA_WITH_AES_128_CBC_SHA
        - ECDHE_RSA_WITH_AES_256_CBC_SHA
        - ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
        - ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
        - ECDHE_RSA_WITH_AES_128_GCM_SHA256
        - ECDHE_RSA_WITH_AES_256_GCM_SHA384
        - ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256
        - ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256
         (default "AES_128_GCM_SHA256:AES_256_GCM_SHA384")
  -v    turn on verbose
```

## Example client conf file:

/etc/nsca-tls-client.conf
```
# This is a comment line
ca = /etc/pki/tls/certs/ca-bundle.crt
cert = /etc/pki/server.pem
command_file = "/dev/shm/nagios.cmd"
delay = 5s
key = /etc/pki/server.pem
port = 5668
server = monitoring.example.com
tls_ciphers = AES_128_GCM_SHA256:AES_256_GCM_SHA384
```

## Example server conf file:

/etc/nsca-tls-server.conf
```
# This is a comment line
allow = /etc/nsca-tls-allow.txt
ca = /etc/pki/ca-trust/extracted
cert = /etc/pki/server.pem
command_file = "/usr/local/var/nagios/rw/nagios.cmd"
delay = 5s
key = /etc/pki/server.pem
listen = ":5668"
max_command_size = 16384
max_queue_size = 100MB
tls_ciphers = AES_128_GCM_SHA256:AES_256_GCM_SHA384
```
