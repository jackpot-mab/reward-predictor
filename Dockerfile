FROM golang:1.21

WORKDIR /app

COPY . .

RUN go build -o rewardpredictor

EXPOSE 8092

CMD ["./rewardpredictor"]