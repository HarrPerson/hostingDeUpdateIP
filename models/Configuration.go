package models

type Configuration struct {
	HostingDe struct {
		Api struct {
			AuthToken   string `json:"AuthToken"`
			Url         string `json:"Url"`
			FindZones   string `json:"FindZones"`
			FindRecords string `json:"FindRecords"`
			UpdateZone  string `json:"UpdateZone"`
		} `json:"Api"`
	} `json:"HostingDe"`
	Domains           []Domains `json:"Domains"`
	ZoneUpdateDefault struct {
		ZoneConfig struct {
			Type         string `json:"Type"`
			Name         string `json:"Name"`
			EmailAddress string `json:"EmailAddress"`
			DnsSecMode   string `json:"DnsSecMode"`
			SoaValues    struct {
				Refresh     int `json:"Refresh"`
				Retry       int `json:"Retry"`
				Expire      int `json:"Expire"`
				Ttl         int `json:"Ttl"`
				NegativeTtl int `json:"NegativeTtl"`
			} `json:"SoaValues"`
		} `json:"ZoneConfig"`
		RecordToAdd struct {
			Ttl int `json:"Ttl"`
		} `json:"RecordToAdd"`
	} `json:"ZoneUpdateDefault"`
	Ipv6Postfix    string `json:"Ipv6Postfix"`
	Ipv4ServiceUrl string `json:"Ipv4ServiceUrl"`
	LogLevel       string `json:"LogLevel"`
	Simulate       bool   `json:"Simulate"`
}

type Domains struct {
	Host       string   `json:"Host"`
	SetHostToo bool     `json:"SetHostToo"`
	Subs       []string `json:"Subs"`
}
