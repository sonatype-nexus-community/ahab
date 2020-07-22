FROM ubuntu:latest

WORKDIR /tmp/ahab

COPY ahab .

RUN chmod +x ahab

RUN apt-get update && apt-get install -y ca-certificates

RUN update-ca-certificates
