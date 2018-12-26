FROM centos

LABEL version=1.0
LABEL name="go-ucast dev"
RUN yum -y install git net-tools tcpdump nmap-ncat; \
    yum -y groupinstall 'Development Tools' ; \
    yum clean all ; \
    mkdir /repo && cd /repo ; \
    git clone https://github.com/elisescu/udpcast.git ; \
    cd udpcast && ./configure && make

WORKDIR /repo/udpcast
CMD [ "/bin/bash/" ]