FROM golang:1.19

LABEL author="Siska Stevani"
LABEL github="github.com/siskastev"

WORKDIR /stev/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

RUN go build -o main ./cmd

CMD ["./main"]
