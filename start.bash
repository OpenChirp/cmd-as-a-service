#!/bin/bash
# Starts the cmd-as-a service locally

export SERVICE_ID="<oc_service_id>"
export SERVICE_TOKEN="<oc_service_token>"
export FRAMEWORK_SERVER="https://api.openchirp.io"
export MQTT_SERVER="tls://mqtt.openchirp.io:8883"
export LOG_LEVEL=5

exec ./cmd-as-a-service