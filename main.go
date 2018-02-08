package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/markbates/inflect"
)

func main() {
	templateName := flag.String("plate", "", "Name of template to use")
	csvName := flag.String("csv", "", "Name of CSV file to use")
	outputName := flag.String("output", "", "Name of output file to use (default: standard out)")
	wrapOutput := flag.Bool("wrap-output", false, "Wraps output in a preview page suitable for display in a web browser")
	flag.Parse()

	if err := parseAndRun(*templateName, *csvName, *outputName, *wrapOutput); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}

const wrapperHTML = `
<html>
	<head>
		<title>Small Plate Preview</title>
		<style>
			body {
				padding: 20px;
				min-width: 100vw;
				min-height: 100vh;
				box-sizing: border-box;
			}
			textarea, iframe {
				resize: both;
				width: 50vw;
				height: 20vh;
				overflow: scroll;
			}
		</style>
	</head>
	<body>
		<label>
			<h2>Code:</h2>
			<textarea id="codebox">{{.}}</textarea>
			<script>
			document.getElementById('codebox').addEventListener('focus', function(e) {
				e.target.setSelectionRange(0, e.target.value.length);
				if (document.execCommand("copy")) {
					alert("Copied");
				}
			});
			</script>

		</label>
		<div>
			<h2>Preview:</h2>
			<iframe src="data:text/html,{{.}}">
		</div>
	</body>
</html>
`

func parseAndRun(templateName, csvName, outputName string, wrapOutput bool) error {
	var output io.Writer = os.Stdout
	if outputName != "" && outputName != "-" {
		f, err := os.Create(outputName)
		if err != nil {
			return err
		}
		defer f.Close()

		output = f
	}

	if !wrapOutput {
		return run(templateName, csvName, output)
	}

	var buf bytes.Buffer
	if err := run(templateName, csvName, &buf); err != nil {
		return err
	}

	wrapper := template.Must(template.New("wrapper").Parse(wrapperHTML))
	return wrapper.Execute(output, buf.String())
}

var funcMap = map[string]interface{}{
	"unescape": func(s string) template.HTML { return template.HTML(s) },
	"groupby":  groupBy,
	"int":      func(s string) int { i, _ := strconv.Atoi(s); return i },
}

func run(templateName, csvName string, output io.Writer) error {
	t := template.New(templateName).Funcs(funcMap).Funcs(inflect.Helpers)
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
	cr.Comment = '#'
	cr.FieldsPerRecord = -1
	cr.ReuseRecord = true

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
