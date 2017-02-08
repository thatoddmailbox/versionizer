package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

var cssRegex *regexp.Regexp
var jsRegex *regexp.Regexp

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func calculateHash(path string) string {
	// TODO: some sort of caching or memoization
	hasher := sha256.New()
	fileData, err := ioutil.ReadFile(path) // TODO: don't read the whole file into memory
	check(err)
	hasher.Write(fileData)
	return hex.EncodeToString(hasher.Sum(nil))[:8]
}

func versionize(filePath string, basePath string) {
	byteData, err := ioutil.ReadFile(filePath) // TODO: don't read the whole file into memory
	check(err)

	// TODO: stop switch between byte arrays and strings
	data := string(byteData)

	cssMatches := cssRegex.FindAllSubmatch(byteData, -1)

	for _, v := range cssMatches {
		fullThing := string(v[0])

		filePath := path.Join(basePath, "css/" + string(v[1]))
		fileHash := calculateHash(filePath)

		oldPath := "css/" + string(v[1])
		newPath := oldPath + "?v=" + fileHash

		newFullThing := strings.Replace(fullThing, oldPath, newPath, -1)

		data = strings.Replace(data, fullThing, newFullThing, -1)
	}

	jsMatches := jsRegex.FindAllSubmatch(byteData, -1)

	for _, v := range jsMatches {
		fullThing := string(v[0])

		filePath := path.Join(basePath, "js/" + string(v[1]))
		fileHash := calculateHash(filePath)

		oldPath := "js/" + string(v[1])
		newPath := oldPath + "?v=" + fileHash

		newFullThing := strings.Replace(fullThing, oldPath, newPath, -1)

		data = strings.Replace(data, fullThing, newFullThing, -1)
	}

	ioutil.WriteFile(filePath, []byte(data), 0777)
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("usage: versionizer <directory>\n");
		return;
	}
	dirPath := os.Args[1]
	pathStat, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		fmt.Printf("Path '%s' doesn't exist!\n", dirPath);
		return;
	}
	check(err)
	if !pathStat.IsDir() {
		fmt.Printf("Path '%s' is not a directory!\n", dirPath);
		return;
	}

	// TODO: make these less specific
	cssRegex, err = regexp.Compile("<link rel=\"stylesheet\" href=\"css/(.*)\" />")
	check(err)

	jsRegex, err = regexp.Compile("<script src=\"js/(.*)\"></script>")
	check(err)

	// TODO: more than just app.html
	appPath := path.Join(dirPath, "app.html")
	_, err = os.Stat(appPath)
	if os.IsNotExist(err) {
		fmt.Printf("Could not find app.html!\n");
		return;
	}
	check(err)

	versionize(appPath, dirPath)
}