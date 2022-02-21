FROM golang:1.17 AS build
ARG SHORT_SHA

WORKDIR /app/fakestock
COPY . .
RUN CGO_ENABLED=0 go build -o /app/fakestock/fakestock -ldflags="-X 'github.com/cloud66-oss/fakestock/main.Commit=${SHORT_SHA}'"

FROM alpine
LABEL maintainer="Cloud 66 Engineering <hello@cloud66.com>"

RUN mkdir /app
# Copy any other files required in the final image here
COPY tickers.csv /app/tickers.csv
COPY --from=build /app/fakestock/fakestock /app
EXPOSE 8080
ENTRYPOINT ["/app/fakestock"]