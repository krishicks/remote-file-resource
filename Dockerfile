FROM busybox

ADD assets/ /opt/resource/

RUN chmod +x /opt/resource/*
