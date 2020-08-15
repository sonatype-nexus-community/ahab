FROM centos:latest

WORKDIR /ahab-docker

RUN dnf upgrade -y && dnf check-update

COPY ahab .

# Spit out these just for easier debugging
RUN dnf list installed

RUN dnf list installed | ./ahab chase --os fedora

