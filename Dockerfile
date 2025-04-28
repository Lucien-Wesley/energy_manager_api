#use officiel Golang image
FROM golang:1.24.0-alpine3.21

# set working directory
WORKDIR /app

# Copy the source code
COPY . .  

#  Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o api .

# Expose the port
EXPOSE 5000

# Run the executable
CMD [ "./api" ]

#  docker compose up   energy_copilot_app  
