package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//we assume the enum type name use CamelCase convention but the first letter is upper case
func extractPrefix(s string) string {
	var sb strings.Builder
	for i, r := range s {
		if i == 0 || unicode.IsUpper(r) {
			sb.WriteRune(unicode.ToUpper(r))
		}
	}
	return sb.String()
}

func upcaseFirstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(strings.ToUpper(string(s[0])))
	sb.WriteString(s[1:])
	return sb.String()
}

func lowcaseFirstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(strings.ToLower(string(s[0])))
	sb.WriteString(s[1:])
	return sb.String()
}

func buildEnumFromValue(s string) string {
	if len(s) == 0 {
		return ""
	}
	var sb strings.Builder
	fields := strings.FieldsFunc(s, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	for _, f := range fields {
		sb.WriteString(upcaseFirstLetter(f))
	}
	return sb.String()
}

func writeTypeSource(f *os.File, typeName string, enumNames []string, enumValues []string) {
	if len(enumNames) != len(enumValues) {
		return
	}
	f.WriteString("package main\n")
	f.WriteString("\nimport (")
	f.WriteString(fmt.Sprintf("\n\t%q", "fmt"))
	f.WriteString("\n)\n")
	f.WriteString(fmt.Sprintf("\n//%s Enumeration type", typeName))
	f.WriteString(fmt.Sprintf("\ntype %s int\n", typeName))

	f.WriteString("\nconst (")
	for i := 0; i < len(enumNames); i++ {
		f.WriteString(fmt.Sprintf("\n\t// %s = %q", enumNames[i], enumValues[i]))
		f.WriteString(fmt.Sprint("\n\t", enumNames[i]))
		if i == 0 {
			f.WriteString(" int = iota")
		}
	}
	f.WriteString("\n)\n")
	f.Sync()
}

func writeTypeSourceForString(f *os.File, typeName string, enumNames []string, enumValues []string) {
	if len(enumNames) != len(enumValues) {
		return
	}
	f.WriteString("package main\n")
	f.WriteString(fmt.Sprintf("\n//%s Enumeration type", typeName))
	f.WriteString(fmt.Sprintf("\ntype %s string\n", typeName))

	f.WriteString("\nconst (")
	for i := 0; i < len(enumNames); i++ {
		f.WriteString(fmt.Sprintf("\n\t// %s = %q", enumNames[i], enumValues[i]))
		f.WriteString(fmt.Sprintf("\n\t%s string = %q", enumNames[i], enumValues[i]))
	}
	f.WriteString("\n)\n")
	f.Sync()
}

func writeConstructorSource(f *os.File, typeName string, enumNames []string, enumValues []string) {
	if len(enumNames) != len(enumValues) {
		return
	}
	f.WriteString(fmt.Sprintf("\n//New%s : Construct a new %s Object", typeName, typeName))
	f.WriteString(fmt.Sprintf("\nfunc New%s(s string) (%s, error) {", typeName, typeName))
	f.WriteString(fmt.Sprintf("\n\tvar r %s", typeName))
	f.WriteString("\n\tswitch s {")
	for i := 0; i < len(enumNames); i++ {
		f.WriteString(fmt.Sprintf("\n\tcase %q:", enumValues[i]))
		f.WriteString(fmt.Sprintf("\n\t\tr = %s(%s)", typeName, enumNames[i]))
	}
	f.WriteString("\n\tdefault:")
	f.WriteString(fmt.Sprintf("\n\t\treturn %s(%s), fmt.Errorf(\"%%q is not a valid %s\", s)", typeName, enumNames[0], typeName))
	f.WriteString("\n\t}")
	f.WriteString("\n\treturn r, nil")
	f.WriteString("\n}\n")
	f.Sync()
}

func writeStringSource(f *os.File, typeName string, enumNames []string, enumValues []string) {
	if len(enumNames) != len(enumValues) {
		return
	}
	f.WriteString(fmt.Sprintf("\nfunc (enum %s) String() string {", typeName))
	f.WriteString(fmt.Sprintf("\n\tvar r string"))
	f.WriteString("\n\tswitch int(enum) {")
	for i := 0; i < len(enumNames); i++ {
		f.WriteString(fmt.Sprintf("\n\tcase %s:", enumNames[i]))
		f.WriteString(fmt.Sprintf("\n\t\tr = %q", enumValues[i]))
	}
	f.WriteString("\n\t}")
	f.WriteString("\n\treturn r")
	f.WriteString("\n}\n")
	f.Sync()
}

func main() {
	//check command line parameter
	if len(os.Args) < 3 {
		fmt.Println("\nusage:\n\tgoEnumGenerator input-file-name 0 or 1 (0 for int enum, 1 for string enum)")
		return
	}

	if !strings.Contains("01", os.Args[2]) {
		fmt.Println("we only support two type enums: 0->int 1->string")
		return
	}

	//read data into memory
	fn := os.Args[1]
	fIn, err := os.Open(fn)
	check(err)
	defer fIn.Close()

	enumType := strings.Compare("0", os.Args[2])
	if enumType == 0 {
		fmt.Println("Generate enum use int as under type")
	} else {
		fmt.Println("Generate enum use string as under type")
	}

	scanner := bufio.NewScanner(fIn)

	var typeName string
	var prefix string
	enumValues := make([]string, 0, 32)
	enumNames := make([]string, 0, 32)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " +,.-~!@#$%^&*();:'")
		if i == 0 {
			typeName = line
			prefix = extractPrefix(typeName)
			i++
		} else {
			enumValues = append(enumValues, line)
			enumNames = append(enumNames, prefix+buildEnumFromValue(line))
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	//output to a new file
	fn = lowcaseFirstLetter(typeName) + ".go"
	fOut, err := os.Create(fn)
	check(err)
	defer fOut.Close()
	if enumType == 0 {
		writeTypeSource(fOut, typeName, enumNames, enumValues)
		writeConstructorSource(fOut, typeName, enumNames, enumValues)
		writeStringSource(fOut, typeName, enumNames, enumValues)
	} else {
		writeTypeSourceForString(fOut, typeName, enumNames, enumValues)
	}

	fOut.Sync()
}
