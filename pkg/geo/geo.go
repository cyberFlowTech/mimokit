package geo

import (
	"github.com/oschwald/geoip2-golang"
	"log"
)

var geoDb *geoip2.Reader

func InitGeoDb(path string) {
	var err error
	geoDb, err = geoip2.Open(path)
	if err != nil {
		log.Fatal(err)
	}
}

func GetGeoDb() *geoip2.Reader {
	return geoDb
}
