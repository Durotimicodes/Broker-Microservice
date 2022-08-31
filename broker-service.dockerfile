#1.get the base go image
#FROM golang:1.18-alpine as builder

#2 run the command on the docker image
#RUN mkdir /app 

#3. copy everything in the current folder into the docker folder
#COPY . /app

#4. Set working directory
#WORKDIR /app

#5 build go code
#RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

#6 to make sure its executable
#RUN chmod +x /app/brokerApp

#7 build a tiny docker image
FROM alpine:latest

#8 run command on a new docker image
RUN mkdir /app

#9 build from broker app and copy to /app
#COPY --from=builder /app/brokerApp /app
COPY brokerApp /app

#10 execute command
CMD [ "/app/brokerApp" ]