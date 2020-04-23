package main

import (
	"crypto/sha1"
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	/*
		key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
		block, _ := kcp.NewAESBlockCrypt(key)

		// wait for server to become ready
		// time.Sleep(time.Second)

		// dial to the echo server
		if sess, err := kcp.DialWithOptions(os.Args[1]+":12345", block, 10, 3); err == nil {
			sess.SetACKNoDelay(true)
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
			// 	// time.Sleep(time.Second)
			// }
			// oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
			// if err != nil {
			// 	panic(err)
			// }
			// defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

			// sess.SetStreamMode(true)
			// sess.SetNoDelay(1, 10, 2, 1)

			go func() { io.Copy(sess, os.Stdin) }()
			io.Copy(os.Stdout, sess)
		} else {
			log.Fatal(err)
		}
	*/
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))
	// static := http.FileServer(http.Dir("."))
	// log.Fatal(http.ListenAndServe(*addr, static))
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	if sess, err := kcp.DialWithOptions(os.Args[1]+":12345", block, 10, 3); err == nil {
		sess.SetACKNoDelay(true)
		sess.Write([]byte("w\n"))

		// op, rd, _ := c.NextReader()
		// println("BinaryMessage", websocket.BinaryMessage, op)
		// wr, _ := c.NextWriter(websocket.TextMessage)
		// go func() {
		// 	// io.Copy(sess, os.Stdin)
		// 	io.Copy(sess, rd)
		// }()
		// // io.Copy(os.Stdout, sess)
		// io.Copy(wr, sess)
		// /*
		go func() {
			for {
				buf := make([]byte, 1024)
				if n, err := io.ReadFull(sess, buf); err == nil {
					log.Println("recv:", string(buf[:n]))
					c.WriteMessage(websocket.TextMessage, buf[:n])
				} else {
					log.Fatal(err)
				}
			}
		}()

		for {
			// op, rd, err := c.NextReader()
			// wr, err := c.NextWriter(op)
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("ws recv %d: %s", mt, message)
			sess.Write(message)
			// err = c.WriteMessage(mt, message)
			// err = c.WriteMessage(mt, []byte("ls\n"))
			// if err != nil {
			// 	log.Println("write:", err)
			// 	break
			// }
		}
		// */
	} else {
		log.Fatal(err)
	}

	// for {
	// 	// op, rd, err := c.NextReader()
	// 	// wr, err := c.NextWriter(op)
	// 	mt, message, err := c.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	log.Printf("recv%d: %s", mt, message)
	// 	sess.Write(message)
	// 	// err = c.WriteMessage(mt, message)
	// 	// err = c.WriteMessage(mt, []byte("ls\n"))
	// 	if err != nil {
	// 		log.Println("write:", err)
	// 		break
	// 	}
	// }
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<link rel="stylesheet" href="https://unpkg.com/xterm@4.5.0/css/xterm.css" />
<script src="https://unpkg.com/xterm@4.5.0/lib/xterm.js"></script>
<script src="./static/xterm-addon-attach.js"></script>
</head>
<body>
<div id="terminal"  style="height: 100%;width:100%" ></div>
<script type="module">  
// import { AttachAddon } from "./static/xterm-addon-attach.js";
window.addEventListener("load", function(evt) {
	var ws;
// var term = new Terminal({
// 	cols: 92,
// 	rows: 22,
// 	cursorBlink: true, // 光标闪烁
// 	cursorStyle: "block", // 光标样式  null | 'block' | 'underline' | 'bar'
// 	scrollback: 800, //回滚
// 	tabStopWidth: 8, //制表宽度
// 	screenKeys: true//
//   });
// term.open(document.getElementById('terminal'));
// term.focus()
//     // term.attachCustomKeyEventHandler(function(ev) {
//     //   //粘贴 ctrl+v
//     //   if (ev.keyCode == 86 && ev.ctrlKey) {
//     //     websocket.send(new TextEncoder().encode("\x00" + this.copy));
//     //   }
// 	// });
// 	term.onData(function (data) {		
// 		term.write(data)
// 		ws.send(data);
//     })

		ws = new WebSocket("{{.}}");
		var terminal = new Terminal();
		terminal.open(document.getElementById('terminal'));
		terminal.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ')
terminal.focus()

		const attachAddon = new AttachAddon.AttachAddon(ws);
terminal.loadAddon(attachAddon);

		// ws.binaryType = "arraybuffer";
        // ws.onopen = function(evt) {
		// 	print("OPEN");
		// 	term.writeln(
		// 		"******************************************************************"
		// 	  );
        // }
        // ws.onclose = function(evt) {
        //     print("CLOSE");
        //     ws = null;
        // }
        // ws.onmessage = function(evt) {
		// 	print("RESPONSE: " + evt.data);
		// 	term.write(evt.data);
        // }
        // ws.onerror = function(evt) {
        //     print("ERROR: " + evt.data);
        // }

    //     ws.close();
});
</script>
</body>
</html>
`))
