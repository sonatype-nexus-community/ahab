FROM ubuntu:latest

RUN apt-get update && apt-get -y install python3-pip

COPY ahab .

RUN dpkg-query --show --showformat='${Package} ${Version}\n'
RUN dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase