package iguana

import (
	"fmt"
	"log"
	"os"
	"strings"
	"regexp"
)

//	TODO ! GENERAL
/*
	ReadFile(*.html)

	parseHTML :
		if IMPORT	: importPackage && importTemplate
		if REQUIRE	: importTemplate && callTemplate
		if CALL		: parseHTML && putTemplate
*/

type Template struct {
	file		string
	name		string
	isAccessed	bool
	content		string
}


const (
	//	Regular expression for searching for <import>
	TAG_IMPORT_REGEX  string = `\s*<import src=["|']\w*\.http["|'] ?\/>`
	//	Regular expression for searching for <require>
	TAG_REQUIRE_REGEX string = `\s*<require template=["|']\w*\.htt["|'] ?\/>`
	//	Regular expression for searching for <call>
	TAG_CALL_REGEX    string = `\s*<call template=["|']\w+\.?\w*["|'] ?\/>`

	//	Regular expression to preserve indentations of the tags
	// LINE_IDENT_REGEX string = `^( +|\t+|)`
	LINE_IDENT_REGEX string = `\s*<`

	//	Regular expression to find file paths
	PATH_REGEX string = `\w+\.\w+`
)


var (
	//	Compilation of all regular expressions for new tags
	TAG_COMPILER_REGEX string = fmt.Sprintf(
		"%s|%s|%s",
		TAG_IMPORT_REGEX,
		TAG_REQUIRE_REGEX,
		TAG_CALL_REGEX,
	)

	RELATIVE_PATH string

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

	INDENTATION_SYMBOL string = "\t"
)

// TODO #1 Verify tag syntax

//	Error handler function
func check(_err error) {
	if _err != nil {
		panic(_err)
	}
}


func GenerateHTML(_relFilePath string, _newFileName string) {
	log.Printf("! [require] package is created")
	IMPORT["require"] = map[string]Template{}

	//	Separates relative path to files and a file name
	pathSeparator := strings.LastIndex(_relFilePath, "/")

	RELATIVE_PATH = _relFilePath[:pathSeparator]
	fileName := _relFilePath[pathSeparator+1:]
	// log.Println(RELATIVE_PATH, fileName)

	log.Printf("! parseHTML '%s' -> begin", fileName)
	var content string = parseHTML(fileName)
	log.Printf("! parseHTML '%s' <- end", fileName)

	// log.Printf("%s", content)

	newFilePath := "./" + _newFileName

	log.Printf("! Creating file '%s'", newFilePath)
	file, err := os.Create(newFilePath)
	check(err)

	log.Printf("! Writing file contents")
	n, err := file.WriteString(content)
	check(err)

	log.Printf("! Saved %d bytes", n)

	log.Printf("! Syncing file")
	file.Sync()

	log.Printf("! Closing file '%s'", newFilePath)
	file.Close()

	log.Printf("! HTML Generated SUCCESSFULLY")

	// unpackDict()
}


//	Reads the HTML file, finds all new tags and work them around
func parseHTML(_fileName string) string {
	data, err := os.ReadFile(RELATIVE_PATH + "/" + _fileName)
	check(err)
	// log.Printf("\n%s", data)

	// var fileLines []string = strings.Split(string(data), "\n")

	//	Create regex Compiler
	regexTagCompiler, err := regexp.Compile(TAG_COMPILER_REGEX)
	check(err)

	//	Search for all indexes that matches our new tags
	regexTagCompilerFinds := regexTagCompiler.FindAllIndex(data, -1)
	// log.Printf("%v", regexTagCompilerFinds)

	regexIndentationCompiler, err := regexp.Compile(LINE_IDENT_REGEX)
	check(err)
	// regexIndentationCompilerFinds := regexIndentationCompiler.FindAllIndex(data, -1)

	// log.Println(regexFinds)

	//	Create new file
	var newFile string = ""
	var currentIndex uint = 0
	//	For every regex find: do
	for _, tagIndexes := range regexTagCompilerFinds {
		//	Everything before tag
		newFile += string(data[currentIndex:tagIndexes[0]])
		currentIndex = uint(tagIndexes[1]) + 1

		//	Get tag
		tag := string(data[tagIndexes[0]:tagIndexes[1]])

		// Compare tag
		if strings.Contains(tag, "<import") {
			//	Import all templates from package
			regexPackageCompiler, err := regexp.Compile(PATH_REGEX)
			check(err)

			packageName := string(regexPackageCompiler.Find([]byte(tag)))
			log.Printf("! importPackage [%s] -> begin", packageName)
			importPackage(packageName)
			log.Printf("! importPackage [%s] <- end", packageName)
		}

		if strings.Contains(tag, "<require") {
			//	Put content of template
			regexRequireCompiler, err := regexp.Compile(PATH_REGEX)
			check(err)

			requireFileName := string(regexRequireCompiler.Find([]byte(tag)))
			log.Printf("! importTemplate [%s] %s '%s' -> begin", "require", requireFileName, "")
			importTemplate("require", requireFileName, "")
			log.Printf("! importTemplate [%s] %s '%s' <- end", "require", requireFileName, "")
		}

		if strings.Contains(tag, "<call") {
			//	Use one of the templates from package
			regexCallCompiler, err := regexp.Compile(PATH_REGEX)
			check(err)

			templateName := string(regexCallCompiler.Find([]byte(tag)))

			regexIndentationCompilerFinds := regexIndentationCompiler.Find([]byte(tag))

			log.Printf("! callTemplate '%s' -> begin", templateName)
			newFile += callTemplate(templateName, len(regexIndentationCompilerFinds)-1)
			log.Printf("! callTemplate '%s' <- end", templateName)
		}

		// print(newFile)
	}

	//	Save everything after tags
	newFile += string(data[currentIndex:])

	if newFile[len(newFile)-1] == '\n' {
		return newFile[:len(newFile)-1]
	}
	return newFile
}


func importPackage(_packageName string) {
	log.Printf("□ %s", _packageName)

	data, err := os.ReadFile(RELATIVE_PATH+"/"+_packageName)
	check(err)

	//	Separate lines
	fileLines := strings.Split(string(data), "\n")

	//	Extract and save package name from 1st line
	packageName := fileLines[0]
	packageName = packageName[1:len(packageName)-1]

	log.Printf("□ Imported as [%s]", packageName)

	IMPORT[packageName] = map[string]Template{}

	//	For every line in package file : do
	for _, line := range fileLines[1:] {
		//	Ignore spaces
		if line == "" {
			continue
		}
		lineBits := strings.Split(line, " ")

		templateMode := lineBits[0]
		templateFile := lineBits[1]
		templateName := lineBits[2][1:len(lineBits[2])-1]


		//	Get the content of template
		if templateMode == "~" || templateMode == "+" {
			log.Printf("! importTemplate [%s] %s '%s' -> begin", packageName, templateFile, templateName)
			importTemplate(packageName, templateFile, templateName)
			log.Printf("! importTemplate [%s] %s '%s' <- end", packageName, templateFile, templateName)
		} else {
			//	Save the template to dictionary
			IMPORT[packageName][templateName] = Template{
				templateFile,
				templateName,
				false,
				string([]byte{}),
			}
		}
	}

	// unpackDict()

	log.Printf("□ [%s] is imported", packageName)
}


func importTemplate(_packageName string, _templateFile string, _templateName string) {
	log.Printf("\t# %s", _templateFile)

	templateContent, err := os.ReadFile(RELATIVE_PATH + "/" + _templateFile)
	check(err)

	if len(_templateName) == 0 {
		_templateName = strings.Split(_templateFile, ".")[0]
	}

	log.Printf("\t# Imported as '%s' template", _templateName)

	//	Save the template to dictionary
	IMPORT[_packageName][_templateName] = Template{
		_templateFile,
		_templateName,
		true,
		string(templateContent),
	}

	log.Printf("\t# '%s' template is imported as part of [%s] package", _templateName, _packageName)

	// unpackDict()
}


func putTemplate() {}


func callTemplate(_templateName string, _indentSize int) string {
	templatePath := strings.Split(_templateName, ".")
	template := IMPORT[templatePath[0]][templatePath[1]]

	log.Printf("! parseHTML '%s' -> begin", template.file)
	parsedTemplate := parseHTML(template.file) // ! No point in saving content if parsing files
	log.Printf("! parseHTML '%s' <- end", template.file)
	

	templateForInsertion := "\n"
	lines := strings.Split(parsedTemplate, "\n")
	for _, line := range lines {
		if (_indentSize - 1) > 0 {
			templateForInsertion += strings.Repeat(INDENTATION_SYMBOL, _indentSize - 1)
		}
		templateForInsertion += line + "\n"
	}

	return templateForInsertion
}
