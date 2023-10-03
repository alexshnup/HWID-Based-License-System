FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server/server ./server/*.go 
# RUN  go build ./server/*.go 
#RUN CGO_ENABLED=0 GOOS=darwin go build -o server/bin ./server/*.go 

FROM alpine:latest
# Copy server binary from builder stage
COPY --from=builder /app/server/server /app/
# Create db directory
RUN mkdir -p /app/db
# Expose port
EXPOSE 9348
# Set entry point
ENTRYPOINT ["/app/server"]