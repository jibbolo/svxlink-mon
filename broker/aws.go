package broker

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/satori/go.uuid"
)

const awsService = "iotdevicegateway"

func NewAWSBroker(region, endpoint, awsAccessKey, awsSecretKey, topic string) (*Broker, error) {
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
	req, _ := http.NewRequest("GET", "wss://"+endpoint+"/mqtt", nil)
	_, err := signer.Presign(req, nil, awsService, region, 5*time.Minute, time.Now())
	if err != nil {
		return nil, fmt.Errorf("can't presign: %v", err)
	}
	opts := mqtt.NewClientOptions().
		AddBroker(req.URL.String()).
		SetClientID(uuid.NewV4().String()).
		SetProtocolVersion(4)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("can't' connect: %v", token.Error())
	}
	return &Broker{client, topic}, nil
}
