FROM centos:7
MAINTAINER The CentOS Project
ENTRYPOINT top -b

RUN yum install -y gcc make wget openssl-devel zlib-devel bzip2-devel libffi-devel libffi-devel cairo-devel pango-devel.x86_64

RUN wget https://www.python.org/ftp/python/3.7.1/Python-3.7.1.tar.xz

RUN tar -xvf Python-3.7.1.tar.xz && cd Python-3.7.1 && ./configure && make && make install

RUN wget https://bootstrap.pypa.io/get-pip.py && python get-pip.py

# Nmap Openssl

RUN pip install cffi

# create app web
RUN mkdir -p /app/sumeru-web

# install python lib env
WORKDIR /app
ADD . /app/sumeru-web
