FROM centos:7
MAINTAINER The CentOS Project
CMD echo "Hello World" && tail -f /var/log/messages
