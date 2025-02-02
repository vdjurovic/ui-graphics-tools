package commands

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"

	"gopkg.in/yaml.v3"
)

const (
	defaultWidth  = "600"
	defaultHeight = "400"
)

var (
	//go:embed templates/splash-template.svg
	splashTemplate string
)

type SplashParams struct {
	ConfigFile  string
	OutputDir   string
	OutFileName string
	Verbose     bool
}

type SplashTemplateData struct {
	Title          string `yaml:"title"`
	Subheading     string `yaml:"subheading"`
	Copyright      string `yaml:"copyright"`
	BackgroundFill string `yaml:"background-fill"`
}

func GenerateSplashScreen(params *SplashParams) error {
	yamlConten, err := os.ReadFile(params.ConfigFile)
	if err != nil {
		return err
	}
	templateData := SplashTemplateData{}
	err = yaml.Unmarshal(yamlConten, &templateData)
	if err != nil {
		return err
	}
	tmpl, err := template.New("splash-template").Parse(splashTemplate)
	if err != nil {
		return err
	}
	tmpFile, err := ioutil.TempFile("", "splash")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	err = tmpl.Execute(tmpFile, templateData)
	outFilepath := path.Join(params.OutputDir, params.OutFileName)
	os.MkdirAll(params.OutputDir, os.ModePerm)
	if err != nil {
		return err
	}
	inkscape := execCommand(inkscapeCmd, "-w", defaultWidth, "-h", defaultHeight, "-o", outFilepath, tmpFile.Name())
	out, err := inkscape.CombinedOutput()
	if err != nil {
		log.Printf("Failed to run inkscape: %s\n", err)
		return err
	}
	if params.Verbose {
		fmt.Println(string(out))
	}
	return err
}
