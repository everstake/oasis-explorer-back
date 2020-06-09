FROM golang:latest AS build-env

RUN mkdir /app
WORKDIR /app
COPY . .
ENV CGO_ENABLED 0
RUN go build

FROM alpine:latest
RUN apk add ca-certificates
COPY --from=build-env /app/oasisTracker /
COPY --from=build-env /app/.secrets /
COPY --from=build-env app/dao/mysql/migrations /dao/mysql/migrations
COPY --from=build-env app/dao/clickhouse/migrations /dao/clickhouse/migrations

CMD ["/oasisTracker"]