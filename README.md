# NSCA-TLS - Nagios Service Check Acceptor TLS Edition

NSCA is a Linux/Unix daemon allows you to integrate passive alerts and checks
from remote machines and applications with Nagios.

Written in GOLang, this service is intended to be setup in a client -> server
manner in which the nsca-tls-server sits on the Nagios side and accepts
incoming connections from the clients.  Every client will have the
nsca-tls-client running and will expose a local FIFO file, /dev/shm/nagios.cmd,
which will act the same way the /usr/local/var/nagios/rw/nagios.cmd file acts
on the Nagios server.

Authentication and authorization is done by two way SSL signed certificates.
The server reloads the certificates every minute from the allowlist file
specified on the command line.

Attempts are made to reconnect if a connection is lost, however back pressure
is sent back to the client when the fifo buffer is unable to write to the
nagios server.

Server startup example:
```
./nsca-tls-server -ca /etc/pki/ca.pem -cert /etc/pki/npe.pem -key /etc/pki/npe.pem
```

Client startup example:
```
./nsca-tls-client -server localhost:5568 -ca /etc/pki/ca.pem -cert /etc/pki/npe.pem -key /etc/pki/npe.pem
```
