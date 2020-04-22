package main

import (
	"crypto/sha1"
	"io"
	"log"
	"os/exec"

	"github.com/creack/pty"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	if listener, err := kcp.ListenWithOptions("0.0.0.0:12345", block, 10, 3); err == nil {
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Fatal(err)
			}
			go handleEcho(s)
		}
	} else {
		log.Fatal(err)
	}
}

// handleEcho send back everything it received
func handleEcho(conn *kcp.UDPSession) {
	// conn.SetNoDelay(1, 10, 2, 1)
	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return
	}
	// Make sure to close the pty at the end.
	defer func() {
		log.Println("client disconnect")
		ptmx.Close()
	}() // Best effort.

	// Handle pty size.
	// ch := make(chan os.Signal, 1)
	// signal.Notify(ch, syscall.SIGWINCH)
	// go func() {
	// 	for range ch {
	// 		if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
	// 			log.Printf("error resizing pty: %s", err)
	// 		}
	// 	}
	// }()
	// ch <- syscall.SIGWINCH // Initial resize.

	// // Set stdin in raw mode.
	// oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() { io.Copy(ptmx, conn) }()
	io.Copy(conn, ptmx)

	// go func(conn *kcp.UDPSession, ptmx *os.File) {
	// 	buf := make([]byte, 4096)
	// 	n, err := ptmx.Read(buf)
	// 	log.Println("pty", string(buf[:n]))
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	conn.Write(buf[:n])
	// }(conn, ptmx)

	// buf := make([]byte, 4096)
	// for {
	// 	n, err := conn.Read(buf)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	log.Println("client", string(buf[:n]), buf[:n])
	// 	// n, err = conn.Write(buf[:n])
	// 	// if err != nil {
	// 	// 	log.Println(err)
	// 	// 	return
	// 	// }
	// 	ptmx.Write(buf[:n])
	// }
}
