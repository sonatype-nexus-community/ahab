FROM ubuntu:latest

WORKDIR /ahab-docker

RUN apt-get update && apt-get -y install python3-pip

COPY ahab .

# Spit out these just for easier debugging
RUN dpkg-query --show --showformat='${Package} ${Version}\n'

RUN dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase --package-manager dpkg

