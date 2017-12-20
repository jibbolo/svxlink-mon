package broker

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const googleBrokerEndpoint = "tls://mqtt.googleapis.com:8883"
const googleUsername = "unused"
const googleMQTTVersion = 4
const clientIDFormat = "projects/%s/locations/%s/registries/%s/devices/%s"

type GoogleBroker struct {
	basicBroker
}

func NewGoogleBroker(region, projectID, registryID, deviceID, certRootPath, privateKeyPath string) (*GoogleBroker, error) {

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(time.Minute * 60).Unix(),
		"aud": projectID,
	})
	privatekey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatekey)
	if err != nil {
		return nil, err
	}
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(rsaKey)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(certRootPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientID := fmt.Sprintf(clientIDFormat, projectID, region, registryID, deviceID)

	opts := mqtt.NewClientOptions().
		AddBroker(googleBrokerEndpoint).
		SetClientID(clientID).
		SetPassword(tokenString).
		SetProtocolVersion(googleMQTTVersion).
		SetUsername(googleUsername)

	//create and start a client using the above ClientOptions
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("networkerror: %v", token.Error())
	}
	return &GoogleBroker{basicBroker{client, "/devices/" + deviceID + "/events"}}, nil
}
