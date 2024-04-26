package cloudflare

type DNSRecord struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl,omitempty"`
	Proxied bool   `json:"proxied,omitempty"`
}

type ZoneResponse struct {
	Result []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"result"`
	Success bool `json:"success"`
}

type GetCfResp struct {
	Result struct {
		ID string `json:"id"`
	} `json:"result"`
	Success bool `json:"success"`
}

type NewSiteResp struct {
	Result struct {
		ID          string   `json:"id"`
		NameServers []string `json:"name_servers"`
	} `json:"result"`
	Success bool `json:"success"`
}
