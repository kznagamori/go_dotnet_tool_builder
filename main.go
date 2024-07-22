package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Project struct {
	XMLName        xml.Name        `xml:"Project"`
	PropertyGroups []PropertyGroup `xml:"PropertyGroup"`
}

type PropertyGroup struct {
	XMLName  xml.Name  `xml:"PropertyGroup"`
	Elements []Element `xml:",any"`
}

type Element struct {
	XMLName xml.Name
	Content string `xml:",chardata"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path to csproj file>")
		return
	}

	csprojPath := os.Args[1]
	csprojDir := filepath.Dir(csprojPath)
	csprojBaseName := filepath.Base(csprojPath)
	toolCommandName := csprojBaseName[:len(csprojBaseName)-len(filepath.Ext(csprojBaseName))]

	// Load csproj file
	csprojData, err := ioutil.ReadFile(csprojPath)
	if err != nil {
		fmt.Println("Error reading csproj file:", err)
		return
	}

	var project Project
	if err := xml.Unmarshal(csprojData, &project); err != nil {
		fmt.Println("Error unmarshalling csproj file:", err)
		return
	}

	// Modify or add elements in PropertyGroup
	for i := range project.PropertyGroups {
		pg := &project.PropertyGroups[i]
		setOrUpdateElement(pg, "OutputType", "Exe")
		setOrUpdateElement(pg, "PackAsTool", "true")
		setOrUpdateElement(pg, "ToolCommandName", toolCommandName)
		setOrUpdateElement(pg, "PackageOutputPath", "./nupkg")
	}

	// Marshal modified project back to XML
	modifiedCsprojData, err := xml.MarshalIndent(project, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling modified csproj file:", err)
		return
	}

	// Write modified csproj file
	if err := ioutil.WriteFile(csprojPath, modifiedCsprojData, 0644); err != nil {
		fmt.Println("Error writing modified csproj file:", err)
		return
	}

	// Create dotnet_tool_publish.bat
	batFilePath := filepath.Join(csprojDir, "dotnet_tool_publish.bat")
	batContent := fmt.Sprintf(`@echo off
dotnet pack
dotnet tool install --global --add-source ./nupkg %s
`, toolCommandName)

	if err := os.WriteFile(batFilePath, []byte(batContent), 0644); err != nil {
		fmt.Println("Error writing batch file:", err)
		return
	}

	fmt.Println("Successfully modified csproj file and created dotnet_tool_publish.bat")
}

func setOrUpdateElement(pg *PropertyGroup, name, value string) {
	for i, elem := range pg.Elements {
		if elem.XMLName.Local == name {
			pg.Elements[i].Content = value
			return
		}
	}
	pg.Elements = append(pg.Elements, Element{
		XMLName: xml.Name{Local: name},
		Content: value,
	})
}
