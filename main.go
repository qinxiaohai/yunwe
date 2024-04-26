package main

import (
        "fmt"
        "log"
        "os"
        "strings"
        "time"

        implCF "github.com/ned123abc/yunwe/cloudflare/impl"
        "github.com/ned123abc/yunwe/telegram"
)

func main() {
        varNewDomain := os.Getenv("NewDomainName")
        varDomainRecord := os.Getenv("DomainRecord")
        fmt.Printf("域名 : %s\n", varNewDomain)
        fmt.Printf("域名 : %s\n", varDomainRecord)

        if varNewDomain == "" || varDomainRecord == "" {
                log.Fatalf("值不能为空")
        }

        domainArray := strings.Split(varNewDomain, "\n")

        // 清理数组中的每个元素，移除空格
        for i := 0; i < len(domainArray); i++ {
                domainArray[i] = strings.TrimSpace(domainArray[i])
                if domainArray[i] != "" {
                        fmt.Printf("域名 #%d: %s\n", i+1, domainArray[i])

                        NS1, NS2, siteStatus := implCF.NewSite(domainArray[i])
                        if !siteStatus {
                                panic("域名没添加Cf")
                        }

                        fmt.Println("Sleeping for 5 seconds...")
                        time.Sleep(5 * time.Second) // 暂停5秒

                        zoneID, err := implCF.GetZoneId(domainArray[i])
                        if err != nil {
                                panic(err)
                        }

                        if zoneID != "" {
                                implCF.UpdateRecord(domainArray[i], zoneID, strings.TrimSpace(varDomainRecord))
                                time.Sleep(2 * time.Second)
                                implCF.UpdateCacheRule(zoneID)
                                time.Sleep(2 * time.Second)
                                implCF.PatchHttpsOn(zoneID)

                                time.Sleep(2 * time.Second)
                                implCF.UpdateFirewallRule(zoneID)
                                time.Sleep(2 * time.Second)
                                implCF.UpdateSecurityLevel(zoneID)
                                time.Sleep(2 * time.Second)
                                implCF.UpdateMinTLSVersion(zoneID)
                                time.Sleep(2 * time.Second)
                                implCF.UpdateRateLimitRule(zoneID)
                                message := "=== 新域名NS解析 ==="
                                message += "\n域名:" + domainArray[i]
                                message += "\nNS1:" + NS1
                                message += "\nNS2:" + NS2
                                fmt.Println(message)
                                telegram.SendMessage(message)
                        }
                }
        }
}
