FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY model/ /app/model
COPY repository/ /app/repository
COPY server/ /app/server
COPY database/ /app/database
COPY main.go .

RUN go build -o userfy-server

RUN cd database
CMD [ "sqlite3", "users.db", ".read users.sql" ]

RUN cd ..

EXPOSE 8080

ENTRYPOINT ["./userfy-server"]