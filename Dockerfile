FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o rewardpredictor

EXPOSE 8092

CMD ["./rewardpredictor"]