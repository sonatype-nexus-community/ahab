FROM ubuntu:latest

WORKDIR /ahab-docker

RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 871920D1991BC93C
RUN apt-get update
RUN chmod 666 /etc/apt/trusted.gpg.d/ubuntu-keyring-2012-cdimage.gpg
RUN chmod 666 /etc/apt/trusted.gpg.d/ubuntu-keyring-2018-archive.gpg
#RUN apt-get update & apt-get install -y gnupg
#RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 871920D1991BC93C
RUN apt-get update && apt-get -y install python3-pip

COPY ahab .

# Spit out these just for easier debugging
RUN dpkg-query --show --showformat='${Package} ${Version}\n'

RUN dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase

