FROM ubuntu:latest

WORKDIR /ahab-docker

RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 871920D1991BC93C
RUN apt-get update && apt-get -y install python3-pip

COPY ahab .

# Spit out these just for easier debugging
RUN dpkg-query --show --showformat='${Package} ${Version}\n'

# Deprecated
RUN dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase --os debian
# New way
RUN dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase --package-manager dpkg
