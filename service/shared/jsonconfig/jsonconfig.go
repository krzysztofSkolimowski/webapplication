package jsonconfig

import (
	"io"
	"os"
	"path/filepath"
	"log"
	"io/ioutil"
)

type Parser interface {
	ParseJSON([]byte) error
}

func Load(config string, p Parser) {
	var err error
	var path string
	var rc = io.ReadCloser(os.Stdin)
	if path, err = filepath.Abs(config); err != nil {
		log.Fatalln(err)
	}

	if rc, err = os.Open(path); err != nil {
		log.Fatalln(err)
	}

	// Read the config file
	jsonBytes, err := ioutil.ReadAll(rc)
	rc.Close()
	if err != nil {
		log.Fatalln(err)
	}

	// Parse the config
	if err := p.ParseJSON(jsonBytes); err != nil {
		log.Fatalln("Could not parse %q: %v", config, err)
	}
}
