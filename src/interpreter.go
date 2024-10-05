package iguana

import (
	"fmt"
	"log"
	"os"
	"strings"

	// "strings"
	"regexp"
)

const (
	//	Regular expression for searching for <import>
	TAG_IMPORT_REGEX  string = `<import src=("|')\w*\.http("|') ?\/>`
	//	Regular expression for searching for <require>
	TAG_REQUIRE_REGEX string = `<require template=("|')\w*\.htt("|') ?\/>`
	//	Regular expression for searching for <call>
	TAG_CALL_REGEX    string = `<call template=("|')\w*\.\w*("|') ?\/>`

	//	Regular expression to preserve indentations of the tags
	LINE_IDENT_REGEX string = `^( +|\t+|)`

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
)

// TODO #1 Verify tag syntax

//	Error handler function
func check(_err error) {
	if _err != nil {
		panic(_err)
	}

}

//	Reads the HTML file, finds all new tags and work them around
func ReadFile(_relFilePath string) {
	//	Separates relative path to files and a file name
	pathSeparator := strings.LastIndex(_relFilePath, "/")
	RELATIVE_PATH	= _relFilePath[:pathSeparator]
	fileName		:= _relFilePath[pathSeparator+1:]
	log.Println(RELATIVE_PATH, fileName)
	
	//	Read HTML file
	data, err := os.ReadFile(_relFilePath)
	check(err)

	// var fileLines []string = strings.Split(string(data), "\n")

	//	Create regex Compiler
	regexTagCompiler, err := regexp.Compile(TAG_COMPILER_REGEX)
	//	Search for all indexes that matches our new tags
	regexTagCompilerFinds := regexTagCompiler.FindAllIndex(data, -1)

	regexIdentationCompiler, err := regexp.Compile(LINE_IDENT_REGEX)
	// regexIdentationCompilerFinds := regexIdentationCompiler.FindAllIndex(data, -1)

	// log.Println(regexFinds)

	//	Create new file
	var newFile string = "\n"
	var currentIndex uint = 0
	//	For every regex find: do
	for _, tagIndexes := range regexTagCompilerFinds {
		//	Everything before tag
		newFile += string(data[currentIndex:tagIndexes[0]])
		currentIndex = uint(tagIndexes[1]) + 1

		//	Get tag
		tag := string(data[tagIndexes[0]:tagIndexes[1]])

		// Compare tag
		if strings.Contains(tag, "import") {
			//	Import all templates from package
			regexImportFileCompiler, err := regexp.Compile(PATH_REGEX)
			check(err)
			importFileName := string(regexImportFileCompiler.Find(data))
			parseImportFile(importFileName)

			continue
		}

		if strings.Contains(tag, "request") {
			//	Put content of template
			regexRequestFileCompiler, err := regexp.Compile(PATH_REGEX)
			check(err)
			requestFileName := string(regexRequestFileCompiler.Find(data))
			parseRequestFile(requestFileName)

			continue
		}

		if strings.Contains(tag, "call") {
			//	Use one of the templates from package
			regexCallCompiler, err := regexp.Compile(PATH_REGEX)
			check(err)
			call := string(regexCallCompiler.Find(data))

			regexIdentationCompilerFinds := regexIdentationCompiler.FindAllIndex(data, -1)


			parseCall(call, regexIdentationCompilerFinds[len(regexIdentationCompilerFinds)-1][1])

			continue
		}
	}

	//	Save everything after tags
	newFile += string(data[currentIndex:])

	//log.Println(newFile)

	unpackDict()

	// for _, line := range fileLines {
	// 	log.Println(line)
	// 	isMatch, err := regexp.MatchString("call|require|import", line)
	// 	check(err)
	// 	if isMatch {
	// 		continue
	// 	} else {
	// 		newFile += line
	// 	}
	// }
}
