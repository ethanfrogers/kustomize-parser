package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type KustomizeFile struct {
	Resources           []string             `yaml:"resources"`
	Patches             []Patch              `yaml:"patches"`
	ConfigMapGenerators []ConfigMapGenerator `yaml:"configMapGenerator"`
}

type ConfigMapGenerator struct {
	Name  string   `yaml:"name"`
	Files []string `yaml:"files"`
}

type Patch struct {
	Path string `yaml:"path"`
}

// FindAndParseKustomizeFile handles building a list of dependencies for a given
// kustomization.yml file
func FetchAndParseKustomizeFile(pth string) (*KustomizeFile, []string, error) {
	basePath := filepath.Dir(pth)
	fileContents, err := ioutil.ReadFile(pth)
	if err != nil {
		return nil, nil, err
	}

	var kfile KustomizeFile
	if err := yaml.Unmarshal(fileContents, &kfile); err != nil {
		return nil, nil, err
	}

	requestedPaths := []string{}
	// append patches if path is non empty.
	for _, patch := range kfile.Patches {
		if patch.Path != "" {
			requestedPaths = append(requestedPaths, filepath.Join(basePath, patch.Path))
		}
	}

	// append config map generator files
	for _, cmg := range kfile.ConfigMapGenerators {
		for _, item := range cmg.Files {
			requestedPaths = append(requestedPaths, filepath.Join(basePath, item))
		}
	}

	// handle resource paths
	for _, resource := range kfile.Resources {
		if ext := filepath.Ext(resource); ext == "" {
			relPath := filepath.Join(basePath, resource)
			// this file can be either .yml, .yaml or a raw file called kustomization
			// in real life we'll need to check for all 3
			searchPath := filepath.Join(
				relPath, "kustomization.yml")
			_, pths, err := FetchAndParseKustomizeFile(searchPath)
			if err != nil {
				return nil, nil, err
			}
			for _, p := range pths {
				requestedPaths = append(requestedPaths, p)
			}
			continue
		}
		requestedPaths = append(requestedPaths, filepath.Join(basePath, resource))
	}

	return &kfile, requestedPaths, nil
}

func main() {
	file := flag.String("file", "/Users/ethanrogers/Scratch/kustomize-test/overlays/kustomization.yml", "input kustomize file")
	flag.Parse()

	_, requestedPaths, err := FetchAndParseKustomizeFile(*file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(requestedPaths)
}
