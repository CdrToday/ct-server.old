package main

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

/// config
func conf() Config {
	doc, _ := ioutil.ReadFile(_genPath(false))
	config := Config{}

	toml.Unmarshal(doc, &config)
	return config
}

/// ------ toml ------
type Config struct {
	Mail Mail
}

type Mail struct {
	Auth MailAuth
	Msg  MailMessage
}

type MailAuth struct {
	Ident string
	User  string
	Pass  string
	Host  string
}

type MailMessage struct {
	Addr    string
	Subject string
	From    string
}

/// ----- utils -------
/// # getPath
/// production for command-line arg
/// development for ../config.toml
func _genPath(pro bool) string {
	var config string

	if pro {
		config = os.Args[0]
	} else {
		config = "../config.toml"
	}

	dir, err := filepath.Abs(config)
	if err != nil {
		log.Fatal(err)
	}

	return dir
}
