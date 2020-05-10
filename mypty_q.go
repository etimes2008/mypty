package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"math/big"
	"os/exec"

	"github.com/creack/pty"
	"github.com/lucas-clemente/quic-go"
)

func main() {
	listener, err := quic.ListenAddr("0.0.0.0:12345", generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			stream, err := sess.AcceptStream(context.Background())
			if err != nil {
				panic(err)
			}

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
			go func() {
				io.Copy(ptmx, stream)
			}()
			io.Copy(stream, ptmx)

		}()

	}
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
