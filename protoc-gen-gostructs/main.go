// The protog-gen-gostructs plugin creates golang types that are easier to work
// with the the structs generated by grpc. Specifically it enables you to use
// custom tags on structs and use embed messages directly into another message
// without nesting them. The structs work in conjunction with the types generated
// by the grpc plugin. Autogenerated conversion functions allow simple conversion
// from and to the grpc structs without the use of reflections.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/dkfbasel/protobuf/protoc-gen-gostructs/plugin"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

func main() {

	gen := generator.New()

	// read the proto file passed from protoc via stdin
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		gen.Error(err, "could not read input")
	}

	err = proto.Unmarshal(data, gen.Request)
	if err != nil {
		gen.Error(err, "could not parse proto definition")
	}

	if len(gen.Request.FileToGenerate) == 0 {
		gen.Fail("no output files to generate")
	}

	// get command line paramenters
	// gen.CommandLineParameters(gen.Request.GetParameter())

	// wrap all descriptor and file descriptors
	gen.WrapTypes()

	// set the package name to be used
	gen.SetPackageNames()

	// build map of types
	gen.BuildTypeNameMap()

	// use a custom plugin to generate the output
	gen.GeneratePlugin(plugin.New())

	regxImport, err := regexp.Compile(`import (.*) "github.com\/dkfbasel\/protobuf\/types\/(.*)"`)
	if err != nil {
		log.Fatalln("could not compile regular expression")
	}

	// go through all input files and define the name of the output
	for i := 0; i < len(gen.Response.File); i++ {

		// modify the content after generation
		// NOTE: ideally this should be adapted in the template, however this
		// does currently not seem to be possible with gogo/protobuf
		var newContent = gen.Response.File[i].GetContent()

		// replace the generator name
		newContent = strings.Replace(newContent, "by protoc-gen-gogo", "by protoc-gen-gotags", -1)

		// remove proto imports
		newContent = strings.Replace(newContent, "import proto \"github.com/gogo/protobuf/proto\"\n", "", -1)

		// remove math package
		newContent = strings.Replace(newContent, "import math \"math\"\n", "", -1)

		// remove underscore definitions
		newContent = strings.Replace(newContent, "var _ = proto.Marshal\n", "", -1)
		newContent = strings.Replace(newContent, "var _ = math.Inf\n", "", -1)

		// fix imports for our types
		// note: this will only work very specifically for the dkfbasel protobuf package

		// get all package names for our types
		dkfPackages := regxImport.FindAllStringSubmatch(newContent, -1)

		for _, importMatch := range dkfPackages {

			// match the type of the package case insensitive
			typeMatch := regexp.MustCompile(fmt.Sprintf(`(?i)dkfbasel_protobuf\.%s`, importMatch[2]))

			// i.e. find dkfbasel_protobuf.Timestamp
			replaceMatch := typeMatch.FindString(newContent)

			// replace the package
			// i.e. dkfbasel_protobuf.Timestamp -> dkfbasel_protobuf1.Timestamp
			replaceWith := strings.Replace(replaceMatch, "dkfbasel_protobuf", importMatch[1], -1)

			// replache all occurances in the content
			newContent = strings.Replace(newContent, replaceMatch, replaceWith, -1)

		}

		gen.Response.File[i].Content = &newContent

		newFileName := strings.Replace(*gen.Response.File[i].Name, ".pb.go", ".alias.go", -1)

		gen.Response.File[i].Name = proto.String(newFileName)

		// gen.Error(fmt.Errorf("tmp"), gen.Response.File[i].GetName())
	}

	// return the return of the plugin to protoc via stdout
	data, err = proto.Marshal(gen.Response)
	if err != nil {
		gen.Error(err, "failed to marshal output proto")
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		gen.Error(err, "failed to pass output to protoc")
	}

}
