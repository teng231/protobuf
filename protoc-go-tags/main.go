package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func main() {
	var directory string
	var path string
	var filepathErr error
	flag.StringVar(&directory, "dir", ".", "directory to search for pb files")
	flag.StringVar(&path, "path", "", "pb file path")
	flag.Parse()

	if len(path) != 0 {
		filepathErr = addGoTags(path)
	} else {
		filepathErr = filepath.Walk(directory, walker)
	}

	if filepathErr != nil {
		log.Printf("could not correctly handle directory: %v\n", filepathErr)
	}
}

// walk over target directory
func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// ignore all files that have not been generated by pb
	if !strings.HasSuffix(path, ".pb.go") {
		return nil
	}

	err = addGoTags(path)
	if err != nil {
		log.Printf("%s: %v\n", path, err)
	}

	return err
}

// addGoTags will add specifically defined tags to the go protobuf structs
func addGoTags(path string) error {

	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("could not parse go file: %v", err)
	}

	var visitor visitor
	ast.Walk(visitor, astFile)

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not open output file: %v", err)
	}
	defer file.Close() // nolint: errcheck

	err = format.Node(file, fileSet, astFile)
	if err != nil {
		return fmt.Errorf("could not write output file: %v", err)
	}

	return nil
}

// define a visitor struct to handle every ast node
type visitor struct{}

// Visit will find all fields in the code and add tags that are placed
// in the comment above the field as struct tags to the field
func (v visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch node.(type) {

	// we are only interested in field types
	case *ast.Field:

		// comments after the field are stripped away by the grpc protoc-generator
		// therefore we need to check the lines above the field which are
		// contained in the comment that is part of the field
		var commentTags []string

		// iterate through the sub-nodes of the field
		ast.Inspect(node, func(child ast.Node) bool {

			// get the content of all comments in the field
			switch comment := child.(type) {
			case *ast.Comment:
				// we are only interested in comments which matches syntax
				match := commentTagExp.FindStringSubmatch(comment.Text)
				if len(match) == 4 {
					commentTags = append(commentTags, match[1])
				}
			}

			return true

		})

		// no comment tags
		if len(commentTags) == 0 {
			break
		}

		// tags are stored as childodes of type ast.BasicLit
		ast.Inspect(node, func(child ast.Node) bool {

			switch basicLit := child.(type) {
			case *ast.BasicLit:

				// unshift current tag
				curTag := strings.Trim(basicLit.Value, "`")
				commentTags = append([]string{curTag}, commentTags...)
				basicLit.Value = override(commentTags)
			}

			return true
		})

	}

	return v
}

// override tags
// the tag comes later will override the before
// example:
// `tag:"default" tag:"1" tag:"2"` will be overrided to `tag:"2"`
func override(tags []string) string {
	var tagMap map[string]string = make(map[string]string)
	var tagNames []string
	var buf []string

	// pick every one tag pair
	for _, match := range tagExp.FindAllStringSubmatch(strings.Join(tags, " "), -1) {
		if len(match) == 3 {
			// override the tag
			tagMap[match[1]] = match[2]
		}
	}

	// keep tags sorted
	for name := range tagMap {
		tagNames = append(tagNames, name)
	}
	sort.Strings(tagNames)

	for _, name := range tagNames {
		buf = append(buf, fmt.Sprintf("%s:%s", name, tagMap[name]))
	}

	return fmt.Sprintf("`%s`", strings.Join(buf, " "))
}

// tagExp will find tag expressions in the comment
var tagExp = regexp.MustCompile(`([_a-z][_\w]*):("[^"]+")`)

// commentTagExp will find tag expressions in comments which match the syntax
// "// `tagName:"tagValues"`"
// "// `tagName:"tagValues" secondTagName:"secondTagValues"`"
var commentTagExp = regexp.MustCompile("^//\\s*`(([_a-z][_\\w]*:\"[^\"]+\")(\\s+[_a-z][_\\w]*:\"[^\"]+\")*)`$")
