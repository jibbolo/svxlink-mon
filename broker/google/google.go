package google

import (
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jibbolo/svxlink-mon/broker"
)

type GoogleBroker struct {
	client mqtt.Client
}

func New() *GoogleBroker {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(time.Minute * 60).Unix(),
		"aud": "test-168217",
	})
	privatekey, err := ioutil.ReadFile("misc/rsa_private.pem")
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

	caCert, err := ioutil.ReadFile("misc/roots.pem")
	if err != nil {
		log.Fatal("3", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	opts := mqtt.NewClientOptions().
		AddBroker("tls://mqtt.googleapis.com:8883").
		SetClientID("projects/test-168217/locations/europe-west1/registries/test/devices/test2").
		SetPassword(tokenString).SetProtocolVersion(4).SetUsername("unused")

	//create and start a client using the above ClientOptions
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("networkerror: %v", token.Error())
	}
	return &GoogleBroker{client}
}

func (g *GoogleBroker) Publish(text string) broker.Token {
	return g.client.Publish("/devices/test2/events", 1, false, text)
}

func (g *GoogleBroker) Close() {
	g.client.Disconnect(250)
}
