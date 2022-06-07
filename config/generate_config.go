package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func writeFile(content []byte) error {
	var err error
	var path string
	path, err = configFromUserHomeDir()
	if err != nil {
		path, err = configFromRunDir()
		if err != nil {
			return err
		}
	}

	if exist, _ := pathExists(path); !exist {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if err = ioutil.WriteFile(path+"ntraceConfig.yml", []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func AutoGenerate() (*tracerConfig, error) {
	token := Token{
		LeoMoeAPI: "LeoOwO",
		IPInfo:    "",
	}

	preference := Preference{
		AlwaysRoutePath:   false,
		TablePrintDefault: false,
		DataOrigin:        "LeoMoeAPI",
	}

	finalConfig := tracerConfig{
		Token:      token,
		Preference: preference,
	}

	yamlData, err := yaml.Marshal(&finalConfig)

	if err != nil {
		return nil, err
	}

	if err = writeFile(yamlData); err != nil {
		return nil, err
	} else {
		return &finalConfig, nil
	}
}

func Generate() error {

	fmt.Println("欢迎使用高阶自定义功能，这是一个配置向导，我们会帮助您生成配置文件。\n您的配置文件会被放在 ~/.nexttrace/ntraceConfig.yml 中，您也可以通过编辑这个文件来自定义配置。")

	tc, err := Read()

	// Initialize Default Config
	if err != nil || tc.DataOrigin == "" {
		if tc, err = AutoGenerate(); err != nil {
			log.Fatal(err)
		}
	}

	for {
		prompt := promptui.Select{
			Label: "请选择功能",
			// Items: []string{"Token设置", "路由跟踪偏好设置", "快速路由测试设置", "保存并退出"},
			Items: []string{"Token设置", "路由跟踪偏好设置", "Web API 设置", "不保存推出", "保存并退出"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("取消设置 %v\n", err)
		}

		switch result {
		case "Token 设置":
			tc.tokenSettings()
		case "路由跟踪偏好设置":
			tc.preferenceSettings()
		case "Web API 设置":
			tc.webSettings()
		case "不保存并退出":
			return nil
		case "保存并退出":
			if err := tc.saveConfig(); err != nil {
				return err
			}
			return nil
		}
	}
}
