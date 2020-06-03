package logic

import (
	"encoding/json"
	"flag"
	"fmt"
	"hostingDeUpdateIP/models"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetExternalIPv6(ipv6postfix string) []string {
	var ipv6adresses []string
	ifaces, err := net.Interfaces()

	log.Debug("Searching for IPv6 with prefix: ", ipv6postfix)

	if err != nil {
		log.Error(err)
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Error(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			var foundIp string = ip.String()
			if strings.Contains(foundIp, ipv6postfix) {
				log.Debug("Added ipv6: " + foundIp)
				ipv6adresses = append(ipv6adresses, foundIp)
			}
		}
	}
	return ipv6adresses
}

func GetExternalIpv4(ipv4serviceurl string) string {
	resp, err := http.Get(ipv4serviceurl)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(content)
}

func FindActualIpv6Adress(zones []models.Zone, ipv4address string, ipv6adresses []string) []models.Zone {

	for i, zone := range zones {
		for j, hostEntry := range zone.Hostentries {
			log.Debug("Hostentry: ", hostEntry)
			// Set IPv4
			zones[i].Hostentries[j].A = ipv4address

			// Set IPv6
			var AAAA string = zones[i].Hostentries[j].AAAAOld
			for _, ipv6address := range ipv6adresses {
				log.Debug("AAAA | ipv6address: ", AAAA, "|", ipv6address)
				if AAAA != ipv6address {
					log.Debug("ipv6 differs: ", ipv6address)
					zones[i].Hostentries[j].AAAA = ipv6address
					break
				}
			}
		}
	}
	return zones
}

func LoadConfig(fileName string) models.Configuration {

	config := DefaultConfigBuilder()

	file, err := os.Open(fileName)
	if err != nil {
		log.Error(err)
		log.Error("Can't open configuration")
		os.Exit(1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)
	if err != nil {
		log.Error(err)
		log.Error("Can't load configuration %s", file.Name())
		os.Exit(1)
	}

	var logLevel log.Level = log.DebugLevel
	switch config.LogLevel {
	case "Error":
		logLevel = log.ErrorLevel
	case "Info":
		logLevel = log.InfoLevel
	case "Debug":
		logLevel = log.DebugLevel
	default:
		logLevel = log.DebugLevel
		log.Error("Loglevel '", config.LogLevel, "' is wrong. Switching to Debug mode")
	}
	log.SetLevel(logLevel)

	log.Debug("Configuration: ", config)
	return config
}

func DefaultConfigBuilder() models.Configuration {

	config := models.Configuration{}

	config.HostingDe.Api.AuthToken = "NOT-KEY-SET"
	config.HostingDe.Api.Url = "https://secure.hosting.de/api/dns/v1/json/"
	config.HostingDe.Api.FindZones = "zonesFind"
	config.HostingDe.Api.FindRecords = "recordsFind"
	config.HostingDe.Api.UpdateZone = "zoneUpdate"

	var domain1 models.Domains = models.Domains{
		Host:       "domain1.tld",
		SetHostToo: false,
		Subs:       []string{"www", "www2"},
	}
	var domain2 models.Domains = models.Domains{
		Host:       "domain2.tld",
		SetHostToo: false,
		Subs:       []string{"blog"},
	}
	config.Domains = append(config.Domains, domain1)
	config.Domains = append(config.Domains, domain2)

	config.ZoneUpdateDefault.ZoneConfig.Type = "NATIVE"
	config.ZoneUpdateDefault.ZoneConfig.EmailAddress = "admin@domain.tld"
	config.ZoneUpdateDefault.ZoneConfig.DnsSecMode = "off"
	config.ZoneUpdateDefault.ZoneConfig.SoaValues.Refresh = 86400
	config.ZoneUpdateDefault.ZoneConfig.SoaValues.Retry = 7200
	config.ZoneUpdateDefault.ZoneConfig.SoaValues.Expire = 3600000
	config.ZoneUpdateDefault.ZoneConfig.SoaValues.Ttl = 172800
	config.ZoneUpdateDefault.ZoneConfig.SoaValues.NegativeTtl = 3600
	config.Ipv4ServiceUrl = "http://v4.ipv6-test.com/api/myip.php"
	config.ZoneUpdateDefault.RecordToAdd.Ttl = 86000
	config.Ipv6Postfix = "::1234:5678"
	config.Ipv4ServiceUrl = "http://v4.ipv6-test.com/api/myip.php"
	config.LogLevel = "Debug"
	config.Simulate = false

	log.Debug("Default Configuration: ", config)
	return config
}

func ConfigLooksGood(config models.Configuration) bool {
	var errorlist []string = []string{}
	if config.Ipv6Postfix == "" {
		errorlist = append(errorlist, "Ipv6Postfix is empty")
	}

	if len(errorlist) > 0 {
		log.Error("Error list is not empty: ", errorlist)
		return false
	}

	return true
}

func ReadArgs() models.ProgramFlags {
	var programFlags models.ProgramFlags = models.ProgramFlags{}

	// flags declaration using flag package
	flag.StringVar(&programFlags.ConfigFile, "cf", "config.json", "ConfigFile")
	flag.BoolVar(&programFlags.ShowHelp, "h", false, "Show Help")
	flag.BoolVar(&programFlags.ShowExampleConfig, "ec", false, "Show Example Config")
	flag.BoolVar(&programFlags.ShowVersion, "v", false, "Show Version")

	flag.Parse()

	return programFlags
}

func ProcessArgs(programFlags models.ProgramFlags, versionString string) {

	if programFlags.ShowHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if programFlags.ShowExampleConfig {
		var config models.Configuration = DefaultConfigBuilder()
		jsonStr, err := json.MarshalIndent(config, "", "   ")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		fmt.Println(string(jsonStr))
		os.Exit(0)
	}
	if programFlags.ShowVersion {
		log.Info(versionString)
		os.Exit(0)
	}

}
