package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	defaultRepository = "https://repo1.maven.org/maven2"
)

type Dependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type Repository struct {
	URL string `xml:"url"`
}

type Project struct {
	Dependencies []Dependency `xml:"dependencies>dependency"`
	Repositories []Repository `xml:"repositories>repository"`
}

func main() {
	args := os.Args[1:]

	if len(args) <= 0 {
		log.Println("Usage: ./downloader <pom.xml>")
		return
	}

	filePath := args[0]

	log.Println("Selected file:", filePath)

	if strings.HasSuffix(filePath, "pom.xml") {
		downloadMaven(filePath)
	} else {
		log.Fatal("Gradle is not supported yet")
	}

}

func downloadMaven(path string) {
	folderPath := strings.Replace(path, "pom.xml", "", 1) + "dependencies/"
	_, err := os.Stat(folderPath)

	if os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var project Project

	decoder := xml.NewDecoder(file)
	if err := decoder.Decode(&project); err != nil {
		log.Fatal(err)
	}

	repositories := mavenRepositories(project)

	for _, dep := range project.Dependencies {
		fmt.Println(dep.GroupID, dep.ArtifactID, dep.Version)
		downloadPath := artifactPath(dep)
		outputPath := folderPath + dep.ArtifactID + "-" + dep.Version + ".jar"

		if err := downloadFromRepositories(repositories, downloadPath, outputPath); err != nil {
			log.Fatal(err)
		}
	}
}

func mavenRepositories(project Project) []string {
	repositories := []string{defaultRepository}
	for _, repo := range project.Repositories {
		repositories = appendRepository(repositories, repo.URL)
	}

	return repositories
}

func appendRepository(repositories []string, repository string) []string {
	repository = strings.TrimSpace(repository)
	if repository == "" {
		return repositories
	}

	repository = strings.TrimRight(repository, "/")
	for _, existing := range repositories {
		if existing == repository {
			return repositories
		}
	}

	return append(repositories, repository)
}

func artifactPath(dep Dependency) string {
	groupPath := strings.ReplaceAll(dep.GroupID, ".", "/")
	return fmt.Sprintf("%s/%s/%s/%s-%s.jar", groupPath, dep.ArtifactID, dep.Version, dep.ArtifactID, dep.Version)
}

func downloadFromRepositories(repositories []string, artifactPath string, filepath string) error {
	var lastErr error
	for _, repository := range repositories {
		downloadURL := repository + "/" + artifactPath
		if err := DownloadFile(downloadURL, filepath); err == nil {
			return nil
		} else {
			lastErr = err
			log.Printf("Could not download from %s: %v", repository, err)
		}
	}

	return fmt.Errorf("failed to download %s from configured repositories: %w", artifactPath, lastErr)
}

func DownloadFile(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
