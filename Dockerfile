FROM centos:7
MAINTAINER The CentOS Project

# create app web
RUN mkdir -p /app/sumeru-web

# install python lib env
WORKDIR /app
ADD . /app/sumeru-web

ENTRYPOINT top -b
