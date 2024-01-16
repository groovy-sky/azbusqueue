#!/bin/bash
rgName="sb-rg"
location="westeurope"
sbQueueName="demo"

echo "[INF] Creating resource group $rgName[location: $location]"
sbName="servicebus$(az group create --name $rgName --location $location -o tsv --query id | md5sum | cut -d ' ' -f 1)"

echo "[INF] Creating Service Bus $sbName in $rgName[location: $location]"
az servicebus namespace create --resource-group $rgName --name $sbName --location $location

echo "[INF] Creating queue $sbQueueName"
az servicebus queue create --resource-group $rgName --namespace-name $sbName --name $sbQueueName

echo "[INF] Granting access and storing it in AZURE_SERVICEBUS_CONNECTION_STRING"
export AZURE_SERVICEBUS_CONNECTION_STRING=$(az servicebus namespace authorization-rule keys list --resource-group $rgName --namespace-name $sbName --name RootManageSharedAccessKey --query primaryConnectionString --output tsv)

echo "[INF] Installing Go package"
go install github.com/groovy-sky/azbusqueue/v2@latest
export PATH="$HOME/go/bin:$PATH"
alias sb="azbusqueue"

echo "[INF] Testing message sending"
sb send --queue $sbQueueName --message "Azure Service Bus $sbName has been successfully created and configured. This message has been sent through $sbQueueName queue."

echo "[INF] Testing message recieving"
sb recieve --queue $sbQueueName