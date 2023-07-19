FROM golang:1.20-alpine AS build

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin cmd/guarde.go

FROM golang:1.20-alpine

COPY --from=build /usr/local/bin /usr/local/bin

CMD ["guarde"]