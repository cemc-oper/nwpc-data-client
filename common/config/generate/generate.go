package main

import (
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

//go:generate go run generate.go
func main() {
	fmt.Printf("This a generate tool to embed YAML config files into nwpc-data-client.\n")

	configDir := filepath.Join("../../../", "conf")

	file, err := os.Create("../config.autogen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	type configItem struct {
		Name    string
		Content string
	}

	var configList []string

	_ = filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != common.ConfigFileBasename {
			return nil
		}

		relFilePath, _ := filepath.Rel(configDir, path)
		dataTypeFile := filepath.ToSlash(relFilePath)
		dataType := dataTypeFile[:len(dataTypeFile)-5]
		fmt.Printf("%s\n", dataType)

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		configList = append(configList, configString(dataType, string(content)))
		return nil
	})

	fmt.Fprint(file, `// this is a auto generated file.
package config

var EmbeddedConfigs = [][2]string{
`)

	for _, item := range configList {
		fmt.Fprint(file, "    ", item, " ,\n")
	}

	fmt.Fprintf(file, "}\n")
}

func configString(name, content string) string {
	return fmt.Sprintf("{`%s`, `%s`}", name, content)
}
