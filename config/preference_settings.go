package config

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func preferenceFuncChoice() string {
	prompt := promptui.Select{
		Label: "请选择要设置的偏好",
		Items: []string{"IP 反向解析（rdns）", "IP 地理位置数据源", "Route-Path", "默认路由跟踪显示模式"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}
	return result
}

func (tc *tracerConfig) preferenceSettings() {
	if rc := preferenceFuncChoice(); rc != "" {
		switch rc {
		case "IP 反向解析（rdns）":
			prompt := promptui.Select{
				Label: "是否默认开启 IP 反向解析",
				Items: []string{"Yes", "No"},
			}

			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
			}
			switch result {
			case "Yes":
				tc.Preference.NoRDNS = false
			case "No":
				tc.Preference.NoRDNS = true
			}

		case "IP 地理位置数据源":
			prompt := promptui.Select{
				Label: "请选择您默认想要使用的 IP 地理位置数据源",
				Items: []string{"LeoMoeAPI", "IPInfo", "IPInsight", "IP.SB", "IPAPI.COM", "IPWHOIS"},
			}

			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
			}
			tc.Preference.DataOrigin = result

		case "Route-Path":
			prompt := promptui.Select{
				Label: "是否默认开启 Route Path 功能",
				Items: []string{"Yes", "No"},
			}

			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
			}
			switch result {
			case "Yes":
				tc.Preference.AlwaysRoutePath = true
			case "No":
				tc.Preference.AlwaysRoutePath = false
			}

		case "默认路由跟踪显示模式":
			prompt := promptui.Select{
				Label: "请选择路由跟踪默认使用的显示模式",
				Items: []string{"实时模式", "制表报告模式"},
			}

			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
			}
			switch result {
			case "制表报告模式":
				tc.Preference.TablePrintDefault = true
			case "实时模式":
				tc.Preference.TablePrintDefault = false
			}
		}
	}
}
