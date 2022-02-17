FROM quay.io/centos/centos:stream

WORKDIR /ahab-docker

RUN dnf upgrade -y && dnf check-update

#Required for autodetect to work
RUN dnf install -y which

COPY ahab .

# Spit out these just for easier debugging
RUN dnf list installed

RUN dnf list installed | ./ahab chase

