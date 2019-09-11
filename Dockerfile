FROM golang:alpine AS builder

# Add all the source code (except what's ignored
# under `.dockerignore`) to the build context.
ADD ./ /go/src/kafkaAPI/

RUN set -ex && \
  cd /go/src/kafkaAPI && \       
  CGO_ENABLED=0 go build \
        -tags netgo \
        -v -a \
        -ldflags '-extldflags "-static"' && \
  mv ./kafkaAPI /usr/bin/kafkaAPI

FROM busybox

# Retrieve the binary from the previous stage
COPY --from=builder /usr/bin/kafkaAPI /usr/local/bin/kafkaAPI

# Set the binary as the entrypoint of the container
EXPOSE 9000/tcp
ENTRYPOINT [ "kafkaAPI" ]
