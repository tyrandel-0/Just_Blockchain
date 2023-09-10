package network

import (
	"net"
	"strings"
	"time"
)

func Listen(address string, handle func(Conn, *Package)) Listener {
	splited := strings.Split(address, ":")
	if len(splited) != 2 {
		return nil
	}
	listener, err := net.Listen("tcp", "0.0.0.0:"+splited[1])
	if err != nil {
		return nil
	}
	go serve(listener, handle)
	return Listener(listener)
}

func serve(listener net.Listener, handle func(Conn, *Package)) {
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		go handleConn(conn, handle)
	}
}

func handleConn(conn net.Conn, handle func(Conn, *Package)) {
	defer conn.Close()
	pack := readPackage(conn)
	if pack == nil {
		return
	}
	handle(Conn(conn), pack)
}

func Handle(option int, conn Conn, pack *Package, handle func(*Package) string) bool {
	if pack.Option != option {
		return false
	}
	_, err := conn.Write([]byte(SerializePackage(&Package{
		Option: option,
		Data:   handle(pack),
	}) + EndBytes))
	if err != nil {
		return false
	}
	return true
}

func Send(address string, pack *Package) *Package {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}
	_, err = connection.Write([]byte(SerializePackage(pack) + EndBytes))
	if err != nil {
		return nil
	}

	var (
		res = new(Package)
		ch  = make(chan bool)
	)

	go func() {
		res = readPackage(connection)
		ch <- true
	}()

	select {
	case <-ch:
	case <-time.After(WaitTime * time.Second):
	}

	return res
}

func readPackage(conn Conn) *Package {
	var (
		data   = ""
		size   = uint(0)
		buffer = make([]byte, Buffsize)
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			return nil
		}

		size += uint(length)
		if size > Dmaxsize {
			return nil
		}

		data += string(buffer[:length])
		if strings.Contains(data, EndBytes) {
			data = strings.Split(data, EndBytes)[0]
			break
		}
	}
	return DeserializePackage(data)
}
