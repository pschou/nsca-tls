# NSCA-TLS - Nagios Service Check Acceptor TLS Edition

NSCA is a Linux/Unix daemon that allows you to integrate passive alerts and checks from remote machines and applications with Nagios.

Written in GOLang, this service is intended to be set up in a client -> server manner in which the nsca-tls-server sits on the Nagios side and accepts incoming connections from the clients. Every client will have the nsca-tls-client running and will expose a local FIFO file, /dev/shm/nagios.cmd, which will act the same way the /usr/local/var/nagios/rw/nagios.cmd file works on the Nagios server.

Two-way SSL signed certificates do authentication and authorization. The server reloads the credentials (certificate common name list) from the allowlist file specified on the command line every minute.

Attempts are made to reconnect if a connection is lost; however back pressure is sent back to the client when the FIFO buffer cannot write to the Nagios server.

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
NSCA-TLS Server, Version 0.1.20231005.1458 (https://github.com/pschou/nsca-tls)
Usage of ./nsca-tls-server:
  -allow string
        file with allowed certificate DNs to accept (default "allowList.txt")
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
        endpoint to listen for messages (default ":5568")
  -max_command_size int
        accept commands of length (default 16384)
  -max_queue_size int
        queue upto this specified number of megabytes (default 1024)
  -tls_ciphers string
        Available ciphers to pick from:
                # TLS 1.0 - 1.2 cipher suites.
                RSA_WITH_RC4_128_SHA
                RSA_WITH_3DES_EDE_CBC_SHA
                RSA_WITH_AES_128_CBC_SHA
                RSA_WITH_AES_256_CBC_SHA
                RSA_WITH_AES_128_CBC_SHA256
                RSA_WITH_AES_128_GCM_SHA256
                RSA_WITH_AES_256_GCM_SHA384
                ECDHE_ECDSA_WITH_RC4_128_SHA
                ECDHE_ECDSA_WITH_AES_128_CBC_SHA
                ECDHE_ECDSA_WITH_AES_256_CBC_SHA
                ECDHE_RSA_WITH_RC4_128_SHA
                ECDHE_RSA_WITH_3DES_EDE_CBC_SHA
                ECDHE_RSA_WITH_AES_128_CBC_SHA
                ECDHE_RSA_WITH_AES_256_CBC_SHA
                ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
                ECDHE_RSA_WITH_AES_128_CBC_SHA256
                ECDHE_RSA_WITH_AES_128_GCM_SHA256
                ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
                ECDHE_RSA_WITH_AES_256_GCM_SHA384
                ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
                ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256
                ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256

                # TLS 1.3 cipher suites.
                AES_128_GCM_SHA256
                AES_256_GCM_SHA384
                CHACHA20_POLY1305_SHA256

         (default "RSA_WITH_AES_128_CBC_SHA:RSA_WITH_AES_128_GCM_SHA256:RSA_WITH_AES_256_GCM_SHA384")
```

## Configuration options for the client
```
$ ./nsca-tls-client -h
NSCA-TLS Client, Version 0.1.20231005.1458 (https://github.com/pschou/nsca-tls)
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
        endpoint port to send messages (default 5568)
  -server string
        endpoint host to send messages (default "my.server")
  -tls_ciphers string
        Available ciphers to pick from:
                # TLS 1.0 - 1.2 cipher suites.
                RSA_WITH_RC4_128_SHA
                RSA_WITH_3DES_EDE_CBC_SHA
                RSA_WITH_AES_128_CBC_SHA
                RSA_WITH_AES_256_CBC_SHA
                RSA_WITH_AES_128_CBC_SHA256
                RSA_WITH_AES_128_GCM_SHA256
                RSA_WITH_AES_256_GCM_SHA384
                ECDHE_ECDSA_WITH_RC4_128_SHA
                ECDHE_ECDSA_WITH_AES_128_CBC_SHA
                ECDHE_ECDSA_WITH_AES_256_CBC_SHA
                ECDHE_RSA_WITH_RC4_128_SHA
                ECDHE_RSA_WITH_3DES_EDE_CBC_SHA
                ECDHE_RSA_WITH_AES_128_CBC_SHA
                ECDHE_RSA_WITH_AES_256_CBC_SHA
                ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
                ECDHE_RSA_WITH_AES_128_CBC_SHA256
                ECDHE_RSA_WITH_AES_128_GCM_SHA256
                ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
                ECDHE_RSA_WITH_AES_256_GCM_SHA384
                ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
                ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256
                ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256

                # TLS 1.3 cipher suites.
                AES_128_GCM_SHA256
                AES_256_GCM_SHA384
                CHACHA20_POLY1305_SHA256

         (default "RSA_WITH_AES_128_CBC_SHA:RSA_WITH_AES_128_GCM_SHA256:RSA_WITH_AES_256_GCM_SHA384")
```

## Example client conf file:

nsca-tls-client.cfg
```
cert = /etc/pki/server.crt
key = /etc/pki/server.key
ca = /etc/pki/ca-trust/extracted
server = monitoring.example.com
tls_ciphers = "RSA_WITH_AES_128_CBC_SHA:RSA_WITH_AES_128_GCM_SHA256:RSA_WITH_AES_256_GCM_SHA384"
delay = 5s
port = 5668
```

## Example server conf file:

nsca-tls-server.cfg
```
command_file = "/rw/nagios.cmd"
listen = ":5668"
tls_ciphers = "RSA_WITH_AES_128_CBC_SHA:RSA_WITH_AES_128_GCM_SHA256:RSA_WITH_AES_256_GCM_SHA384"
delay = 5s
timeout = 10s
```
