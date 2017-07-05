# csv-mailer
Simple utility to send emails from a CSV file by using the https://mailgun.com API.
To use this utility, you must have a Mailgun account.

## Installation from source
To install from source, you must have [Go](https://golang.org) installed.

```sh
$ go get -u github.com/myrenett/csv-mailer
```

You can optionally use [dep](https://github.com/golang/dep) to install dependencies.


Binary distributions are not provided at this point in time.

## Basic usage

```sh
$ export MG_API_KEY=my-api-key
$ export MG_DOMAIN=my-domain
$ echo 'email,name,phone' > mail.csv
$ echo 'email1@example.com,"Norman, Ola",123' >> mail.csv
$ echo 'email2@example.com,John Doe,124' >> mail.csv
% echo 'Hello %recipient.name%.' > mail.tmpl
% echo 'You have email %recipient.email% and phone number %recipient.phone% registered with us.' >> mail.tmpl
$ csv-mailer -subject "test email"  # perfoms a test only, check your log at https://mailgun.com
$ csv-mailer -subject "test email" -send # Sends out the email(s)!


```
TODO:
- Support more CSV seperators such as colon and tab
- Support tagging of emails for better compaign management


# License
Copyright 2017 Sindre Myren

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this program or any of it's source code except in
compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
