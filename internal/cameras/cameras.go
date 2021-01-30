package cameras

import (
	"io"
	"log"
	"net"
)

func Proxy(inAddr, outAddr string, done chan struct{}) {
	ln, err := net.Listen("tcp", inAddr)
	if err != nil {
		log.Printf("could not listen on %s: %s", inAddr, err.Error())
		return
	}
	log.Printf("listening on %s", ln.Addr().String())

	for {
		select {
		case <-done:
			return
		default:
			srcConn, err := ln.Accept()
			if err != nil {
				log.Printf("error accepting connection: %s", err.Error())
			} else {
				log.Printf("new connection from %s", srcConn.RemoteAddr().String())
				dstConn, err := net.Dial("tcp", outAddr)
				if err != nil {
					log.Printf("failed to access destination: %s", err.Error())
					close(done)
					break
				}
				_, err = io.Copy(dstConn, srcConn)
				if err != nil {
					log.Printf("failed to forward connection: %s", err.Error())
					close(done)
				}
			}
		}
	}
}
