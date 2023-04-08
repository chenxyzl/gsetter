package examplepb

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const codeTemplate = `package {{.Package}}

import (
	"google.golang.org/protobuf/proto"
)

type DirtyCallback func(field int32)

type DirtyTracker interface {
	IsDirty(field int32) bool
	SetDirty(field int32)
	ResetDirty(field int32)
	ResetAllDirty()
	SetOnDirtyCallback(callback DirtyCallback)
}

{{range .Messages}}
type dirty{{.Name}} struct {
	{{.Name}}
	dirtyFields map[int32]bool
	onDirtyCallback DirtyCallback
}

{{range .Fields}}
func (x *{{$.Name}}) Set{{.Name}}(value {{.GoType}}) {
	if x.Get{{.Name}}() != value {
		x.{{.Name}} = value
		x.setDirty({{.Number}})
	}
}
{{end}}

func (x *{{.Name}}) IsDirty(field int32) bool {
	if dx, ok := x.extraData.(*dirty{{.Name}}); ok {
		return dx.dirtyFields[field]
	}
	return false
}

func (x *{{.Name}}) SetDirty(field int32) {
	if dx, ok := x.extraData.(*dirty{{.Name}}); ok {
		dx.dirtyFields[field] = true
	}
}

func (x *{{.Name}}) ResetDirty(field int32) {
	if dx, ok := x.extraData.(*dirty{{.Name}}); ok {
		delete(dx.dirtyFields, field)
	}
}

func (x *{{.Name}}) ResetAllDirty() {
	if dx, ok := x.extraData.(*dirty{{.Name}}); ok {
		dx.dirtyFields = make(map[int32]bool)
	}
}

func (x *{{.Name}}) SetOnDirtyCallback(callback DirtyCallback) {
	if dx, ok := x.extraData.(*dirty{{.Name}}); ok {
		dx.onDirtyCallback = callback
	}
}

func (x *{{.Name}}) setDirty(field int32) {
	if dx, ok := x.extraData.(*dirty{{.Name}}); ok {
		dx.dirtyFields[field] = true
		if dx.onDirtyCallback != nil {
			dx.onDirtyCallback(field)
		}
	}
}

{{end}}
`

func Gen() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path-to-pb-file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(data, fds); err != nil {
		panic(err)
	}

	fd := fds.GetFile()[0]

	tmpl, err := template.New("code").Parse(codeTemplate)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, fd)
	if err != nil {
		panic(err)
	}
}
