// Following package is CLI tool, which you can read and send message to Azure Service Bus Queue
// For authentication you need to set environment variable AZURE_SERVICEBUS_CONNECTION_STRING
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/spf13/cobra"
)

type busCLI struct {
	client *azservicebus.Client
}

// Function for ServiceBusClient initialization
func (bus *busCLI) Init() (err error) {
	// Get connection string from environment variable
	connStr := os.Getenv("AZURE_SERVICEBUS_CONNECTION_STRING")
	switch {
	case connStr == "":
		return fmt.Errorf("AZURE_SERVICEBUS_CONNECTION_STRING environment variable not set")
	}

	// Create a Service Bus client
	bus.client, err = azservicebus.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		log.Fatal("\n[ERR] Failed to create client: ", err)
	}

	return err
}

// Read message from queue
func (bus *busCLI) readMessage(queue string) (message azservicebus.ReceivedMessage, err error) {
	// Create a receiver for the queue
	receiver, err := bus.client.NewReceiverForQueue(queue, nil)
	if err != nil {
		return message, err
	}

	defer receiver.Close(context.Background())

	// Create a context to limit how long we will try to receive messages
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	messages, err := receiver.ReceiveMessages(ctxWithTimeout, 1, nil)
	if err != nil {
		if err != context.DeadlineExceeded {
			log.Fatal(err)
		} else {
			return message, nil
		}
	}

	for _, msg := range messages {
		// Complete the message. This will delete the message from the queue.
		message = *msg
		receiver.CompleteMessage(ctxWithTimeout, msg, nil)
	}
	return message, err
}

// Send message from queue
func (bus *busCLI) sendMessage(queue, message, subject, replyTo string) (err error) {

	// Create context
	ctx := context.Background()

	// Create a sender for the queue
	sender, err := bus.client.NewSender(queue, nil)
	if err != nil {
		return err
	}

	defer sender.Close(ctx)

	// Create a message
	busMessage := &azservicebus.Message{
		Body:    []byte(message),
		ReplyTo: &replyTo,
		Subject: &subject,
	}

	// Send the message
	err = sender.SendMessage(ctx, busMessage, nil)
	if err != nil {
		return err
	}

	return err
}

// CLI struct for cobra
type CLI struct {
	Send struct {
		QueueName string `arg:"" name:"queue" help:"Queue name,omitempty"`
		Message   string `arg:"" name:"message" help:"Message to send,omitempty"`
		Subject   string `arg optional name:"subject" help:"Subject of message"`
		ReplyTo   string `arg optional name:"replyto" help:"ReplyTo of message"`
	} `cmd:"" help:"Send message to queue"`
	Recieve struct {
		QueueName string `arg:"" name:"queue" help:"Queue name"`
	} `cmd:"" help:"recieve message from queue"`
}

func main() {
	var bus busCLI
	var err error
	var message azservicebus.ReceivedMessage

	rootCmd := &cobra.Command{
		Use:   "busq",
		Short: "CLI tool for Azure Service Bus Queue",
	}

	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send message to queue",
		Run: func(cmd *cobra.Command, args []string) {
			err = bus.Init()
			if err != nil {
				log.Fatal("Failed to initialize client: ", err)
			}

			queueName, _ := cmd.Flags().GetString("queue")
			message, _ := cmd.Flags().GetString("message")
			subject, _ := cmd.Flags().GetString("subject")
			replyTo, _ := cmd.Flags().GetString("replyto")

			err = bus.sendMessage(queueName, message, subject, replyTo)
			if err != nil {
				log.Fatal("Failed to send message: ", err)
			} else {
				fmt.Println("Message sent successfully")
			}
		},
	}

	sendCmd.Flags().StringP("queue", "q", "", "Queue name")
	sendCmd.Flags().StringP("message", "m", "", "Message to send")
	sendCmd.Flags().StringP("subject", "s", "", "Subject of message")
	sendCmd.Flags().StringP("replyto", "r", "", "ReplyTo of message")
	sendCmd.MarkFlagRequired("queue")
	sendCmd.MarkFlagRequired("message")

	recieveCmd := &cobra.Command{
		Use:   "recieve",
		Short: "Recieve message from queue",
		Run: func(cmd *cobra.Command, args []string) {
			err := bus.Init()
			if err != nil {
				log.Fatal("Failed to initialize client: ", err)
			}

			queueName, _ := cmd.Flags().GetString("queue")

			outputFormat, _ := cmd.Flags().GetString("output")

			message, err = bus.readMessage(queueName)
			if err != nil {
				log.Fatal("Failed to read message: ", err)
			}

			// Print message. Default output format is tsv, which returns only message body. JSON returns Body, Subject and ReplyTo fields
			switch outputFormat {
			case "tsv":
				fmt.Println(string(message.Body))
			case "json":
				fmt.Printf("{\n\"body\": \"%s\",\n\"subject\": \"%s\",\n\"replyto\": \"%s\"\n}\n", string(message.Body), *message.Subject, *message.ReplyTo)
			default:
				fmt.Println(message)
			}
		},
	}

	recieveCmd.Flags().StringP("queue", "q", "", "Queue name")
	recieveCmd.Flags().StringP("output", "o", "tsv", "Output format")
	recieveCmd.MarkFlagRequired("queue")

	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(recieveCmd)

	if err = rootCmd.Execute(); err != nil {
		log.Fatal("Failed to execute command: ", err)
	}
}
