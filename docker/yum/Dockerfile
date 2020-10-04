FROM centos:7

WORKDIR /ahab-docker

RUN yum -y install epel-release python3-pip

COPY ahab .

# Spit out these just for easier debugging
RUN yum list installed

# Deprecated
RUN yum list installed | ./ahab chase --os fedora
# New way
RUN yum list installed | ./ahab chase --package-manager yum
