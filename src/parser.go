package iguana

import (
	"os"
	"strings"
	"log"
)

//	HTML template datatype
type Template struct {
	file		string
	name		string
	isAccessed	bool
	content		string
}

var (
	//	Template dictionary
	IMPORT map[string]map[string]Template = make(map[string]map[string]Template)
	/*
		IMPORT {
			"packageName_1": {
				"templateName_1": <TemplateObject>,
				"templateName_2": <TemplateObject>,
				"templateName_3": <TemplateObject>,
				"templateName_4": <TemplateObject>,
			},
			"packageName_2": {
				...
			}
		}
	*/
)

//	Read the template file, save all the templates
func parseImportFile(_fileName string) {
	log.Println("📦 ", _fileName)
	data, err := os.ReadFile(RELATIVE_PATH+"/"+_fileName)
	check(err)

	//	Separate lines
	fileLines := strings.Split(string(data), "\n")

	//	Extract and save package name from 1st line
	packageName := fileLines[0]
	packageName = packageName[1:len(packageName)-1]
	IMPORT[packageName] = map[string]Template{}

	//	For every line in package file : do
	for _, line := range fileLines[1:] {
		//	Ignore spaces
		if line == "" {
			continue
		}
		lineBits := strings.Split(line, " ")

		//	Get the content of template
		var templateContent []byte
		if lineBits[0] == "~" || lineBits[0] == "+" {
			log.Println("#️⃣ ", lineBits[1])
			templateContent, err = os.ReadFile(RELATIVE_PATH + "/" + lineBits[1])
			check(err)
		} else {
			templateContent = []byte{}
		}

		//	Get the name of the template
		var templateName string
		if len(lineBits[2]) > 2 {
			templateName = lineBits[2][1:len(lineBits[2])-1]
		} else {
			templateName = strings.Split(lineBits[1], ".")[0]
		}

		//	Save the template to dictionary
		IMPORT[packageName][templateName] = Template{
			lineBits[1],
			templateName,
			lineBits[0] == "~" || lineBits[0] == "!",
			string(templateContent),
		}
	}
}

func parseRequestFile(_fileName string) {}

func parseCall(_fileName string, _ident int) {}
