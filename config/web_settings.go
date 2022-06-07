package config

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func webAPIFuncChoice() string {
	prompt := promptui.Select{
		Label: "请选择要设置的偏好",
		Items: []string{"Web API Token 设置"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("取消设置 %v\n", err)
		return ""
	}
	return result
}

func (tc *tracerConfig) webSettings() {
	if cr := webAPIFuncChoice(); cr != "" {
		switch cr {
		case "Web API Token 设置":
			prompt := promptui.Prompt{
				Label: "Web API Token",
			}

			result, err := prompt.Run()

			if err != nil {
				fmt.Printf("取消设置 %v\n", err)
				return
			}

			if result == "" {
				result = "NextTrace"
			}

			tc.WebAPI.APIToken = result
		}
	}
}
