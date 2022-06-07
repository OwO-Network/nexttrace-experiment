package config

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

func webAPIFuncChoice() string {
	prompt := promptui.Select{
		Label: "请选择要设置的偏好",
		Items: []string{"Web API Token 设置", "监听端口设置"},
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
		case "监听端口设置":
			validate := func(input string) error {
				p, err := strconv.ParseFloat(input, 64)
				if err != nil {
					return errors.New("端口必须为纯数字！")
				}
				if p < 1 || p > 65535 {
					return errors.New("端口号不合法！")
				}
				return nil
			}

			prompt := promptui.Prompt{
				Label:    "请输入端口号",
				Validate: validate,
			}

			result, err := prompt.Run()

			if err != nil {
				fmt.Printf("取消设置 %v\n", err)
				return
			}
			tc.WebAPI.ListenPort, _ = strconv.Atoi(result)
		}
	}
}
