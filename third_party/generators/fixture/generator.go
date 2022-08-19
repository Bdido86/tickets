package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var fileName = flag.String("filename", "", "Set file name to generate code for")

const (
	packageTpl = `package %s
`

	instanceTpl = `
type %[1]sBuilder struct {
	instance *%[1]s
}

func %[1]sFixture() *%[1]sBuilder {
	return &%[1]sBuilder{
		instance: &%[1]s{},
	}
}
`
	fieldTpl = `
func (b *%[1]sBuilder) %[2]s(v %[3]s) *%[1]sBuilder {
	b.instance.%[2]s = v
	return b
}
`
	pointerValueTpl = `
func (b *%[1]sBuilder) P() *%[1]s {
	return b.instance
}

func (b *%[1]sBuilder) V() %[1]s {
	return *b.instance
}
`
)

type fileStruct struct {
	filePackage string
	filePath    string
	files       []fileInfo
}

type fileInfo struct {
	name   string
	fields []fileField
}

type fileField struct {
	name  string
	value string
}

func main() {
	flag.Parse()

	if len(*fileName) == 0 {
		log.Fatalln("Empty flag filename. Example --filename ./file.go")
	}

	if _, err := os.Stat(*fileName); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("File not exist: %s", *fileName)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, *fileName, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	var fileStructVar fileStruct
	var fileInfoVar fileInfo
	var fileFieldVar fileField
	var fileFields []fileField

	filePath := filepath.Dir(*fileName)
	if filePath != "." {
		fileStructVar.filePath = filePath
	}

	fileStructVar.filePackage = node.Name.Name

	for _, f := range node.Decls {
		genD, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genD.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			currStruct, ok := currType.Type.(*ast.StructType)
			if !ok {
				continue
			}

			fileInfoVar = fileInfo{}
			fileInfoVar.name = currType.Name.Name

			fileFields = []fileField{}
			for _, field := range currStruct.Fields.List {
				if len(field.Names) == 0 {
					continue
				}
				fileFieldVar = fileField{}
				fileFieldVar.name = field.Names[0].Name
				fileFieldVar.value = field.Type.(*ast.Ident).Name

				fileFields = append(fileFields, fileFieldVar)
			}

			fileInfoVar.fields = fileFields
			fileStructVar.files = append(fileStructVar.files, fileInfoVar)
		}
	}

	if len(fileStructVar.files) == 0 {
		log.Fatalf("File %s not found any struct", *fileName)
	}

	for _, file := range fileStructVar.files {
		file.generateFile(fileStructVar)
	}

	fmt.Println("successful fixture generation")
}

func (file *fileInfo) generateFile(fileStructVar fileStruct) {
	fNameOut := strings.ToLower(file.name) + "_fixture.go"
	if fileStructVar.filePath != "" {
		fNameOut = fileStructVar.filePath + "/" + fNameOut
	}

	out, _ := os.Create(fNameOut)
	defer out.Close()

	file.writePackage(fileStructVar.filePackage, out)
	file.writeInstance(out)
	file.writeFields(out)
	file.writePointerValue(out)
}

func (file *fileInfo) writePackage(filePackage string, out *os.File) {
	fmt.Fprint(out, fmt.Sprintf(packageTpl, filePackage))
}

func (file *fileInfo) writeInstance(out *os.File) {
	fmt.Fprint(out, fmt.Sprintf(instanceTpl, convertToUpperFirstChar(file.name)))
}

func (file *fileInfo) writeFields(out *os.File) {
	for _, field := range file.fields {
		fmt.Fprint(out, fmt.Sprintf(fieldTpl, convertToUpperFirstChar(file.name), convertToUpperFirstChar(field.name), field.value))
	}
}

func (file *fileInfo) writePointerValue(out *os.File) {
	fmt.Fprint(out, fmt.Sprintf(pointerValueTpl, convertToUpperFirstChar(file.name)))
}

func convertToUpperFirstChar(input string) string {
	return strings.ToUpper(input[:1]) + input[1:]
}
