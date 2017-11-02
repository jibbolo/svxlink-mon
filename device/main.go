package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	//import the Paho Go mqtt library
	"github.com/dgrijalva/jwt-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//define a function for the default message handler
var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	now := time.Now()
	mqtt.ERROR = log.New(os.Stderr, "ERR ", log.LstdFlags)
	mqtt.DEBUG = log.New(os.Stdout, "*** ", log.LstdFlags)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(time.Minute * 60).Unix(),
		"aud": "test-168217",
	})
	privatekey, err := ioutil.ReadFile("rsa_private.pem")
	if err != nil {
		log.Fatal("1", err)
	}
	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatekey)
	if err != nil {
		log.Fatal("1.1 ", err)
	}
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(rsaKey)
	if err != nil {
		log.Fatal("2 ", err)
	}

	caCert, err := ioutil.ReadFile("roots.pem")
	if err != nil {
		log.Fatal("3", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	opts := mqtt.NewClientOptions().
		AddBroker("tls://mqtt.googleapis.com:8883").
		SetClientID("projects/test-168217/locations/europe-west1/registries/test/devices/test1").
		SetDefaultPublishHandler(f).
		SetPassword(tokenString).SetProtocolVersion(4).SetUsername("unused")

	//create and start a client using the above ClientOptions
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("networkerror: %v", token.Error())
	}

	// //subscribe to the topic /go-mqtt/sample and request messages to be delivered
	// //at a maximum qos of zero, wait for the receipt to confirm the subscription
	// if token := c.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
	// 	fmt.Println(token.Error())
	// 	os.Exit(1)
	// }

	//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		println(text)
		if !c.IsConnected() {
			println("not conn")
			continue
		}
		token := c.Publish("/devices/test1/events", 0, false, text)
		if token.Wait() && token.Error() != nil {
			log.Printf("publisherror: %v", token.Error())
		}
	}

	time.Sleep(3 * time.Second)

	// //unsubscribe from /go-mqtt/sample
	// if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
	// 	fmt.Println(token.Error())
	// 	os.Exit(1)
	// }

	c.Disconnect(250)
}
