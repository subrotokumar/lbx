FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o lbx cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/lbx .
RUN chmod +x lbx
CMD [ "/app/lbx" ]