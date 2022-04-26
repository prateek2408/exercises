package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/zclconf/go-cty/cty"
)

func main() {
	fmt.Println("hello WOrld")
	basePath, err := os.Getwd()
	// var newFileContents []byte
	if err != nil {
		fmt.Println("unable to get current directory")
	}
	// if newFileContents, err = ChangeValue("main.tf", basePath, "apigee-x-core", "source", "mysource"); err != nil {
	// 	fmt.Println("error in ChangeValue ", err.Error())
	// }
	// os.WriteFile("main.tf", newFileContents, 0644)
	GetTFKeys(basePath)
}

func ChangeValue(fileName, filePath, label, key, value string) ([]byte, error) {
	targetFile := filepath.Join(filePath, fileName)
	file, err := ioutil.ReadFile(targetFile)
	var indexOfBlock int8 = 0

	if err != nil {
		return nil, err
	}

	hclFile, hclErrors := hclwrite.ParseConfig(file, fileName, hcl.InitialPos)
	if hclErrors.HasErrors() {
		return nil, hclErrors
	}

	if indexOfBlock = getBlockIndex(hclFile, label); indexOfBlock == -1 {
		return nil, fmt.Errorf("unable to find any block with given label")
	}

	blockToChange := hclFile.Body().Blocks()[indexOfBlock]
	blockToChange.Body().SetAttributeValue(key, cty.StringVal(value))

	return hclFile.Bytes(), nil
}

func getBlockIndex(file *hclwrite.File, label string) int8 {
	var blockIndex int8 = -1
	for index, block := range file.Body().Blocks() {
		if contains(label, block.Labels()) {
			blockIndex = int8(index)
			break
		}
	}
	return blockIndex
}

func contains(word string, words []string) bool {
	i := sort.SearchStrings(words, word)
	return i < len(words) && words[i] == word
}

func GetTFKeys(filePath string) (map[string]string, error) {

	module, diags := tfconfig.LoadModule(filePath)
	if diags.HasErrors() {
		return nil, diags
	}
	tfKeys := make(map[string]string, len(module.Variables))
	for _, variable := range module.Variables {
		tfKeys[variable.Name] = variable.Type
	}
	return tfKeys, nil
}
