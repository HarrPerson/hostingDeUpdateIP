package simulation

import (
	"hostingDeUpdateIP/models"
)

func DemoValues(config models.Configuration) ([]models.Zone, []string, string) {

	var hostEntries = []models.HostEntry{
		{
			"existing.domain.dev",
			"127.0.0.1",
			"127.0.0.1",
			"2001:abcd:aa:2e03::3:50",
			"fe80::0",
			config.ZoneUpdateDefault.RecordToAdd.Ttl,
		},
		{
			"new.domain.dev",
			"127.0.0.1",
			"",
			"2001:abcd:aa:2e03::3:50",
			"",
			config.ZoneUpdateDefault.RecordToAdd.Ttl,
		},
		{
			"nochange.domain.dev",
			"127.0.0.1",
			"127.0.0.1",
			"2001:abcd:aa:2e03::3:50",
			"2001:abcd:aa:2e03::3:50",
			config.ZoneUpdateDefault.RecordToAdd.Ttl,
		},
	}

	var ipv6adresses = []string{"2001:abcd:aa:2e03::3:50", "2001:abcd:aa:1234::3:50"}
	var ipv4address string = "127.0.0.1"

	var zones []models.Zone = []models.Zone{models.Zone{Name: "domain.dev", Hostentries: hostEntries}}
	return zones, ipv6adresses, ipv4address
}
