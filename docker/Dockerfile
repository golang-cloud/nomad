#FROM alpine:3.6
FROM anapsix/alpine-java
#FROM hub.skyinno.com/common/jre8u121-debian:latest
LABEL maintainer="DJ Enriquez <denrie.enriquezjr@gmail.com> (@djenriquez)"

RUN addgroup nomad && \
    adduser -S -G nomad nomad

ADD mars /usr/local/bin/
ADD server.hcl /etc/nomad/

RUN chmod +x /usr/local/bin/mars

RUN mkdir -p /nomad/data && \
    mkdir -p /etc/nomad && \
    chown -R nomad:nomad /nomad

RUN ln -s  /opt/jdk/bin/java /usr/bin/java

EXPOSE 4646 4647 4648 4648/udp

ENTRYPOINT ["/usr/local/bin/mars"]