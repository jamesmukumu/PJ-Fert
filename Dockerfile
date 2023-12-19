FROM golang:1.21.5


WORKDIR /github.com/jamesmukumu/backup


COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

EXPOSE 7000

CMD ["go", "run", "main.go"]


