FROM alpine:latest

ENV SLACK_TOKEN = ""

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY ./kitty-slack .
RUN chmod +x kitty-slack

CMD ["./kitty-slack"]