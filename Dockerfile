FROM golang:1.13.0 AS builder
 
# Add all the source code (except what's ignored
# under `.dockerignore`) to the build context.
ADD ./ /go/src/kafkaAPI
ADD ./kafkaUtils /go/src/kafkaAPI/kafkaUtils
 
#RUN  apt-get install bash
RUN go get -u github.com/kardianos/govendor
 
RUN set -ex && \
 cd /go/src/kafkaAPI && \
  CGO_ENABLED=0 govendor init && \
  CGO_ENABLED=0 govendor fetch +missing && \
  CGO_ENABLED=0 go build \
        -tags netgo \
        -v -a \
        -ldflags '-extldflags "-static"' && \
  mv ./kafkaAPI /usr/bin/kafkaAPI
 
FROM alpine:3.10.2
 
# Retrieve the binary from the previous stage
COPY --from=builder /usr/bin/kafkaAPI /usr/local/bin/kafkaAPI
 
# Set the binary as the entrypoint of the container
EXPOSE 9000/tcp
ENTRYPOINT [ "kafkaAPI" ]
