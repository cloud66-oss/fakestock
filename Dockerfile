FROM golang:1.17 AS build
ARG SHORT_SHA

WORKDIR /app/fakestock
COPY . .
RUN CGO_ENABLED=0 go build -o /app/fakestock/fakestock -ldflags="-X 'main.Commit=${SHORT_SHA}'"

FROM alpine
LABEL maintainer="Cloud 66 Engineering <hello@cloud66.com>"

RUN mkdir /app

COPY nasdaq.csv /app/nasdaq.csv
COPY nyse.csv /app/nyse.csv

COPY --from=build /app/fakestock/fakestock /app
ENV FAKESTOCK_PATH=/app
EXPOSE 8080
ENTRYPOINT ["/app/fakestock"]