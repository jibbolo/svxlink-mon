package aws

import (
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jibbolo/svxlink-mon/broker"
	"github.com/satori/go.uuid"
)

const awsService = "iotdevicegateway"

type AWSBroker struct {
	client mqtt.Client
}

func New(endpoint, awsRegion, awsAccessKey, awsSecretKey string) (*AWSBroker, error) {

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.StaticProvider{
				Value: credentials.Value{
					AccessKeyID:     awsAccessKey,
					SecretAccessKey: awsSecretKey,
				},
			},
		})

	signer := v4.NewSigner(creds)
	req, _ := http.NewRequest("GET", endpoint, nil)
	_, err := signer.Presign(req, nil, awsService, awsRegion, 5*time.Minute, time.Now())
	if err != nil {
		log.Fatalf("expect no error, got %v", err)
	}
	opts := mqtt.NewClientOptions().
		AddBroker(req.URL.String()).
		SetClientID(uuid.NewV4().String()).
		SetProtocolVersion(4)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &AWSBroker{client}, nil
}

func (g *AWSBroker) Publish(topic, text string) broker.Token {
	return g.client.Publish(topic, 1, false, text)
}

func (g *AWSBroker) Close() {
	g.client.Disconnect(250)
}
