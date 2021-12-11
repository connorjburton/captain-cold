FROM golang:1.16-alpine as build
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build

FROM alpine
COPY --from=build /app/captain-cold . 
CMD ./captain-cold