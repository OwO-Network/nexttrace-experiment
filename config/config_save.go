package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func (tc *tracerConfig) saveConfig() (error) {
	yamlData, err := yaml.Marshal(&tc)

	if err != nil {
		return err
	}

	if err = writeFile(yamlData); err != nil {
		return err
	} else {
		fmt.Println("配置文件已经更新，在下次路由跟踪时，将会使用您的偏好。")
		return nil
	}
}
