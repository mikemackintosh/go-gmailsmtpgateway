# go-gmailsmtpgateway
---

Gmail SMTP Gateway is designed to be a Gmail API gateway driven by SMTP. In short,
you would use this system to receive SMTP messages, which when received, forward
the message to the GMail API via the authenticated user.

## Usage

To get started, you should create a new Google Project, and enable the **Gmail API**.
Once created, download your `client.json` file.

After you build the binary, you can execute it with `./bin/gmailsmtpd`. Supply the
`-o` flag, and provide the path to `client.json`.

Example:

    ./bin/gmailsmtpd -o client.json

When the script initializes, it will give you a link to authorize access to Gmail.

**NOTE** The account that you use will always be the authenticated GMail account, and will
override the `MAIL FROM` command from the SMTP server.

## Installation

    make build

## Credits
Credits to other sources used in the creation of this utility:

- [@mhale's **smtpd**](https://github.com/mhale/smtpd) package
- [@google's **gmail**](https://google.golang.org/api/gmail) package
