package gfunc

import (
	"bufio"
	"os"
	"strings"

	"github.com/lisijie/wordfilter/trie"
	log "github.com/sirupsen/logrus"
)

var (
	myTrie *trie.Trie
)

// LoadSensitiveWordDictionary 加载敏感词文件
func LoadSensitiveWordDictionary(filename string) {
	if filename == "" {
		return
	}

	if myTrie == nil {
		myTrie = trie.NewTrie()
	}

	log.Println("Load Sensitiveword, filename:", filename)

	file, err := os.Open(filename)
	if err != nil {
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var contentString = scanner.Text()
		var ws = strings.Split(contentString, ",")
		for _, w := range ws {
			word := strings.TrimSpace(w)
			if word != "" {
				myTrie.Add(word)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return
}

// ReplaceSensitiveWord 替换敏感词
func ReplaceSensitiveWord(word string) (bool, string) {
	if myTrie == nil {
		return false, word
	}

	keyword, _ := myTrie.Replace(word)

	return false, keyword
}
