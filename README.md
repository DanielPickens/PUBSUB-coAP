PubSub client/server over CoAP(Constrained Application Protocol)


It is a Sub/Pub server and client using CoAP protocol.

Note
It will keep a heart beat signal from client to server if you subscription a topic to remain your UDP port channel.

Install
go get github.com/DanielPickens/CoapPubsub

Usage
Server side example
Create a 1024 buffer for pub/sub server and listen 5683 (default port for CoAP)

package main

import (
	"log"

	. "github.com/DanielPickens/PubsubcoAP"
)

func main() {
	log.Println("Server start....")
	serv := NewCoapPubsubServer(1024)
	serv.ListenAndServe(":5683")
}
Client side example
Create a client to read input flag to send add/remove subscription to server.

package main

import (
	"flag"
	"fmt"
	"log"

	. "github.com/DanielPickens/PubsubcoAP"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 3 {
		fmt.Println("Need more arg: cmd topic msg")
		return
	}

	cmd := flag.Arg(0)
	topic := flag.Arg(1)
	msg := flag.Arg(2)

	fmt.Println(cmd, topic, msg)

	client := NewCoapPubsubClient("localhost:5683")
	if client == nil {
		log.Fatalln("Cannot connect to server, please check your setting.")
	}

	if cmd == "ADDSUB" {
		ch, err := client.AddSub(topic)
		log.Println(" ch:", ch, " err=", err)
		log.Println("Got pub from topic:", topic, " pub:", <-ch)
	}
	log.Println("Done")
}
ex: Add subscription on topic "t1"

client ADDSUB t1 msg
ex: Remove subscription on topic "t1"

client REMSUB t1 msg
ex: Publish "mmmmm" to subscription topic "t1"

client PUB t1 mmmmm
TODO
Hadle for UDP packet lost condition
Gracefully network access
Benchmark
TBD

Inspired
MQTT and CoAP, IoT Protocols
RFC 7252
https://github.com/DanielPickens/go-coap
CoAP an introduction
Project52
It is one of my project 52.


