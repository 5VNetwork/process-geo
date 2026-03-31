package main

import (
	"log"
	"net/http"
	"os"

	"github.com/5vnetwork/vx-core/app/util"
	"github.com/5vnetwork/vx-core/common/geo"
	"github.com/5vnetwork/vx-core/common/geo/memconservative"
	"github.com/golang/protobuf/proto"
)

var geoIpUrl = "https://cdn.jsdelivr.net/gh/Loyalsoldier/v2ray-rules-dat@release/geoip.dat"
var geoSiteUrl = "https://cdn.jsdelivr.net/gh/Loyalsoldier/v2ray-rules-dat@release/geosite.dat"

func main() {
	cn()
	ru()
}

func cn() {
	var err error
	err = util.DownloadToFile(geoIpUrl, http.DefaultClient, "geoip.dat")
	if err != nil {
		log.Fatal(err)
	}
	err =
		util.DownloadToFile(geoSiteUrl, http.DefaultClient, "geosite.dat")
	if err != nil {
		log.Fatal(err)
	}

	geocode := []string{"cn", "apple-cn", "google-cn", "google", "tld-cn", "private", "category-games", "gfw"}
	geoipcode := []string{"private", "cn", "telegram", "google", "facebook", "twitter", "tor"}
	geositePath := "geosite.dat"
	geoipPath := "geoip.dat"
	dstGeositePath := "simplified_geosite.dat"
	dstGeoipPath := "simplified_geoip.dat"

	l := memconservative.NewMemConservativeLoader()
	geositeList := &geo.GeoSiteList{
		Entry: []*geo.GeoSite{},
	}
	geoIpList := &geo.GeoIPList{
		Entry: []*geo.GeoIP{},
	}
	for _, code := range geocode {
		site, err := l.LoadSite(geositePath, code)
		if err != nil {
			log.Fatal(err)
		}
		geositeList.Entry = append(geositeList.Entry, site)
	}
	for _, code := range geoipcode {
		cidr, err := l.LoadIP(geoipPath, code)
		if err != nil {
			log.Fatal(err)
		}
		geoIpList.Entry = append(geoIpList.Entry, cidr)
	}

	// write into files
	// Marshal the geo data to protobuf format
	geositeBytes, err := proto.Marshal(geositeList)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(dstGeositePath, geositeBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

	geoipBytes, err := proto.Marshal(geoIpList)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(dstGeoipPath, geoipBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

const russiaGeoSiteUrl = "https://raw.githubusercontent.com/runetfreedom/russia-v2ray-rules-dat/release/geosite.dat"
const russiaGeoIpUrl = "https://raw.githubusercontent.com/runetfreedom/russia-v2ray-rules-dat/release/geoip.dat"

func ru() {
	err := util.DownloadToFile(russiaGeoSiteUrl, http.DefaultClient, "geosite_ru.dat")
	if err != nil {
		log.Fatal(err)
	}

	geositeList := &geo.GeoSiteList{
		Entry: []*geo.GeoSite{},
	}

	l := memconservative.NewMemConservativeLoader()
	ruSite, err := l.LoadSite("geosite_ru.dat", "ru-blocked")
	if err != nil {
		log.Fatal(err)
	}
	geositeList.Entry = append(geositeList.Entry, ruSite)

	// ruInside, err := l.LoadSite("geosite_ru.dat", "ru-available-only-inside")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// geositeList.Entry = append(geositeList.Entry, ruInside)

	geositeBytes, err := proto.Marshal(geositeList)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("simplified_geosite_ru.dat", geositeBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
