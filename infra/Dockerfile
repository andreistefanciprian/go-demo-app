## Build
FROM golang:1.19-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

ENV GOOS=linux \
GOARCH=386

RUN go build -a -o go-demo-app

## Deploy
FROM gcr.io/distroless/base-debian11

COPY --from=build /app/go-demo-app .

USER nonroot:nonroot

ENTRYPOINT ["./go-demo-app"]