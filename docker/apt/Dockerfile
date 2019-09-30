FROM ubuntu:latest

RUN apt-get update && apt-get -y install python3-pip

COPY ahab .

RUN apt list --installed
RUN apt list --installed | ./ahab chase