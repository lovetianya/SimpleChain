package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s type(server/client) [port]\r\n", os.Args[0])
		os.Exit(1)
	}
	hello()
	if os.Args[1] == "client" {
		go client()
	} else {
		go serve(os.Args[2])
	}
	for {
	}
}

func client() {
	var local_addr string
	var remote_addr string
	fmt.Print("Enter local host:port ")
	fmt.Scanf("%s", &local_addr)
	fmt.Print("Enter remote host:port ")
	fmt.Scanf("%s", &remote_addr)
	raddr, err := net.ResolveTCPAddr("tcp", remote_addr)
	checkErr(err)
	laddr, err := net.ResolveTCPAddr("tcp", local_addr)
	checkErr(err)

	conn, err := net.DialTCP("tcp", laddr, raddr)
	conn.SetKeepAlive(true)
	for {
		// var input string
		fmt.Printf(">>> ")
		// fmt.Scanf("%s", &input)
		reader := bufio.NewReader(os.Stdin)
		input, _, err := reader.ReadLine()
		checkErr(err)
		if string(input) == "q" || string(input) == "quit" {
			break
		}
		_, err = conn.Write(append(input, '\n'))
		checkErr(err)
		var buf [512]byte

		n, err := conn.Read(buf[0:])
		checkErr(err)
		fmt.Println(string(buf[0:n]))
		// response, _ := readFully(conn)
		// fmt.Println(string(response))
	}
	conn.Close()
	os.Exit(0)
}

func serve(port string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+port)
	checkErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
	os.Exit(0)
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// var buf [512]byte
	for {
		// n, err := conn.Read(buf[0:])
		// buf, err := readFully(conn)
		buf, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil && err != io.EOF || len(buf) == 0 {
			return
		}
		fmt.Println(string(buf[:len(buf)-1]))
		_, err2 := conn.Write([]byte("success!"))
		if err2 != nil {
			return
		}
	}
}

func hello() {
	fmt.Printf("\r\nHello! This is simple chain\r\n")
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s\r\n", err.Error())
		os.Exit(2)
	}
}

func readFully(conn net.Conn) ([]byte, error) {
	defer conn.Close()

	result := bytes.NewBuffer(nil)
	var buf [512]byte

	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return result.Bytes(), nil
}

func addEOF(input []byte) []byte {
	buf := make([]byte, len(input)+1)
	copy(buf, input)
	buf[len(input)] = '\n'

	return buf
}
