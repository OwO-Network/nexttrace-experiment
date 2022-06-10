package config

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func tokenFuncChoice() string {
	prompt := promptui.Select{
		Label: "请选择要设置的Token",
		Items: []string{"LeoMoeAPI", "IPInfo", "IPInsight", "IP地理位置校准密钥(开发者使用)"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}
	return result
}

func (tc *tracerConfig) tokenSettings() {

	if cr := tokenFuncChoice(); cr != "" {
		switch cr {
		case "LeoMoeAPI":
			prompt := promptui.Prompt{
				Label: "LeoMoeAPI Token",
			}

			result, err := prompt.Run()

			if err != nil {
				fmt.Printf("取消设置 %v\n", err)
				return
			}

			if result == "" {
				result = "LeoOwO"
			}

			tc.Token.LeoMoeAPI = result
		case "IPInfo":
			prompt := promptui.Prompt{
				Label: "IPInfo Token",
			}

			result, err := prompt.Run()

			if err != nil {
				fmt.Printf("取消设置 %v\n", err)
				return
			}

			tc.Token.IPInfo = result
		case "IPInsight":
			prompt := promptui.Prompt{
				Label: "IPInsight Token",
			}

			result, err := prompt.Run()

			if err != nil {
				fmt.Printf("取消设置 %v\n", err)
				return
			}

			tc.Token.IPInsight = result
		case "IP地理位置校准密钥(开发者使用)":
			prompt := promptui.Prompt{
				Label: "LeoMoe Update 密钥",
			}

			result, err := prompt.Run()

			if err != nil {
				fmt.Printf("取消设置 %v\n", err)
				return
			}

			tc.Token.LeoMoeUpdateKey = result
		}
	}
}
