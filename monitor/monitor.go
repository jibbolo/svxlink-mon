package main

// Thu Sep 28 15:32:09 2017: IW0HKS: Login OK from 195.94.189.122:63358
// Fri Oct  6 20:01:16 2017: IR0UFQ: Client 44.208.124.17:38551 disconnected: Connection closed by remote peer
import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

// Reflector is the main struct
type Reflector struct {
}

// RadioLink is the struct for radiolinks
type RadioLink struct {
	IP string
}

func main() {
	// Sets your Google Cloud Platform project ID.
	projectID := "test-168217"

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	sub := client.Subscription("sub-name")

	log.Println("receiving...")
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message: %s", m.Data)
		m.Ack()
	})
	if err != nil {
		log.Fatal(err)
	}
}
