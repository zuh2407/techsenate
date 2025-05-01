package network

import (
    "fmt"
    "net"
)

func StartServer(port string) {
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        fmt.Println("Server error:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Server listening on port", port)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Connection error:", err)
            continue
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Read error:", err)
        return
    }
    fmt.Println("Message received:", string(buffer[:n]))
}