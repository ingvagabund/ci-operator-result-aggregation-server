#!/bin/sh

function writeError() {
  curl \
    --header "Content-Type: application/json" \
    -X PUT http://localhost:8080/write \
    -d "{ \"reason\": \"$1\", \"job_name\": \"string\", \"state\": \"string\", \"type\": \"string\", \"cluster\": \"string\" }"
}

for n in {1..5}; do writeError "failed to acquire lease: resources not found"; done
for n in {1..8}; do writeError "Entrypoint received interrupt: terminated"; done
for n in {1..7}; do writeError "the build src failed after 30s with reason DockerBuildFailed: Docker build strategy has failed"; done
for n in {1..6}; do writeError "ContainerFailed one or more containers exited"; done
