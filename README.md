# NSCA-TLS - Nagios Service Check Acceptor TLS Edition

NSCA is a Linux/Unix daemon that allows you to integrate passive alerts and checks from remote machines and applications with Nagios.

Written in GOLang, this service is intended to be set up in a client -> server manner in which the nsca-tls-server sits on the Nagios side and accepts incoming connections from the clients. Every client will have the nsca-tls-client running and will expose a local FIFO file, /dev/shm/nagios.cmd, which will act the same way the /usr/local/var/nagios/rw/nagios.cmd file works on the Nagios server.

Two-way SSL signed certificates do authentication and authorization. The server reloads the credentials (certificate common name list) from the allowlist file specified on the command line every minute.

Attempts are made to reconnect if a connection is lost; however back pressure is sent back to the client when the FIFO buffer cannot write to the Nagios server.


Server startup example:
```
./nsca-tls-server -cert ../pki/tests/npe1_cert_DONOTUSE.pem -key ../pki/tests/npe1_key_DONOTUSE.pem -ca ../pki/tests/ca_cert_DONOTUSE.pem -command_file /dev/shm/nagios.cmd -c nsca-tls-server.cfg
```

Client startup example:
```
./nsca-tls-client -server 127.0.0.1 -port 5668 -cert ../pki/tests/npe2_cert_DONOTUSE.pem -key ../pki/tests/npe2_key_DONOTUSE.pem -ca ../pki/tests/ca_cert_DONOTUSE.pem -command_file /dev/shm/client.cmd
```
