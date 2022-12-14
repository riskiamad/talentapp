# build stage
FROM golang:1.19-alpine AS build

WORKDIR /go/src/app
COPY . .

RUN go mod tidy

RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/api ./main.go

# final stage
FROM alpine:latest as deploy
RUN apk update && apk add --no-cache tzdata
ENV TZ Asia/Jakarta
WORKDIR /usr/app

COPY --from=build /go/src/app/bin /go/bin
COPY --from=build /go/src/app/.env ./.env
COPY --from=build /go/src/app/templates ./templates
COPY --from=build /go/src/app/assets ./assets
COPY --from=build /go/src/app/driver/db/mysql/migrations ./migrations

EXPOSE 8000

ENTRYPOINT /go/bin/api