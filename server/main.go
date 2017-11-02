package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

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
