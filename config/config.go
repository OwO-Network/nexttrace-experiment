package config

import "os"

type tracerConfig struct {
	Token      `yaml:"Token"`
	Preference `yaml:"Preference"`
	WebAPI     `yaml:"WebAPI"`
}

type WebAPI struct {
	APIToken string `yaml:"APIToken"`
}

type Token struct {
	LeoMoeAPI string `yaml:"LeoMoeAPI"`
	IPInfo    string `yaml:"IPInfo"`
	IPInsight string `yaml:"IPInsight"`
}

type Preference struct {
	NoRDNS            bool   `yaml:"NoRDNS"`
	DataOrigin        string `yaml:"DataOrigin"`
	AlwaysRoutePath   bool   `yaml:"AlwaysRoutePath"`
	TablePrintDefault bool   `yaml:"TablePrintDefault"`
	TraceMethod       string `yaml:"TraceMethod"`
}

type configPath func() (string, error)

func configFromRunDir() (string, error) {
	return "./", nil
}

func configFromUserHomeDir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return dir + "/.nexttrace/", nil
}
