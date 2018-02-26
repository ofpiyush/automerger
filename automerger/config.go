package automerger

import (
	"flag"
	"io/ioutil"
	"log"
)

type Assignees []string

func (a *Assignees) String() string {
	return "my string representation"
}

func (a *Assignees) Set(value string) error {
	*a = append(*a, value)
	return nil
}

type Config struct {
	IntegrationID string
	Secret        []byte
	Key           []byte
	Address       string
	Ref           string
	Branch        string
	Assignees     Assignees
	Token         string
	GetToken      func(int) (string, []error)
	ApiURL        string
	MergeMethod   string
}

func ConfigureOrDie() *Config {
	config := &Config{}
	var privateKeyFile, token, secret string
	var err error

	flag.StringVar(&config.IntegrationID, "integration-id", "9441", "Github Integration's id.")
	flag.StringVar(&secret, "secret", "", "Github Integration's secret.")
	flag.StringVar(&privateKeyFile, "key-file", "private_key.pem", "Full path to the Github Integration's private key.\nEither this or token should be present.")
	flag.StringVar(&token, "token", "", "Usable token with access to repo.\nEither this or the key file should be present.")
	flag.StringVar(&config.Address, "address", "0.0.0.0:3000", "Address to listen on.")
	flag.StringVar(&config.Branch, "branch", "master", "Branch to start PRs from.")
	flag.StringVar(&config.ApiURL, "api-url", "https://api.github.com", "URL of github installation.")
	flag.StringVar(&config.MergeMethod, "merge-method", "squash", "How to merge changes on default branch to current branch.")
	flag.Var(&config.Assignees, "assignee", "Assignees if PR merge fails\nDefaults to the pusher.")

	flag.Parse()

	if secret == "" {
		log.Fatalln("Need secret to verify github payloads.\nSee: https://developer.github.com/webhooks/securing/")
	}
	config.Secret = []byte(secret)

	if token == "" {
		config.Key, err = ioutil.ReadFile(privateKeyFile)
		fatalErr(err)
		config.GetToken = (&Token{Config: config}).GetorDie
	} else {
		config.GetToken = func(int) (string, []error) { return token, nil }
	}
	// Keep a cache for faster checks
	config.Ref = "refs/heads/" + config.Branch
	return config
}
