/*
 * @Author: lingguohua
 * @Date: 2019-08-12 20:04:41
 * @Description:
 */

// Package gamecfg 游戏的策划配置
package gamecfg

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type loadfn func(*bufio.Reader, *parser)

var (
	loadFuncMap = make(map[string]loadfn)
)

// LoadAll load all config from specific path
func LoadAll(dir string) {
	logrus.Printf("LoadAll, dir:%s", dir)

	// read dir, scan each csv file
	csvFiles := scanCSVFiles(dir)

	p := newParser()
	for _, csvFile := range csvFiles {
		loadCsvFile(csvFile, p)
	}
}

func loadCsvFile(csvFile string, p *parser) {
	filen := filepath.Base(csvFile)
	lf, ok := loadFuncMap[filen]
	if !ok {
		logrus.Warn("no loadfn for:", csvFile)
		return
	}

	fileReader, err := os.OpenFile(csvFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer fileReader.Close()

	br := bufio.NewReader(fileReader)
	r, _, err := br.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	if r != '\uFEFF' {
		br.UnreadRune() // Not a BOM -- put the rune back
	}

	lf(br, p)

	logrus.Printf("load csvFile:%s", csvFile)
}
