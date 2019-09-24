## KafkaAPI

### A test API that exposes some RESTful APIs amongst one that is indeed a kafka REST-api.

Build using 
1. docker build . -t kafkarestapi:latest -t kafkarestapi:0.2 -f ./DockerfileRESTapi
2. docker build . -t kafkamsvc:latest -t kafkamsvc:0.2 -f ./DockerfileMSVC
