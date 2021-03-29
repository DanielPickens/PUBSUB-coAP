package CoapPubsub

import (
	"errors"
	"log"
	"time"

	"github.com/daniel/go-coap"
)

type subConnection struct {
	channel   chan string
	clientCon *coap.Conn
}

type CoapPubsubClient struct {
	msgIndex uint16
	serAddr  string
	subList  map[string]subConnection
}

// Create a pubsub client for CoAP protocol
// It will connect to server and make sure its alive and start heart beat
// To keep udp port open, we will send heart beat event to server every minutes
func NewCoapPubsubClient(servAddr string) *CoapPubsubClient {
	c := new(CoapPubsubClient)
	c.subList = make(map[string]subConnection, 0)
	c.serAddr = servAddr

	//TODO: connection check if any error

	//Start heart beat
	c.msgIndex = GetIPv4Int16() + GetLocalRandomInt()
	log.Println("Init msgID=", c.msgIndex)
	go c.heartBeat()
	return c
}

//AddSub adds Subscription on topic and return a channel for user to wait data
func (c *CoapPubsubClient) AddSub(topic string) (chan string, error) {
	if val, exist := c.subList[topic]; exist {
		//if topic already exist in sub, return and not send to server
		return val.channel, nil
	}

	conn, err := c.sendPubsubReq("ADDSUB", topic)
	if err != nil {
		return nil, err
	}

	subChan := make(chan string)
	go c.waitSubResponse(conn, subChan, topic)

	//Add client connection into member variable for heart beat
	clientConn := subConnection{channel: subChan, clientCon: conn}
	c.subList[topic] = clientConn
	return subChan, nil
}

//RemoveSub removes Subscribetion on topic
func (c *CoapPubsubClient) RemoveSub(topic string) error {
	if _, exist := c.subList[topic]; !exist {
		//if topic not in sub list, return and not send to server
		return nil
	}

	_, err := c.sendPubsubReq("REMSUB", topic)
	return err
}

func (c *CoapPubsubClient) sendPubsubReq(cmd string, topic string) (*coap.Conn, error) {
	Req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: c.getMsgID(),
		Payload:   []byte(""),
	}

	Req.SetOption(coap.ETag, cmd)
	Req.SetPathString(topic)

	conn, err := coap.Dial("udp", c.serAddr)
	if err != nil {
		log.Printf(cmd, ">>Error dialing: %v \n", err)
		return nil, errors.New("Dial failed")
	}
	conn.Send(Req)
	return conn, err
}

func (c *CoapPubsubClient) waitSubResponse(conn *coap.Conn, ch chan string, topic string) {
	var rv *coap.Message
	var err error
	var keepLoop bool
	keepLoop = true
	for keepLoop {
		if rv != nil {
			if err != nil {
				log.Fatalf("Error receiving: %v", err)
			}
			log.Printf("Got %s", rv.Payload)
		}
		rv, err = conn.Receive()

		if err == nil {
			ch <- string(rv.Payload)
		}

		time.Sleep(time.Second)
		if _, exist := c.subList[topic]; !exist {
			//sub topic already remove, leave loop
			log.Println("Loop topic:", topic, " already remove leave loop")
			keepLoop = false
		}
	}
}

func (c *CoapPubsubClient) getMsgID() uint16 {
	c.msgIndex = c.msgIndex + 1
	return c.msgIndex
}

func (c *CoapPubsubClient) heartBeat() {
	log.Println("Starting heart beat loop call")
	hbReq := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: c.getMsgID(),
		Payload:   []byte("Heart beat msg."),
	}

	hbReq.SetOption(coap.ETag, "HB")

	for {

		for k, conn := range c.subList {
			conn.clientCon.Send(hbReq)
			log.Println("Send the heart beat in topic ", k)
		}

		time.Sleep(time.Minute)
	}
}
