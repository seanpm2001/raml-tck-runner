package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type manifest struct {
	Description string   `json:"description"`
	FilePaths   []string `json:"filePaths"`
}

// ListRamls lists RAML file in :folderPath: folder;
// Uses raml-tck manifest.json file
func ListRamls(folderPath string) ([]string, error) {
	fileList := []string{}
	manifestPath := filepath.Join(folderPath, "manifest.json")
	manifestFile, err := os.Open(manifestPath)
	defer manifestFile.Close()
	if err != nil {
		return fileList, err
	}
	byteValue, _ := ioutil.ReadAll(manifestFile)
	var m manifest
	json.Unmarshal(byteValue, &m)
	for _, fp := range m.FilePaths {
		fileList = append(fileList, filepath.Join(folderPath, fp))
	}
	return fileList, nil
}

// CloneTckRepo clones raml-tck repo and returns cloned repo path
func CloneTckRepo() string {
	targetDir := fmt.Sprintf("%s/raml-tck", os.TempDir())
	_ = os.RemoveAll(targetDir)
	fmt.Printf("Cloning raml-tc repo to %s\n", targetDir)
	gitRepo := "git@github.com:raml-org/raml-tck.git"
	cmd := exec.Command(
		"git", "clone", "-b", "rename-cleanup", gitRepo, targetDir)
	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("Failed to clone repo %s", gitRepo))
	}
	return targetDir
}

// SaveReport writes parsing run report as JSON file
func SaveReport(report *Report, outdir string) {
	outdir, err := filepath.Abs(outdir)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to get absolute path to output dir %s: %s",
			outdir, err.Error()))
	}
	err = os.MkdirAll(outdir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to create output dir at %s: %s",
			outdir, err.Error()))
	}
	repFilePath := filepath.Join(
		outdir, fmt.Sprintf("%s.json", report.Parser))
	reportJson, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to marshal report: %s", err.Error()))
	}
	err = ioutil.WriteFile(repFilePath, reportJson, 0644)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to write report to %s: %s",
			repFilePath, err.Error()))
	}
}
