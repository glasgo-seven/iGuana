package iguana

// Implementation of the Ternary Operator
func ternary(statement bool, ifTrue interface{}, ifFalse interface{}) interface{} {
	if statement {
		return ifTrue
	} else {
		return ifFalse
	}
}

//	Prints the contents of the Template Dictionary
func unpackDict() {
	for packageName, templates := range IMPORT {
		println(packageName, " {")
		for templateName, template := range templates {
			println("\t", templateName, " : ", template.file, template.name, template.isAccessed)
		}
		println("}")
	}
}
