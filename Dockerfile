FROM golang:1.17-alpine as build
RUN mkdir /src
WORKDIR /src
COPY . .
RUN go build -o /dist

FROM alpine
COPY --from=build /dist /dist
ENTRYPOINT ["/dist"]