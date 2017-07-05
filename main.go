package main

import (
	"fmt"
	"os"
	"strings"

	"io/ioutil"

	"github.com/namsral/flag.git"
	"gopkg.in/mailgun/mailgun-go.v1"
)

type config struct {
	Domain    string
	APIKey    string
	APIPubKey string

	Subject  string
	From     string
	ToHeader string
	Send     bool

	TemplateFile string
	DataFile     string
	Offset       int
	Limit        int
}

const usage = `Usage: csv-mailer [-domain <domain>] [-api-key '<API key>']
  [-public-api-key '<Public API key>'] [-from <user>[@<domain>]] [-to <index>]
  [-send] [-l <limit>] [-tmpl <template file>] [-csv <data file>]

Send multiple-emails based on CSV data. All options can be either provided at
the command-line, set in a config file, or read from the environment by using
uppercase, undescores and prefixing by 'MG_'. E.g. '-api-key=value' can also be
set by export MG_API_KEY=value. Values that are not specified tend to use sane
defaults that can be listed by '-help'.

Options:
`

func main() {
	cfg := parseInput()

	var msgBody []byte
	if f, err := os.Open(cfg.TemplateFile); err != nil {
		fmt.Println("ERROR: can't open template file:", err)
		os.Exit(1)
	} else {
		msgBody, err = ioutil.ReadAll(f)
		if errClose := f.Close(); errClose != nil {
			fmt.Println("WARNING: issue closing template file:", err)
		}
		if err != nil {
			fmt.Println("ERROR: issue reading template file:", err)
			os.Exit(0)
		}
	}

	mg := mailgun.NewMailgun(cfg.Domain, cfg.APIKey, cfg.APIPubKey)
	msg := mailgun.NewMessage(cfg.From, cfg.Subject, string(msgBody))
	if !cfg.Send {
		msg.EnableTestMode()
	}

	var data []map[string]interface{}
	if f, err := os.Open(cfg.DataFile); err != nil {
		fmt.Println("ERROR: can't open data file:", err)
		os.Exit(1)
	} else {
		data, err = readData(f, cfg.Offset, cfg.Limit)
		if errClose := f.Close(); errClose != nil {
			fmt.Println("WARNING: issue closing data file:", err)
		}
		if err != nil {
			fmt.Println("ERROR: error parsing data file:", err)
			os.Exit(1)
		}
	}

	for k := range data {
		to := data[k][cfg.ToHeader].(string)
		if !strings.Contains(to, "@") {
			fmt.Println("SKIPPED: missing email for row:", data[k])
			continue
		}
		msg.AddRecipientAndVariables(to, data[k])
	}
	fmt.Println(mg.Send(msg))

}

func parseInput() *config {
	var cfg config
	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "MG", 0)
	fs.Usage = func() {
		fmt.Print(usage)
		fs.PrintDefaults()
		os.Exit(0)
	}
	fs.StringVar(&cfg.Domain, "domain", "", "domain name to use for sending your email")
	fs.StringVar(&cfg.APIKey, "api-key", "", "your MailGun API key")
	fs.StringVar(&cfg.APIPubKey, "public-api-key", "", "your MailGun public key")

	fs.StringVar(&cfg.Subject, "subject", "", "email subject")
	fs.StringVar(&cfg.From, "from", "no-reply", "sender username or full address")
	fs.StringVar(&cfg.ToHeader, "to-header", "email", "title for column with recipient addresses")
	fs.BoolVar(&cfg.Send, "send", false, "If not supplied, the emails will be sent to mailgun in test mode")

	fs.StringVar(&cfg.TemplateFile, "tmpl", "mail.tmpl", "path to email template file")
	fs.StringVar(&cfg.DataFile, "csv", "mail.csv", "path to data file")
	fs.IntVar(&cfg.Offset, "offset", 0, "skip the first N rows")
	fs.IntVar(&cfg.Limit, "limit", 0, "if above zero, stop after N rows")

	fs.Parse(os.Args[1:])
	if !strings.Contains(cfg.From, "@") {
		cfg.From += "@" + cfg.Domain
	}

	return &cfg
}
