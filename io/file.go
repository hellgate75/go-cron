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
	"reflect"
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
var DefaultEncodingFormatString = "json"

func EncodingFromValue(value string) Encoding {
	switch strings.TrimSpace(strings.ToLower(value)) {
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

func EncodeTextFormatSummary(in interface{}) ([]byte, error) {
	var err error
	var out = make([]byte, 0)

	if strings.Contains(fmt.Sprintf("%T", in), "[]") {
		// It's an array of elements
		list, err := splitInterfaceArray(in)
		if err != nil {
			var text = fmt.Sprintf("Error: %v", err)
			out = append(out, []byte(text)...)
		} else {
			var text = ""
			if len(list) == 0 {
				text += "No results ...\n"
			}
			for idx, elem := range list {
				header, values, err := decomposeElement(elem)
				//b, _ := EncodeValue(&header, EncodingYaml)
				//fmt.Println("Index:", idx, "Header:", string(b))
				//b, _ = EncodeValue(&values, EncodingYaml)
				//fmt.Println("Index:", idx, "Values:", string(b))
				if err != nil {
					text += fmt.Sprintf("Error in line %v: %v\n", idx, err)
				} else {
					if idx == 0 {
						text += calculateHeaderLines(header)
					}
					text += calculateValueLines(values)
				}
			}
			out = append(out, []byte(text)...)
		}
	} else {
		// It's a single object
		header, values, err := decomposeElement(in)
		//b, _ := EncodeValue(&header, EncodingYaml)
		//fmt.Println("Header:", string(b))
		//b, _ = EncodeValue(&values, EncodingYaml)
		//fmt.Println("Values:", string(b))
		if err != nil {
			var text = fmt.Sprintf("Error: %v\n", err)
			out = append(out, []byte(text)...)
		} else {
			var text = calculateHeaderLines(header)
			text += calculateValueLines(values)
			out = append(out, []byte(text)...)
		}
	}
	return out, err
}

func splitInterfaceArray(in interface{}) ([]interface{}, error) {
	var err error
	var out = make([]interface{}, 0)
	var kind = reflect.TypeOf(in).Kind()
	if kind == reflect.Slice ||  kind == reflect.Array {
		v := reflect.Indirect(reflect.ValueOf(&in)).Elem()
		var length = v.Len()
		for i := 0; i < length; i++ {
			out = append(out, v.Index(i).Interface())
		}
	} else {
		err = errors.New(fmt.Sprintf("Invalid slice or list type: %s", kind.String()))
	}
	return out, err
}

func calculateHeaderLine(header headerSet) (string, []headerSet, []string, bool) {
	var hs = make([]headerSet, 0)
	var spc = make([]string, 0)
	var hasMore = false
	var out = ""
	for _, hl := range header.Columns {
		var length = len(hl.Name)
		out += hl.Name
		hs = append(hs, headerSet{hl.Columns})
		if len(hl.Columns) > 0 {
			hasMore = true
		}
		spc = append(spc, strings.Repeat(" ", length))
	}
	return out, hs, spc, hasMore
}

func calculateSubHeaderLine(headers []headerSet, spaces []string) (string, []headerSet, []string, bool) {
	var hs = make([]headerSet, 0)
	var spc = make([]string, 0)
	var hasMore = false
	var out = ""
	for idx, hsi := range headers {
		var spcLen = 0
		if len(hsi.Columns) == 0 {
			out += spaces[idx]
			spcLen = len(spaces[idx])
			hs = append(hs, headerSet{hsi.Columns})
			spc = append(spc, strings.Repeat(" ", spcLen))

		} else {
			line, hsC, spC, more := calculateHeaderLine(hsi)
			out += line
			hs = append(hs, hsC...)
			spc = append(spc, spC...)
			if more {
				hasMore = true
			}
		}
	}
	return out, hs, spc, hasMore
}

func calculateHeaderLines(header headerSet) string {
	var linesArr = make([]string, 0)
	line, hs, spc, more := calculateHeaderLine(header)
	//fmt.Println(more)
	//fmt.Println(hs)
	//fmt.Println(spc)
	//fmt.Println(line)
	linesArr = append(linesArr, line)
	for more {
		line, hs, spc, more = calculateSubHeaderLine(hs, spc)
		linesArr = append(linesArr, line)
		//fmt.Println(more)
		//fmt.Println(hs)
		//fmt.Println(spc)
		//fmt.Println(line)
	}
	return strings.Join(linesArr, "\n") + "\n"
}

func calculateValueLine(value valuesSet) (string, []valuesSet, []string, bool) {
	var hs = make([]valuesSet, 0)
	var spc = make([]string, 0)
	var hasMore = false
	var out = ""
	for _, hl := range value.Values {
		var length = len(hl.Value)
		out += hl.Value
		hs = append(hs, valuesSet{hl.SubValues})
		if len(hl.SubValues) > 0 {
			hasMore = true
		}
		spc = append(spc, strings.Repeat(" ", length))
	}
	return out, hs, spc, hasMore
}

func calculateSubValueLine(values []valuesSet, spaces []string) (string, []valuesSet, []string, bool) {
	var hs = make([]valuesSet, 0)
	var spc = make([]string, 0)
	var hasMore = false
	var out = ""
	for idx, hsi := range values {
		var spcLen = 0
		if len(hsi.Values) == 0 {
			out += spaces[idx]
			spcLen = len(spaces[idx])
			hs = append(hs, valuesSet{hsi.Values})
			spc = append(spc, strings.Repeat(" ", spcLen))

		} else {
			line, hsC, spC, more := calculateValueLine(hsi)
			out += line
			hs = append(hs, hsC...)
			spc = append(spc, spC...)
			if more {
				hasMore = true
			}
		}
	}
	return out, hs, spc, hasMore
}


func calculateValueLines(values valuesSet) string {
	var linesArr = make([]string, 0)
	line, hs, spc, more := calculateValueLine(values)
	linesArr = append(linesArr, line)
	for more {
		line, hs, spc, more = calculateSubValueLine(hs, spc)
		linesArr = append(linesArr, line)
	}
	return strings.Join(linesArr, "\n") + "\n"
}

var runeA = byte('a')
var diff = byte('a') - runeA

func formatColumn(s string) string {
	var out = ""
	for idx, c := range s {
		if idx == 0 {
			if byte(c) >= runeA {
				c = rune(byte(c) - diff)
			}
			out += fmt.Sprintf("%c", c)
		} else {
			if byte(c) < runeA {
				out += " "
			}
			out += fmt.Sprintf("%c", c)
		}
	}
	return out
}

func trimLen(s string, n int) string {
	if len(s) > n {
		return s[:n]
	} else if len(s) < n {
		return s + strings.Repeat(" ", n-len(s)-1)
	}
	return s
}


func decomposeElement(in interface{}) (headerSet, valuesSet, error) {
	var hSet = headerSet{make([]*headerElem, 0)}
	var vSet = valuesSet{make([]*valueItem, 0)}
	var err error
	if reflect.TypeOf(in).Kind() == reflect.Struct {
		var e = reflect.Indirect(reflect.ValueOf(&in)).Elem()
		var typeOfT = e.Type()
		for i := 0; i < e.NumField(); i++ {
			f := e.Field(i)
			name := typeOfT.Field(i).Name
			fName := formatColumn(name)
			var value interface{}
			if e.FieldByName(name).Type().Kind() != reflect.Interface || ! e.FieldByName(name).IsNil() {
				value = e.FieldByName(name).Interface()
			}
			fValue := ""
			if value != nil {
				fValue = fmt.Sprintf("%v", value)
			}
			vType := f.Type()
			var length = len(fValue)
			if length > 0 && len(fName) > len(fValue) {
				length = len(fName)
				fName = trimLen(fName, length + 2)
				fValue = trimLen(fValue, length + 2)
			} else if length > 0 {
				fName = trimLen(fName, length + 2)
				fValue = trimLen(fValue, length + 2)
			}
			length += 2
			headItem := headerElem{
				Value: name,
				Name: fName,
				Columns: make([]*headerElem, 0),
				Size: length,
			}
			valueItem := valueItem{
				Value: fValue,
				SubValues: make([]*valueItem, 0),
			}
			if vType.Kind() == reflect.Struct || vType.Kind() == reflect.Interface {
				var hs headerSet
				var vs valuesSet
				hs, vs, err = decomposeElement(value)
				if err == nil {
					length = 0
					for _, c := range hs.Columns {
						length += c.Size
						headItem.Columns = append(headItem.Columns, c)
					}
					headItem.Size = length
					for _, c := range vs.Values {
						valueItem.SubValues = append(valueItem.SubValues, c)
					}
					valueItem.Value = strings.Repeat(" ", headItem.Size)
				}
			}
			hSet.Columns = append(hSet.Columns, &headItem)
			vSet.Values = append(vSet.Values, &valueItem)
		}
	} else {
		err = errors.New(fmt.Sprintf("Cannot decompose non strcuture element %+v", in))
	}
	return hSet, vSet, err
}


type headerElem struct {
	Name		string
	Value		string
	Size		int
	Columns		[]*headerElem
}

type headerSet struct {
	Columns		[]*headerElem
}

type valueItem struct {
	Value		string
	SubValues	[]*valueItem
}

type valuesSet struct {
	Values		[]*valueItem
}