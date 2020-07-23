# Dockerfile References: https://docs.docker.com/engine/reference/builder/

#first stage - builder
# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Thang Le <lhthang.98@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


#second stage
FROM alpine:latest

WORKDIR /root/

#You may need to inject CA root certs into the Alpine image as well. you can do this by using apk.
#Uncomment below line
RUN apk --update add ca-certificates

##Copy files you need
COPY --from=builder /app/main .
COPY --from=builder /app/serviceAccountKey.json .

RUN mkdir -p template
COPY --from=builder /app/template template/.

#Command to run the executable
CMD ["./main"]

# Expose port 9000 to the outside world
EXPOSE 8585