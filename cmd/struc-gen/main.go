package main

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"

	"github.com/m29h/struc-gen/cmd/struc-gen/internal/structag"

	"github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

func main() {
	sourceTypePackage := os.Getenv("GOFILE")
	fmt.Printf("file=%s\n", sourceTypePackage)

	// Inspect package and use type checker to infer imported types

	pkg := loadPackage(sourceTypePackage)

	// Lookup the given source type name in the package declarations
	parsetypes := make(map[string]*types.Struct)
	for _, sourceTypeName := range pkg.Types.Scope().Names() {

		obj := pkg.Types.Scope().Lookup(sourceTypeName)

		if obj == nil {
			failErr(fmt.Errorf("%s not found in declared types of %s",
				sourceTypeName, pkg))
		}

		// We check if it is a declared type
		if _, ok := obj.(*types.TypeName); !ok {
			continue
			//failErr(fmt.Errorf("%v is not a named type", obj))
		}
		// We expect the underlying type to be a struct
		structType, ok := obj.Type().Underlying().(*types.Struct)
		if !ok {
			continue
		}
		parsetypes[sourceTypeName] = structType

	}
	err := generate(parsetypes)
	if err != nil {
		failErr(err)
	}
}

func generate(t map[string]*types.Struct) error {

	// Get the package of the file with go:generate comment
	goPackage := os.Getenv("GOPACKAGE")
	// Start a new file in this package
	f := jen.NewFile(goPackage)

	// Add a package comment, so IDEs detect files as generated
	f.PackageComment("Code generated by generator, DO NOT EDIT.")

	// Collect Function block serializing whole struct
	for sourceTypeName, structType := range t {
		var funcBlock []jen.Code
		var funcBlockUnpack []jen.Code
		var funcBlockSize []jen.Code

		// 2. Build "m := make(map[string]interface{})"
		funcBlock = append(funcBlock, jen.Id("m").Op(":=").Lit(0))
		funcBlockUnpack = append(funcBlockUnpack, jen.Id("m").Op(":=").Lit(0))
		funcBlockSize = append(funcBlockSize, jen.Id("m").Op(":=").Lit(0))
		fmt.Printf("->%s\n", sourceTypeName)
		for i := 0; i < structType.NumFields(); i++ {

			tag := structag.NewStrucTag(structType.Field(i), structType.Tag(i))
			// Generate code for each changeset field
			p, _ := tag.Pack()
			u, _ := tag.Unpack()
			s, _ := tag.GetSize()
			funcBlock = append(funcBlock, p)
			funcBlockUnpack = append(funcBlockUnpack, u)
			funcBlockSize = append(funcBlockSize, s)

		}

		// 4. Build return statement
		funcBlock = append(funcBlock, jen.Return(jen.Id("m")))
		funcBlockUnpack = append(funcBlockUnpack, jen.Return(jen.Id("m")))
		funcBlockSize = append(funcBlockSize, jen.Return(jen.Id("m")))

		// 5. Build method
		f.Func().Params(
			structag.RootStructName().Op("*").Id(sourceTypeName),
		).Id("MarshalBinary").Params(jen.Id("b").Index().Byte()).Int().Block(
			funcBlock...,
		)
		f.Func().Params(
			structag.RootStructName().Op("*").Id(sourceTypeName),
		).Id("UnmarshalBinary").Params(jen.Id("b").Index().Byte()).Int().Block(
			funcBlockUnpack...,
		)
		f.Func().Params(
			structag.RootStructName().Op("*").Id(sourceTypeName),
		).Id("SizeOf").Params().Int().Block(
			funcBlockSize...,
		)

	}

	// Build the target file name
	goFile := os.Getenv("GOFILE")
	ext := filepath.Ext(goFile)
	baseFilename := goFile[0 : len(goFile)-len(ext)]
	targetFilename := baseFilename + "_gen.go"

	// Write generated file
	return f.Save(targetFilename)
}

func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedDeps | packages.NeedTypesInfo | packages.NeedTypes | packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		failErr(fmt.Errorf("loading packages for inspection: %v", err))
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0]
}

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
