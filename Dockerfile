FROM golang:1.22.1 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 make

FROM alpine

COPY --from=builder /app/build/server /app/server
COPY --from=builder /app/build/client /app/client

ENV ADDR=:1313
ENV ID=username

EXPOSE 1313

CMD ["/app/server"]
