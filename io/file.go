package io

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Encoding string

func (e Encoding) String() string {
	return string(e)
}

const (
	EncodingUnknown = Encoding("")
	EncodingYaml    = Encoding("yaml")
	EncodingXml     = Encoding("xml")
	EncodingJson    = Encoding("json")
)

var EncodingList = "json, yaml, xml"

var DefaultEncodingFormat = EncodingJson

func EncodingFromValue(value string) Encoding {
	switch strings.ToLower(value) {
	case "yaml", "yml":
		return EncodingYaml
	case "xml":
		return EncodingXml
	case "json":
		return EncodingJson
	default:
		return EncodingUnknown
	}
}

func loadFileBytes(file string) ([]byte, error) {
	var err error
	var data = make([]byte, 0)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	f, err := os.Open(file)
	if err == nil {
		data, err = ioutil.ReadAll(f)
	}
	return data, err
}

func saveFileBytes(file string, data []byte, perm os.FileMode) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	err = ioutil.WriteFile(file, data, perm)
	return err
}

func getReaderFrom(file string) io.Reader {
	if fi, err := os.Stat(file); err == nil && !fi.IsDir() {
		f, errF := os.Open(file)
		if errF == nil {
			return f
		}
	}
	return nil
}

func getWriterFrom(file string, perm os.FileMode) io.Writer {
	if fi, err := os.Stat(file); err != nil && !fi.IsDir() {
		f, errF := os.OpenFile(file, os.O_WRONLY, perm)
		if errF == nil {
			return f
		}
	}
	return nil
}

// Load Configuration from given file
func ReadConfig(enc Encoding, file string, config interface{}) error {
	var err error
	var data = make([]byte, 0)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	switch enc {
	case EncodingJson:
		data, err = loadFileBytes(file)
		if err == nil {
			err = json.Unmarshal(data, &config)
		}
	case EncodingXml:
		data, err = loadFileBytes(file)
		if err == nil {
			err = xml.Unmarshal(data, &config)
		}
	case EncodingYaml:
		data, err = loadFileBytes(file)
		if err == nil {
			err = yaml.Unmarshal(data, &config)
		}
	default:
		err = errors.New(fmt.Sprintf("Unknown encoding format: %v", enc))
	}
	return err
}

// Load Configuration from given file
func ReadNative(file string, config interface{}) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	err = gob.NewDecoder(getReaderFrom(file)).Decode(&config)
	return err
}

// Load Configuration from given file
func SaveNative(file string, config interface{}) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	err = gob.NewEncoder(getWriterFrom(file, 0777)).Encode(&config)
	return err
}

// Save Configuration to given file
func SaveConfig(enc Encoding, file string, config interface{}) error {
	var err error
	var data = make([]byte, 0)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	switch enc {
	case EncodingJson:
		data, err = json.Marshal(&config)
		if err == nil {
			err = saveFileBytes(file, data, 0777)
		}
	case EncodingXml:
		data, err = xml.Marshal(&config)
		if err == nil {
			err = saveFileBytes(file, data, 0777)
		}
	case EncodingYaml:
		data, err = yaml.Marshal(&config)
		if err == nil {
			err = saveFileBytes(file, data, 0777)
		}
	default:
		err = errors.New(fmt.Sprintf("Unknown encoding format: %v", enc))
	}
	return err
}
