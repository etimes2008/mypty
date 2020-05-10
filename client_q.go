package main

import (
	"context"
	"crypto/tls"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/lucas-clemente/quic-go"
)

func main() {
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

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr(os.Args[1]+":12345", tlsConf, nil)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for {

		go func() {
			for {
				buf := make([]byte, 1024)
				// if n, err := io.ReadFull(stream, buf); err == nil {
				if n, err := stream.Read(buf); err == nil {
					log.Println("recv:", string(buf[:n]))
					c.WriteMessage(websocket.TextMessage, buf[:n])
				} else {
					log.Fatal(err)
				}
			}
		}()

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("ws recv %d: %s", mt, message)
			stream.Write(message)
		}
	}

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
