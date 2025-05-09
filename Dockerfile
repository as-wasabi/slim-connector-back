FROM golang:1.22.6-bullseye AS builder
LABEL authors="kigawa"

WORKDIR /app

COPY ./ /app

RUN go build slim-connector-back/cmd/tasq && chmod +x /app/tasq
RUN ls ./
#FROM gcr.io/distroless/static-debian12
FROM debian

WORKDIR /app
COPY --from=builder /app/tasq /app/

#CMD ["ls"]
CMD ["/app/tasq"]
