FROM golang:1.17-alpine as build
RUN mkdir /src
WORKDIR /src
COPY . .
RUN go build -o /dist

FROM alpine
COPY --from=build /dist /dist
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
RUN chmod 755 /usr/bin/aws-lambda-rie
COPY entry.sh /
RUN chmod 755 /entry.sh
ENTRYPOINT ["/entry.sh"]
CMD ["/dist"]