FROM golang:1.24.4 AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .


# Stage 2 : create a distroless image


FROM gcr.io/distroless/base

COPY --from=base /app/main main

COPY --from=base /app/.env .env

EXPOSE 3000

CMD [ "./main" ]
