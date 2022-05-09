FROM alpine:latest

ADD email_connector_linux ./email
RUN mkdir /root/tmp_images

ADD .env ./env
EXPOSE 9200
ENTRYPOINT ["./email"]
