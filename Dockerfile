FROM golang:1.16-alpine 

RUN mkdir /app 
WORKDIR /app 
COPY . . 

RUN go build -o sample-api 
EXPOSE 3000 3000
CMD ./sample-api