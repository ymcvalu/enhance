package main

import (
	"strings"
	"text/template"
)

func JoinParams(params []Param) string {
	pstr := make([]string, len(params))
	for i := range params {
		pstr[i] = strings.Join(params[i].Names, ", ")
		if pstr[i] != "" {
			pstr[i] += " " + params[i].Type
		} else {
			pstr[i] = params[i].Type
		}
	}
	return strings.Join(pstr, ", ")
}

func ParamList(params []Param) string {
	pstr := make([]string, len(params))
	for i := range params {
		pstr[i] = strings.Join(params[i].Names, ", ")
	}
	return strings.Join(pstr, ", ")
}

func JoinAnnos(annos []Anno) string {
	astr := make([]string, len(annos))
	for i := range annos {
		astr[i] = strings.Join(annos[i].Tags, ", ")
	}
	return strings.Join(astr, ", ")
}

var temp *template.Template

const tmpl = `
func {{if .Recv.Type}}({{.Recv.Name}} {{.Recv.Type}}) {{end}}{{.Name}}({{.Ins | joinParams}}){{if .Outs | len}}({{.Outs | joinParams}}){{end}}{
    {{if .Outs |len}}return {{end}}enhancer.Enhance2({{if .Recv.Name}}{{.Recv.Name}}.{{end}}_{{.Name}},{{.Annos | joinAnnos}}).(func ({{.Ins | joinParams}}){{if .Outs | len}}({{.Outs | joinParams}}){{end}})({{.Ins | paramList}})
}
`

func init() {
	temp = template.New("tmpl")
	temp.Funcs(map[string]interface{}{
		"joinParams": JoinParams,
		"joinAnnos":  JoinAnnos,
		"paramList":  ParamList,
	})
	_, err := temp.Parse(tmpl)

	if err != nil {
		panic(err)
	}
}
