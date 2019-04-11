FROM alpine:3.8

RUN apk upgrade --update --no-cache

USER nobody

ADD qcloud-agent_linux_amd64 /usr/local/bin/qcloud-agent

CMD ["/usr/local/bin/qcloud-agent"]