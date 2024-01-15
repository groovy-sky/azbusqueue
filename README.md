# Azure Service Queue CLI

![](logo.svg)

## Introduction
This Go package provides a command-line interface (CLI) tool for interacting with Azure Service Bus. It allows you to send/receive messages from/to a queue. 

## Installation

If you want to try out the tool, you'll need to:
1. Deploy Azure Service Bus and create a queue.
2. Open Azure Cloud Shell and execute command below.

```
go install github.com/groovy-sky/service-bus-queue-cliv2@latest
export PATH="$HOME/go/bin:$PATH"
```

## Usage

At first you'll need to set the AZURE_SERVICEBUS_CONNECTION_STRING environment variable to your Azure Service Bus connection string. Here's how you can do it:

```
export AZURE_SERVICEBUS_CONNECTION_STRING="your-connection-string"
```

After that you can use the tool to send/receive messages from/to a queue. 

To send a message(subject and reply-to are optional):

``` 
service-bus-queue-cli send --queue your-queue-name --message "Your message" --subject "Your subject" --replyto "Your reply-to"
```

To receive a message from a queue:

```
service-bus-queue-cli read --queue your-queue-name
```

