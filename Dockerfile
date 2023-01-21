FROM golang:1.18 as builder
WORKDIR /app/src/leaderboard-service
ENV GOPATH=/app
COPY . /app/src/leaderboard-service
RUN go get -d -v ./...
RUN go build -o main .
CMD [ "./main" ]
