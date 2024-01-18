# Azure Service Bus Queue CLI

![](logo.svg)

## Introduction
This Go package provides a simple command-line interface (CLI) tool for interacting with [Azure Service Bus Queue](https://azure.microsoft.com/en-us/products/service-bus). It allows you to send/receive messages from/to a queue. 

## Getting started
### Prerequisites
- Go, version 1.18 or higher - [Install Go](https://go.dev/doc/install)
- Azure subscription - [Create a free account](https://azure.microsoft.com/free/)
- Service Bus namespace - [Create a namespace](https://learn.microsoft.com/azure/service-bus-messaging/service-bus-create-namespace-portal)
- Service Bus queue - [Create a queue](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quickstart-portal#create-a-queue)

## Installation

If you want to try out the tool, you'll need to:
1. Deploy Azure Service Bus and create a queue.
2. Install the tool.

Easiest way how-to do that is to use the [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest) and execute [script.sh](script.sh) file:

```
wget https://raw.githubusercontent.com/groovy-sky/service-bus-queue-cli/main/script.sh
chmod +x script.sh
./script.sh
```

## Examples

Once you have created a queue, you can use the tool to send/receive messages from/to the queue.
At first you'll need to set the AZURE_SERVICEBUS_CONNECTION_STRING environment variable to your Azure Service Bus connection string. Here's how you can do it:

```
export AZURE_SERVICEBUS_CONNECTION_STRING="your-connection-string"
```

After that you can use the tool to send/receive messages from/to a queue. 

### Send message

To send a message(subject and reply-to are optional):

``` 
azbusqueue send --queue your-queue-name --message "Your message" --subject "Your subject" --replyto "Your reply-to"
```
### Receive message

To receive a message from a queue:

```
azbusqueue read --queue your-queue-name
```

