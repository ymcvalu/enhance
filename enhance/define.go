package main

type File struct {
	Funs      []Func
	ImportPos int
	ImportKey bool
}

type Func struct {
	Name   string
	Ins    []Param
	Outs   []Param
	Lbrace int
	Rbrace int
	Annos  []Anno
	Recv   Recv
}

type Param struct {
	Names []string
	Type  string
}
type Recv struct {
	Name string
	Type string
}

type Anno struct {
	Line int
	Tags []string
}
