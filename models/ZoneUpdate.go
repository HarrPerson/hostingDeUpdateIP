package models

type ZoneUpdate struct {
	AuthToken  string `json:"authToken"`
	ZoneConfig struct {
		Type         string `json:"type"`
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		DnsSecMode   string `json:"dnsSecMode"`
		SoaValues    struct {
			Refresh     int `json:"refresh"`
			Retry       int `json:"retry"`
			Expire      int `json:"expire"`
			Ttl         int `json:"ttl"`
			NegativeTtl int `json:"negativeTtl"`
		} `json:"soaValues"`
	} `json:"zoneConfig"`
	RecordsToAdd    []RecordsToAdd    `json:"recordsToAdd"`
	RecordsToDelete []RecordsToDelete `json:"recordsToDelete"`
}

type RecordsToAdd struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
}

type RecordsToDelete struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type ZoneUpdateRecordResponse struct {
	Response struct {
		Data []struct {
			Name    string `json:"name"`
			Content string `json:"content"`
			Type    string `json:"type"`
		} `json:"data"`
	} `json:"response"`
	Status string
}

// Builder
func ZoneUpdateBuilder(config Configuration, zone Zone) ZoneUpdate {
	zoneUpdate := ZoneUpdate{}

	zoneUpdate.AuthToken = config.HostingDe.Api.AuthToken
	zoneUpdate.ZoneConfig.Type = config.ZoneUpdateDefault.ZoneConfig.Type
	zoneUpdate.ZoneConfig.Name = zone.Name
	zoneUpdate.ZoneConfig.EmailAddress = config.ZoneUpdateDefault.ZoneConfig.EmailAddress
	zoneUpdate.ZoneConfig.DnsSecMode = config.ZoneUpdateDefault.ZoneConfig.DnsSecMode
	zoneUpdate.ZoneConfig.SoaValues.Refresh = config.ZoneUpdateDefault.ZoneConfig.SoaValues.Refresh
	zoneUpdate.ZoneConfig.SoaValues.Retry = config.ZoneUpdateDefault.ZoneConfig.SoaValues.Retry
	zoneUpdate.ZoneConfig.SoaValues.Expire = config.ZoneUpdateDefault.ZoneConfig.SoaValues.Expire
	zoneUpdate.ZoneConfig.SoaValues.Ttl = config.ZoneUpdateDefault.ZoneConfig.SoaValues.Ttl
	zoneUpdate.ZoneConfig.SoaValues.NegativeTtl = config.ZoneUpdateDefault.ZoneConfig.SoaValues.NegativeTtl

	zoneUpdate.RecordsToAdd = []RecordsToAdd{}
	zoneUpdate.RecordsToDelete = []RecordsToDelete{}

	return zoneUpdate
}
