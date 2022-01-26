package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func handleSteam(connSrc net.Conn, connDest net.Conn) {
	defer connSrc.Close()
	for {
		//var bufferSrc []byte
		bufferSrc := make([]byte, 10 * 1024)
		srcRSize, errSrcR := connSrc.Read(bufferSrc)
		if errSrcR == io.EOF || errSrcR != nil {
			log.Println("Read error:", errSrcR)
			break
		}
		srcWSize, errSrcW := connDest.Write(bufferSrc)
		if errSrcW != nil {
			log.Println("Write error:", errSrcW)
			return
		}
		log.Println("size", srcRSize, srcWSize)
		if srcRSize < 1024{
			//log.Println("bufferSrc", string(bufferSrc))
		}
	}
}

func handleConnection(connSrc net.Conn, rAddr string) {
	connDest, err := net.Dial("tcp", rAddr)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	go handleSteam(connSrc, connDest)
	go handleSteam(connDest, connSrc)
}

func initProxy()  {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// -r 127.0.0.1:8080 -l 127.0.0.1:11020
// http://127.0.0.1:11020
func parseArgs() (string, string) {
	remoteAddr := flag.String("r", "", "remote address")
	localAddr := flag.String("l", ":11020", "listen local address")
	flagVer := flag.Bool("v", false, "print version")

	flag.Parse()

	if *flagVer {
		fmt.Println("v1.0.0")
		os.Exit(0)
	}

	if *remoteAddr == "" {
		flag.Usage()
		os.Exit(0)
	}

	log.Println(*remoteAddr, *localAddr)
	return *remoteAddr, *localAddr
}

/**
 * Proxy
 * http://127.0.0.1:11020
 */
func main() {
	initProxy()
	rAddr, lAddr := parseArgs()

	// Creates servers
	ln, err := net.Listen("tcp", lAddr)
	if err != nil {
		log.Println("Listen error:", err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			return
		}
		go handleConnection(conn, rAddr)
	}
}

