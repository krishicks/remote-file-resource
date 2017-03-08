FROM alpine

RUN apk --update upgrade && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD assets/ /opt/resource/

RUN chmod +x /opt/resource/*
