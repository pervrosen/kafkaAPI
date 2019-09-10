FROM alpine:3.10.2
COPY  . /app
RUN "go build ." /app
CMD /app/kafkaAPI
