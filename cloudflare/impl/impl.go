package impl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ned123abc/yunwe/cloudflare"
)

// 生产
const (
	baseURL   = "https://api.cloudflare.com/client/v4"
	authEmail = "xiaohailisa@gmail.com"                  // 替换为你的 Cloudflare 邮箱
	authKey   = "100d77f4735565e80dd6a140a86e934c18989" // 替换为你的 Cloudflare API 密钥
	authAccId = "e76a7a23da86f0bf73c5e322227d0b14"
)


// 获取域名ZoneID
func GetZoneId(varDomain string) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL+"/zones?name="+varDomain, nil)
	req.Header.Set("X-Auth-Email", authEmail)
	req.Header.Set("X-Auth-Key", authKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err // 返回错误
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err // 返回错误
	}

	var zoneResp cloudflare.ZoneResponse
	if err := json.Unmarshal(body, &zoneResp); err != nil {
		return "", err // 返回错误
	}

	if zoneResp.Success {
		for _, zone := range zoneResp.Result {
			if zone.Name == varDomain {
				fmt.Printf("%+v\n", zone.ID)
				return zone.ID, nil
			}
		}
		return "", fmt.Errorf("domain %s not found", varDomain) // 未找到域名时返回错误
	} else {
		return "", fmt.Errorf("failed to fetch zones") // 获取区域失败时返回错误
	}
}

// -----------------------------------------  添加域名解析
func UpdateRecord(varDomain string, cfZoneId string, varRecord string) string {
	// 初始化 DNS 记录数组
	records := []cloudflare.DNSRecord{
		{
			Type:    "CNAME",
			Name:    "@",
			Content: varRecord,
			TTL:     1,
			Proxied: true,
		},
		{
			Type:    "CNAME",
			Name:    "*",
			Content: varRecord,
			TTL:     1,
			Proxied: true,
		},
	}

	client := &http.Client{}
	baseURL := "https://api.cloudflare.com/client/v4"
	var responseBodies []string

	for _, record := range records {
		// 转换为 JSON
		recordJSON, err := json.Marshal(record)
		if err != nil {
			fmt.Printf("Error marshalling DNS record: %+v\n", err)
			return "error"
		}

		// 创建请求
		req, err := http.NewRequest("POST", baseURL+"/zones/"+cfZoneId+"/dns_records", bytes.NewBuffer(recordJSON))
		if err != nil {
			fmt.Printf("Error creating request: %+v\n", err)
			return "error"
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Auth-Email", authEmail) // 替换为实际的邮箱
		req.Header.Set("X-Auth-Key", authKey)     // 替换为实际的认证密钥

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %+v\n", err)
			return "error"
		}
		defer resp.Body.Close()

		// 读取响应
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %+v\n", err)
			return "error"
		}
		responseBodies = append(responseBodies, string(body))
	}

	// 处理响应
	for _, responseBody := range responseBodies {
		fmt.Println("Response:", responseBody)
	}

	return "done"
}

func NewSite(varDomain string) (string, string, bool) {

	url := "https://api.cloudflare.com/client/v4/zones"
	// 创建一个 HTTP 客户端

	// 创建 JSON 数据
	jsonData := fmt.Sprintf(`{
		"account": {
			"id": "%s"
		},
		"name": "%s",
		"type": "full"
	}`, authAccId, varDomain)
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(jsonData))
	if err != nil {
		panic(err)
	}

	// 添加请求头
	req.Header.Add("X-Auth-Email", authEmail)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析 JSON
	var response cloudflare.NewSiteResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// 打印响应
	//fmt.Println("NewSite Response: ", string(body))
	fmt.Println("NewSite Nameserver1: ", response.Result.NameServers[0])
	fmt.Println("NewSite Nameserver2: ", response.Result.NameServers[1])
	fmt.Println("NewSite Success Status: ", response.Success)
	return response.Result.NameServers[0], response.Result.NameServers[1], response.Success
}

// 创建缓存规则ID
func CreateCacheRuleId(zoneID string) string {

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/rulesets", zoneID)
	// 创建一个 HTTP 客户端

	// 创建 JSON 数据
	jsonData := []byte(`{
			"description": "",
			"kind": "zone",
			"name": "default",
			"phase": "http_request_cache_settings"
		}`)
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// 添加请求头
	req.Header.Add("X-Auth-Email", authEmail)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析 JSON
	var response cloudflare.GetCfResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// 打印响应
	fmt.Println("Response: ", string(body))
	fmt.Println("缓存规则ID: ", response.Result.ID)
	return response.Result.ID
}

// ----------------------------------------- 添加缓存规则       配置 /xxxx/  不缓存 绕过缓存
func UpdateCacheRule(zoneID string) string {

	cacheRuleId := CreateCacheRuleId(zoneID)

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/rulesets/%s", zoneID, cacheRuleId)
	// 创建一个 HTTP 客户端

	// 创建 JSON 数据
	jsonData := []byte(`{
			"rules": [
				{
					"expression": "(http.request.uri.path eq \"/xxxx/\")",
					"description": "no_cache",
					"action": "set_cache_settings",
					"action_parameters": {
						"cache": false
					}
				}
			]
		}`)
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// 添加请求头
	req.Header.Add("X-Auth-Email", authEmail)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析 JSON
	var response cloudflare.GetCfResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// 打印响应
	fmt.Println("Response:", string(body))
	return "done"
}

// ------------------------------- 创建防火墙规则ID-------------------------------
func CreateFirewallRuleId(zoneID string) string {

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/rulesets", zoneID)
	// 创建一个 HTTP 客户端

	// 创建 JSON 数据
	jsonData := []byte(`{
			"name": "limit",
			"kind": "zone",
			"description": "country access",
			"phase": "http_request_firewall_custom"
		}`)
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// 添加请求头
	req.Header.Add("X-Auth-Email", authEmail)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析 JSON
	var response cloudflare.GetCfResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// 打印响应
	fmt.Println("Response: ", string(body))
	fmt.Println("缓存规则ID: ", response.Result.ID)
	return response.Result.ID
}

// ----------------------------------------- 添加Waf防火墙规则    防止其他国家攻击  创建一个limit 拒绝所有 国家，允许 巴西、中国 访问
func UpdateFirewallRule(zoneID string) string {

	firewallRuleId := CreateFirewallRuleId(zoneID)

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/rulesets/%s", zoneID, firewallRuleId)
	// 创建一个 HTTP 客户端

	// 创建 JSON  如果是允许所有 拒绝 巴西       (ip.geoip.country in {"BR"})
	jsonData := []byte(`{
        "kind": "zone",
        "description": "This ruleset executes a managed ruleset.",
        "source": "firewall_custom",
        "rules": [
            {
                "action": "block",
                "expression": "(not ip.geoip.country in { \"CN\" \"JP\" })",
                "description": "limit",
                "enabled": true
            }
        ],
        "phase": "http_request_firewall_custom"
    }`)
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// 添加请求头
	req.Header.Add("X-Auth-Email", authEmail)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析 JSON
	var response cloudflare.GetCfResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// 打印响应
	fmt.Println("Response:", string(body))
	return "done"
}

// ----------------------------------------- 开启https跳转
func PatchHttpsOn(zoneID string) string {

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/settings/always_use_https", zoneID)
	// 创建一个 HTTP 客户端

	// 创建 JSON 数据
	jsonData := []byte(`{
			"value": "on"
		}`)
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// 添加请求头
	req.Header.Add("X-Auth-Email", authEmail)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析 JSON
	var response cloudflare.GetCfResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// 打印响应
	fmt.Println("Response:", string(body))
	return "done"
}

func MainCloudflare(domainArray []string) string {

	for i := 0; i < len(domainArray); i++ {
		domainArray[i] = strings.TrimSpace(domainArray[i])
		if domainArray[i] != "" {
			fmt.Printf("域名 #%d: %s\n", i+1, domainArray[i])
		}
	}

	return ""
}

func UpdateSecurityLevel(zoneID string) string {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/settings/security_level", zoneID)

	// JSON data for setting the security level to high
	jsonData := []byte(`{"value": "high"}`)

	client := &http.Client{}

	// Create a new request with the method PATCH
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err) // Handle error according to your error policy
	}

	// Add necessary headers to the request
	req.Header.Add("X-Auth-Email", authEmail)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err) // Handle error according to your error policy
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err) // Handle error according to your error policy
	}

	fmt.Println("Response:", string(body))
	return "done"
}

func UpdateMinTLSVersion(zoneID string) string {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/settings/min_tls_version", zoneID)

	// JSON data for setting the minimum TLS version
	jsonData := []byte(`{"value": "1.2"}`)

	client := &http.Client{}

	// Create a new request with the method PATCH
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err) // Proper error handling is advised here
	}

	// Set necessary headers
	req.Header.Add("X-Auth-Email", authEmail)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err) // Proper error handling is advised here
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err) // Proper error handling is advised here
	}

	fmt.Println("Response:", string(body))
	return "done"
}

// ------------------------------- 网限限速-------------------------------
func CreateRateLimitRuleId(zoneID string) string {

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/rulesets", zoneID)
	// 创建一个 HTTP 客户端

	// 创建 JSON 数据
	jsonData := []byte(`{
		"name": "limit",
		"kind": "zone",
		"description": "This ruleset executes a managed ruleset.",
		"phase": "http_ratelimit"
		}`)
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// 添加请求头
	req.Header.Add("X-Auth-Email", authEmail)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析 JSON
	var response cloudflare.GetCfResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// 打印响应
	fmt.Println("Response: ", string(body))
	fmt.Println("缓存规则ID: ", response.Result.ID)
	return response.Result.ID
}

// --------------------- 添加Waf防火墙规则   配置 速率限制规则 防止cc 攻击。 10秒钟 20 次
func UpdateRateLimitRule(zoneID string) string {

	rateLimitRuleId := CreateRateLimitRuleId(zoneID)

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/rulesets/%s", zoneID, rateLimitRuleId)
	// 创建一个 HTTP 客户端

	// 创建 JSON 数据      
	jsonData := []byte(`{
		"kind": "zone",
		"description": "This ruleset executes a managed ruleset.",
		"source": "rate_limit",
		"rules": [
		  {
			"action": "block",
			"ratelimit": {
				"characteristics": [
				  "ip.src",
				  "cf.colo.id"
				],
				"period": 10,
				"requests_per_period": 20,
				"mitigation_timeout": 10
			  },
			"expression": "(cf.bot_management.verified_bot)",
			"description": "limit",
			"enabled": true
		  }
		],
		"phase": "http_ratelimit"
    }`)
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// 添加请求头
	req.Header.Add("X-Auth-Email", authEmail)
	req.Header.Add("X-Auth-Key", authKey)
	req.Header.Add("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析 JSON
	var response cloudflare.GetCfResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// 打印响应
	fmt.Println("Response:", string(body))
	return "done"
}
