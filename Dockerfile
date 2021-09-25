FROM golang:1.16.5

WORKDIR $GOPATH/src/user-access-management

COPY . .

ENV POSTGRESQL_CONN_STRING $POSTGRESQL_CONN_STRING
ENV TOKEN_SERVICE_URL $TOKEN_SERVICE_URL
ENV JWT_TOKEN_SECRET $JWT_TOKEN_SECRET

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# REST Server PORT
EXPOSE 3540

# Run the executable
CMD go run main.go
