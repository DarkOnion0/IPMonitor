FROM docker.io/library/golang:1.17 as builder

WORKDIR /src/app

# Copy project file
COPY . .

# Donwload and Install project
RUN go get -d -v ./...
RUN env CGO_ENABLED=0 go build -o IPChecker main.go

# Create a new very lightweight image for the runtime
FROM docker.io/library/alpine:latest

WORKDIR /src/app

# Copy the executable build i nthe previous step
COPY --from=builder /src/app/IPChecker /src/app/

# Env variables
ENV DEBUG "false"

ENV MINUTES "*/10"
ENV HOURS "*"
ENV MONTH_DAY "*"
ENV MONTH "*"
ENV WEEK_DAY "*"

CMD ./IPChecker -debug $DEBUG -cron "$MINUTES $HOURS $MONTH_DAY $MONTH $WEEK_DAY"
