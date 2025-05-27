FROM golang:1.24.2
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main .

EXPOSE 8080
##|| Environment variables for Ethereum network configuration ||##
##################################################################
ENV NODE_COUNT=${NODE_COUNT}
ENV NODE_ADDRESSES=${NODE_ADDRESSES}
ENV NODE_PORTS=${NODE_PORTS}

##|| Run the WebSocket Gateway with the environment variables passed as flags ||##
##################################################################################
CMD ./main --nodes=${NODE_COUNT} --addresses=${NODE_ADDRESSES} --ports=${NODE_PORTS}