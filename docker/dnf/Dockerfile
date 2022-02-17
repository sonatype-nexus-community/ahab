FROM quay.io/centos/centos:stream

WORKDIR /ahab-docker

RUN dnf upgrade -y && dnf check-update

COPY ahab .

# Spit out these just for easier debugging
RUN dnf list installed

# Deprecated
RUN dnf list installed | ./ahab chase --os fedora
# New way
RUN dnf list installed | ./ahab chase --package-manager dnf
