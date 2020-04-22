package main

import (
	"crypto/sha1"
	"io"
	"log"
	"os"

	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	// wait for server to become ready
	// time.Sleep(time.Second)

	// dial to the echo server
	if sess, err := kcp.DialWithOptions("61.164.110.198:12345", block, 10, 3); err == nil {
		sess.Write([]byte("w\n"))
		// for {
		// 	data := time.Now().String()
		// 	buf := make([]byte, len(data))
		// 	log.Println("sent:", data)
		// 	if _, err := sess.Write([]byte(data)); err == nil {
		// 		// read back the data
		// 		if _, err := io.ReadFull(sess, buf); err == nil {
		// 			log.Println("recv:", string(buf))
		// 		} else {
		// 			log.Fatal(err)
		// 		}
		// 	} else {
		// 		log.Fatal(err)
		// 	}
		// 	time.Sleep(time.Second)
		// }
		// oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
		// if err != nil {
		// 	panic(err)
		// }
		// defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

		// sess.SetStreamMode(true)
		// sess.SetNoDelay(1, 10, 2, 1)
		sess.SetACKNoDelay(true)

		go func() { io.Copy(sess, os.Stdin) }()
		io.Copy(os.Stdout, sess)
	} else {
		log.Fatal(err)
	}
}
