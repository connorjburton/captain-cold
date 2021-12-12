FROM golang:1.17-alpine as build
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build

FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /app/captain-cold.exe .
ENTRYPOINT ["./captain-cold.exe"]