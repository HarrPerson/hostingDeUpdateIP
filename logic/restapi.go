package logic

import (
	"bytes"
	"encoding/json"
	"hostingDeUpdateIP/models"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func getZonesFromConfig(config models.Configuration) []models.Zone {
	var zones []models.Zone = []models.Zone{}
	var zone models.Zone = models.Zone{}
	var hostEntry models.HostEntry = models.HostEntry{}

	for _, domain := range config.Domains {
		zone = models.Zone{Name: domain.Host}
		zone.Hostentries = []models.HostEntry{}
		if domain.SetHostToo == true {
			hostEntry = models.HostEntry{
				Domain: domain.Host,
				Ttl:    config.ZoneUpdateDefault.RecordToAdd.Ttl}
			zone.Hostentries = append(zone.Hostentries, hostEntry)
		}
		for _, sub := range domain.Subs {

			hostEntry = models.HostEntry{
				Domain: sub + "." + domain.Host,
				Ttl:    config.ZoneUpdateDefault.RecordToAdd.Ttl}
			zone.Hostentries = append(zone.Hostentries, hostEntry)
		}
		zones = append(zones, zone)
	}
	return zones
}

func FindHostEntries(config models.Configuration) []models.Zone {
	log.Debug("Find host entries and old ips")
	var zones []models.Zone = getZonesFromConfig(config)

	log.Debug("Zones: ", zones)

	for i, zone := range zones {
		log.Info("Fetching Zone: ", zone.Name)
		for j, hostEntry := range zone.Hostentries {
			log.Info("Fetching Domain: ", hostEntry.Domain)
			zones[i].Hostentries[j] = findRecordsApiCall(config, hostEntry)
			log.Debug("Hostentry: ", hostEntry)
		}
	}

	log.Debug("Zones: ", zones)
	return zones
}

func findRecordsApiCall(config models.Configuration, hostEntry models.HostEntry) models.HostEntry {

	log.Debug("Call recordsFind for hostEntry ", hostEntry)

	var url = config.HostingDe.Api.Url + config.HostingDe.Api.FindRecords
	var findRecord = models.FindRecords{}
	findRecord.AuthToken = config.HostingDe.Api.AuthToken
	findRecord.Filter.Field = "RecordName"
	findRecord.Filter.Value = hostEntry.Domain
	log.Debug("Found record: ", findRecord)

	jsonStr, err := json.Marshal(findRecord)
	if err != nil {
		log.Panic("Json Marshal failed", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	var client *http.Client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Debug("response Headers: ", resp.Header)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panic(err)
		}
		log.Debug("response Body: ", string(body))

		var findResponse models.ZoneUpdateRecordResponse = models.ZoneUpdateRecordResponse{}

		err = json.Unmarshal([]byte(body), &findResponse)
		if err != nil {
			log.Panic(err)
		}
		log.Debug("Response Data: ", findResponse.Response.Data)

		if findResponse.Status != "success" && findResponse.Status != "pending" {
			log.Error("Return was not successfull")
			log.Error("response Body: ", string(body))
			os.Exit(1)
		}

		for _, d := range findResponse.Response.Data {
			switch strings.ToLower(d.Type) {
			case "a":
				hostEntry.Domain = d.Name
				hostEntry.AOld = d.Content
			case "aaaa":
				hostEntry.Domain = d.Name
				hostEntry.AAAAOld = d.Content
			}
		}
		return hostEntry
	}

	log.Panic("Something went wrong!")
	return hostEntry

}

func UpdateHostEntries(config models.Configuration, zoneUpdate models.ZoneUpdate) {

	var url = config.HostingDe.Api.Url + config.HostingDe.Api.UpdateZone

	if len(zoneUpdate.RecordsToAdd) > 0 {
		log.Info("Updating Entries: ", zoneUpdate.RecordsToAdd)
		jsonStr, err := json.Marshal(zoneUpdate)
		if err != nil {
			log.Panic("Json Marshal failed", err)
		}
		log.Debug("Request JSON: ", string(jsonStr))

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		var client *http.Client = &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			log.Debug("response Headers: ", resp.Header)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Panic(err)
			}
			log.Debug("response Body: ", string(body))

			var findResponse models.ZoneUpdateRecordResponse = models.ZoneUpdateRecordResponse{}

			err = json.Unmarshal([]byte(body), &findResponse)
			if err != nil {
				log.Panic(err)
			}
			log.Debug("Response Data: ", findResponse.Response.Data)

			if findResponse.Status != "success" && findResponse.Status != "pending" {
				log.Error("Return was not successfull for ", zoneUpdate.ZoneConfig.Name)
				log.Error("response Body: ", string(body))
				os.Exit(1)
			} else {
				log.Info("Update finished for Zone: ", zoneUpdate.ZoneConfig.Name)
				os.Exit(0)
			}
		}
	} else {
		log.Info("No Update necessary for ", zoneUpdate.ZoneConfig.Name, ", recordsToAdd is empty. Means IPs didn't change.")
	}
}

func BuildZoneUpdate(zoneUpdate models.ZoneUpdate, zone models.Zone) models.ZoneUpdate {

	var hostEntries []models.HostEntry = zone.Hostentries

	for _, hostEntry := range hostEntries {
		var recordToAdd models.RecordsToAdd
		var recordToDelete models.RecordsToDelete

		// Fill recordsToAdd
		if hostEntry.A != "" && hostEntry.A != hostEntry.AOld {
			recordToAdd = models.RecordsToAdd{
				Name:    hostEntry.Domain,
				Type:    "A",
				Content: hostEntry.A,
				Ttl:     hostEntry.Ttl,
			}
			zoneUpdate.RecordsToAdd = append(zoneUpdate.RecordsToAdd, recordToAdd)

		}
		if hostEntry.AAAA != "" && hostEntry.AAAA != hostEntry.AAAAOld {
			recordToAdd = models.RecordsToAdd{
				Name:    hostEntry.Domain,
				Type:    "AAAA",
				Content: hostEntry.AAAA,
				Ttl:     hostEntry.Ttl,
			}
			zoneUpdate.RecordsToAdd = append(zoneUpdate.RecordsToAdd, recordToAdd)

		}

		// Fill recordsToDelete
		if hostEntry.AAAAOld != "" && hostEntry.AOld != "" {
			if hostEntry.A != "" && hostEntry.A != hostEntry.AOld {
				recordToDelete = models.RecordsToDelete{
					Name:    hostEntry.Domain,
					Type:    "A",
					Content: hostEntry.AOld,
				}
				zoneUpdate.RecordsToDelete = append(zoneUpdate.RecordsToDelete, recordToDelete)

			}
			log.Debug("Domain AAAA AAAAold: ", hostEntry.Domain, " ", hostEntry.AAAA, " ", hostEntry.AAAAOld)
			if hostEntry.AAAA != "" && hostEntry.AAAA != hostEntry.AAAAOld {
				recordToDelete = models.RecordsToDelete{
					Name:    hostEntry.Domain,
					Type:    "AAAA",
					Content: hostEntry.AAAAOld,
				}
				zoneUpdate.RecordsToDelete = append(zoneUpdate.RecordsToDelete, recordToDelete)
			}
		}
	}

	log.Debug("ZoneUpdate: ", zoneUpdate)
	return zoneUpdate
}
