package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path-to-csproj>")
		return
	}

	csprojPath := os.Args[1]
	baseName := strings.TrimSuffix(filepath.Base(csprojPath), filepath.Ext(csprojPath))

	// .csprojファイルの解析
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(csprojPath); err != nil {
		fmt.Printf("Error reading csproj file: %v\n", err)
		return
	}

	project := doc.SelectElement("Project")
	if project == nil {
		fmt.Println("No <Project> element found in the csproj file.")
		return
	}

	// <PropertyGroup>の子エレメントを追加または更新
	propertyGroups := project.SelectElements("PropertyGroup")
	for _, propertyGroup := range propertyGroups {
		updateOrCreateElement(propertyGroup, "OutputType", "Exe")
		updateOrCreateElement(propertyGroup, "PackAsTool", "true")
		updateOrCreateElement(propertyGroup, "ToolCommandName", baseName)
		updateOrCreateElement(propertyGroup, "PackageOutputPath", "./nupkg")
	}

	// 変更を保存 (インデントを設定)
	doc.Indent(2)
	if err := doc.WriteToFile(csprojPath); err != nil {
		fmt.Printf("Error writing csproj file: %v\n", err)
		return
	}

	// バッチファイルの作成
	batchFilePath := filepath.Join(filepath.Dir(csprojPath), "dotnet_tool_publish.bat")
	batchFileContent := fmt.Sprintf(`dotnet pack
dotnet tool install --global --add-source ./nupkg %s`, baseName)

	if err := os.WriteFile(batchFilePath, []byte(batchFileContent), 0644); err != nil {
		fmt.Printf("Error writing batch file: %v\n", err)
		return
	}

	fmt.Println("dotnet_tool_publish.bat created successfully.")
}

func updateOrCreateElement(parent *etree.Element, tag, text string) {
	element := parent.SelectElement(tag)
	if element == nil {
		element = parent.CreateElement(tag)
	}
	element.SetText(text)
}
