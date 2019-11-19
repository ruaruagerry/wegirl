/*
 * @Author: lingguohua
 * @Date: 2019-08-12 20:04:14
 * @Description:
 */

package gamecfg

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func scanCSVFiles(root string) []string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".csv" {
			return nil
		}

		files = append(files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}

func name2Field(name string) string {
	// first char upcase
	name = string(strings.ToUpper(name)[0]) + name[1:]

	name = strings.ReplaceAll(name, "Id", "ID")

	return strings.ReplaceAll(name, "_", "")
}

func fileName2StructName(filename string) string {
	fn := filepath.Base(filename)
	name := fn[:len(fn)-4]

	return name
}

func commentClean(comment string) string {
	comment = strings.ReplaceAll(comment, "\r\n", " ")
	comment = strings.ReplaceAll(comment, "\n", " ")

	return comment
}

func type2GoType(typestring string) string {
	typestring = strings.TrimSpace(typestring)
	tsarray := strings.Split(typestring, ",")
	var ts string
	// always use last qualifier
	ts = strings.TrimSpace(tsarray[len(tsarray)-1])

	switch ts {
	case "unumber":
		return "UNumber"
	case "int":
		return "int"
	case "boolean":
		return "bool"
	case "number":
		return "float32"
	case "float":
		return "float32"
	}

	// default
	return "string"
}

func genSubStructType(sb1 *strings.Builder, structName string, begin int, columns int,
	comments []string, typestrings []string, names []string) {

	sb1.WriteString(fmt.Sprintf("// %s TODO: game config struct\n", structName))
	sb1.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	for i := begin; i < (begin + columns); i++ {
		name := names[i]
		if name == "" {
			continue
		}

		row := fmt.Sprintf("    %s %s `csv:\"%s\"` // %s\n",
			name2Field(name),
			type2GoType(typestrings[i]),
			name,
			commentClean(comments[i]))
		sb1.WriteString(row)
	}

	sb1.WriteString(fmt.Sprintf("}\n"))
}

func genDefinitionFromFile(csvFile string, structName string) string {
	csvFileReader, err := os.OpenFile(csvFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFileReader.Close()

	br := bufio.NewReader(csvFileReader)
	rune1, _, err := br.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	if rune1 != '\uFEFF' {
		br.UnreadRune() // Not a BOM -- put the rune back
	}

	r := csv.NewReader(br)
	// skip first row
	comments, err := r.Read()
	if err != nil {
		logrus.Errorf("skip first line, read csv file:%s failed:%v", csvFile, err)
		return ""
	}

	typestrings, err := r.Read()
	if err != nil {
		logrus.Errorf("read type string, read csv file:%s failed:%v", csvFile, err)
		return ""
	}

	names, err := r.Read()
	if err != nil {
		logrus.Errorf("read column name, read csv file:%s failed:%v", csvFile, err)
		return ""
	}

	var sb1 strings.Builder
	var sb2 strings.Builder
	sb2.WriteString(fmt.Sprintf("// %s TODO: game config struct\n", structName))
	sb2.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	for i := 0; i < len(names); i++ {
		name := names[i]
		if name == "" {
			continue
		}

		// if name is has prefix 'arr_', then is part of array
		if strings.HasPrefix(name, "arr_") {
			h1 := name[4:]
			index := strings.Index(h1, "_")
			arrName := h1[:index]
			fieldType := name2Field(arrName)

			// skip all array column
			j := i + 1
			for ; j < len(names); j++ {
				skip := names[j]
				if len(skip) < (4 + index) {
					break
				}

				if skip[:4+index] == name[:4+index] {
					continue
				}

				break
			}

			// find each array element span columns
			columns := 1
			commentPrefix := comments[i]
			// find numberic character
			for numbericIdx, r := range commentPrefix {
				if r == '1' {
					commentPrefix = commentPrefix[:numbericIdx+1]
				}
			}

			for k := i + 1; k < len(comments); k++ {
				if strings.HasPrefix(comments[k], commentPrefix) {
					columns++

					continue
				}

				break
			}

			var typeName string
			if columns > 1 {
				subStructType := fmt.Sprintf("%sOf%s", fieldType, structName)
				// need to generate sub-struct-type
				genSubStructType(&sb1, subStructType, i, columns,
					comments, typestrings, names)
				// pointer
				typeName = "*" + subStructType
			} else {
				// normal built-in type, e.g. string
				typeName = type2GoType(typestrings[i])
			}

			row := fmt.Sprintf("    %sArr []%s `csv:\"%s\"` // %s\n",
				fieldType,
				typeName,
				arrName,
				commentClean(comments[i]))
			sb2.WriteString(row)

			i = j - 1
		} else {
			row := fmt.Sprintf("    %s %s `csv:\"%s\"` // %s\n",
				name2Field(name),
				type2GoType(typestrings[i]),
				name,
				commentClean(comments[i]))
			sb2.WriteString(row)
		}
	}

	sb2.WriteString(fmt.Sprintf("}\n"))
	return sb1.String() + sb2.String()
}

// Gen generate go file
func Gen(dir string, outdir string) {
	logrus.Printf("Gen, dir:%s, outdir:%s", dir, outdir)

	// read dir, scan each csv file
	csvFiles := scanCSVFiles(dir)
	var sb strings.Builder
	sb.WriteString("package gamecfg\n\n")

	validCsvFiles := make([]string, 0, len(csvFiles))

	for _, csvFile := range csvFiles {
		if strings.Contains(csvFile, "string") {
			continue
		}

		logrus.Printf("Gen, csv file:%s", csvFile)

		structNameLower := fileName2StructName(csvFile)
		structNameUpper := string(strings.ToUpper(structNameLower)[0]) + structNameLower[1:]

		gen := genDefinitionFromFile(csvFile, structNameUpper)

		logrus.Println(gen)

		if len(gen) > 0 {
			sb.WriteString(gen)

			sb.WriteString("\n")

			validCsvFiles = append(validCsvFiles, csvFile)
		}
	}

	final := sb.String()
	if len(final) > 0 {
		err := ioutil.WriteFile(path.Join(outdir, "game_cfgs.go"), []byte(final), 0644)
		if err != nil {
			logrus.Panic("Gen write to file failed:", err)
		}
		logrus.Println("generate game_cfgs.go file completed")
	}

	// gen load functions
	if len(validCsvFiles) > 0 {
		genLoadFunc(validCsvFiles, outdir)
	}
}

func genLoadFunc(csvFiles []string, outdir string) {
	var sb strings.Builder

	sb.WriteString("package gamecfg\n\n")
	// imports
	sb.WriteString("import (\n")
	sb.WriteString("    \"encoding/csv\"\n")
	sb.WriteString("    \"bufio\"\n")
	sb.WriteString("    \"github.com/sirupsen/logrus\"\n")
	sb.WriteString(")\n")

	// gen variables
	sb.WriteString("var (\n")
	for _, csvFile := range csvFiles {
		structNameLower := fileName2StructName(csvFile)
		structNameUpper := string(strings.ToUpper(structNameLower)[0]) + structNameLower[1:]
		variableName := structNameUpper + "Map"
		variableArrayName := structNameUpper + "Array"
		sb.WriteString(fmt.Sprintf("    // %s TODO\n", variableName))
		sb.WriteString(fmt.Sprintf("    %s map[string]*%s\n", variableName, structNameUpper))
		sb.WriteString(fmt.Sprintf("    %s []*%s\n", variableArrayName, structNameUpper))
	}
	sb.WriteString(")\n")

	// load func
	for _, csvFile := range csvFiles {
		structNameLower := fileName2StructName(csvFile)
		structNameUpper := string(strings.ToUpper(structNameLower)[0]) + structNameLower[1:]
		variableName := structNameUpper + "Map"
		variableArrayName := structNameUpper + "Array"

		sb.WriteString(fmt.Sprintf("func load%s(file *bufio.Reader, p *parser) {\n", structNameUpper))
		sb.WriteString("    csvReader := csv.NewReader(file)\n")
		sb.WriteString("    csvReader.Read()\n")
		sb.WriteString("    csvReader.Read()\n")

		sb.WriteString("    header, err := csvReader.Read()\n")
		sb.WriteString("    if err != nil {\n")
		sb.WriteString("        logrus.Panic(err)\n")
		sb.WriteString("    }\n")

		// to map
		sb.WriteString(fmt.Sprintf("    %s = make(map[string]*%s)\n", variableName, structNameUpper))
		sb.WriteString("    for {\n")
		sb.WriteString("        row, err := csvReader.Read()\n")
		sb.WriteString("        if err != nil {\n")
		sb.WriteString("            break\n")
		sb.WriteString("        }\n\n")
		sb.WriteString(fmt.Sprintf("        av := &%s{}\n", structNameUpper))
		sb.WriteString("        p.unmarshalRow(header, row, av)\n")
		sb.WriteString(fmt.Sprintf("        %s[av.CfgID] = av\n", variableName))
		sb.WriteString(fmt.Sprintf("        %s = append(%s,av)\n", variableArrayName, variableArrayName))
		sb.WriteString("    }\n")
		sb.WriteString("}\n")
	}

	// load function map
	sb.WriteString("func init(){\n")
	for _, csvFile := range csvFiles {
		filename := filepath.Base(csvFile)
		structNameLower := fileName2StructName(csvFile)
		structNameUpper := string(strings.ToUpper(structNameLower)[0]) + structNameLower[1:]

		sb.WriteString(fmt.Sprintf("    loadFuncMap[\"%s\"] = load%s\n", filename, structNameUpper))
	}

	sb.WriteString("}\n")

	// output to file
	final := sb.String()
	if len(final) > 0 {
		err := ioutil.WriteFile(path.Join(outdir, "game_cfgs_load.go"), []byte(final), 0644)
		if err != nil {
			logrus.Panic("Gen write to file failed:", err)
		}
		logrus.Println("generate game_cfgs_load.go file completed")
	}
}
