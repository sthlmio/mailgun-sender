# mailgun-sender

Simple golang http server that accepts a `POST` and sends an email using mailgun.

Need to set these env vars for it to work (example):

```sh
ACCESS_CONTROL_ALLOW_ORIGIN: *
MAILGUN_API_BASE: https://api.eu.mailgun.net/v3
MAILGUN_DOMAIN:
MAILGUN_API_KEY:
MAIL_SUBJECT:
MAIL_RECIPIENT:
PORT:
```

## Multi arch

Built for `linux/amd64` and `linux/arm64/v8`.
