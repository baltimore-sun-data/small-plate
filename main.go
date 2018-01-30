package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	templateName := flag.String("plate", "", "Name of template to use")
	csvName := flag.String("csv", "", "Name of CSV file to use")
	outputName := flag.String("output", "", "Name of output file to use (default: standard out)")
	flag.Parse()

	if err := parseAndRun(*templateName, *csvName, *outputName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}

func parseAndRun(templateName, csvName, outputName string) error {
	output := os.Stdout
	var err error
	if outputName != "" && outputName != "-" {
		output, err = os.Create(outputName)
		if err != nil {
			return err
		}
	}

	return run(templateName, csvName, output)
}

var funcMap = map[string]interface{}{
	"unescape": func(s string) template.HTML { return template.HTML(s) },
	"groupby":  groupBy,
}

func run(templateName, csvName string, output io.Writer) error {
	t := template.New(templateName).Funcs(funcMap)
	contents, err := ioutil.ReadFile(templateName)
	if err != nil {
		return err
	}
	t, err = t.Parse(string(contents))
	if err != nil {
		return err
	}

	f, err := os.Open(csvName)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := makeData(f)
	if err != nil {
		return err
	}

	return t.Execute(output, data)
}

func makeData(r io.Reader) (data []map[string]string, err error) {
	cr := csv.NewReader(r)
	fields, err := cr.Read()

	// Save headers for each row of dict
	dataHeader := make(map[int]string, len(fields))
	for i, field := range fields {
		dataHeader[i] = field
	}

	for {
		fields, err = cr.Read()
		if err == io.EOF {
			return data, nil
		}

		if err != nil {
			return nil, err
		}

		datum := make(map[string]string, len(fields))
		for i, val := range fields {
			datum[dataHeader[i]] = val
		}
		data = append(data, datum)
	}
}

type object = map[string]string
type groupedObj = struct {
	Key   string
	Items []object
}

func groupBy(key string, objs []object) []groupedObj {
	if len(objs) < 1 {
		return nil
	}

	var ret []groupedObj

	lastKey := objs[0][key]
	cur := groupedObj{Key: lastKey, Items: []object{objs[0]}}

	for _, obj := range objs[1:] {
		if obj[key] == lastKey {
			cur.Items = append(cur.Items, obj)
		} else {
			ret = append(ret, cur)
			lastKey = obj[key]
			cur = groupedObj{Key: lastKey, Items: []object{obj}}
		}
	}
	ret = append(ret, cur)

	return ret
}
