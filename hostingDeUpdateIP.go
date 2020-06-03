package main

import (
	"encoding/json"
	"fmt"
	"hostingDeUpdateIP/logic"
	"hostingDeUpdateIP/models"
	"hostingDeUpdateIP/simulation"
	"os"

	log "github.com/sirupsen/logrus"
)

const version = "0.5.1"

var buildTime = "<buildTime not set>"

var versionString = "Version " + version + " | Build: " + buildTime

func main() {

	var programFlags models.ProgramFlags = logic.ReadArgs()

	logic.ProcessArgs(programFlags, versionString)

	var config models.Configuration = logic.LoadConfig(programFlags.ConfigFile)

	if !logic.ConfigLooksGood(config) {
		os.Exit(1)
	}

	var zones []models.Zone
	var ipv6adresses []string
	var ipv4address string

	if !config.Simulate {
		ipv6adresses = logic.GetExternalIPv6(config.Ipv6Postfix)
		ipv4address = logic.GetExternalIpv4(config.Ipv4ServiceUrl)

		log.Debug("IPv6s | IPv4: ", ipv6adresses, " | ", ipv4address)

		zones = logic.FindHostEntries(config)
		log.Debug("Original HostingEntries: ", zones)

		zones = logic.FindActualIpv6Adress(zones, ipv4address, ipv6adresses)
		log.Debug("Modified HostingEntries: ", zones)
	} else {
		log.Warn("Simulating only!")
		log.Debug("Original HostingEntries: ", zones)
		zones, ipv6adresses, ipv4address = simulation.DemoValues(config)
		log.Debug("Modified HostingEntries: ", zones)
		log.Debug("IPv6s | IPv4: ", ipv6adresses, " | ", ipv4address)
	}

	log.Debug("Len: ", len(zones))
	log.Debug("Zone: ", zones)
	for _, zone := range zones {
		zoneUpdate := models.ZoneUpdateBuilder(config, zone)
		log.Debug(zoneUpdate)

		zoneUpdate = logic.BuildZoneUpdate(zoneUpdate, zone)

		if config.LogLevel == "Debug" {
			jsonStr, _ := json.MarshalIndent(zoneUpdate, "", "    ")
			fmt.Println(string(jsonStr))
		}

		if !config.Simulate {
			logic.UpdateHostEntries(config, zoneUpdate)
		}
	}
	os.Exit(0)
}
