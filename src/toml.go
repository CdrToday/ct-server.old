package main

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func conf() *toml.Tree {
	doc, _ := ioutil.ReadFile(_genPath(false))
	config, _ := toml.Load(string(doc))

	return config
}

/// ----- utils -------
/// # getPath
/// production for command-line arg
/// development for ../config.toml
func _genPath(pro bool) string {
	var config string

	if pro {
		config = os.Args[1]
	} else {
		config = "../config.toml"
	}

	dir, err := filepath.Abs(config)
	if err != nil {
		log.Fatal(err)
	}

	return dir
}
