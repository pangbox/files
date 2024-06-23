package region

import (
	"log"
	"strings"

	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"github.com/pangbox/pangfiles/pak"
)

var xteaKeys = []pyxtea.Key{
	pyxtea.KeyUS,
	pyxtea.KeyJP,
	pyxtea.KeyTH,
	pyxtea.KeyEU,
	pyxtea.KeyID,
	pyxtea.KeyKR,
}

var regionToKey = map[string]pyxtea.Key{
	"us": pyxtea.KeyUS,
	"jp": pyxtea.KeyJP,
	"th": pyxtea.KeyTH,
	"eu": pyxtea.KeyEU,
	"id": pyxtea.KeyID,
	"kr": pyxtea.KeyKR,
}

var keyToRegion = map[pyxtea.Key]string{
	pyxtea.KeyUS: "us",
	pyxtea.KeyJP: "jp",
	pyxtea.KeyTH: "th",
	pyxtea.KeyEU: "eu",
	pyxtea.KeyID: "id",
	pyxtea.KeyKR: "kr",
}

func Key(regionCode string) pyxtea.Key {
	key, ok := regionToKey[regionCode]
	if !ok {
		log.Fatalf("Invalid region %q (valid regions: us, jp, th, eu, id, kr)", regionCode)
	}
	return key
}

func ForKey(key pyxtea.Key) string {
	region, ok := keyToRegion[key]
	if !ok {
		panic("programming error: unexpected key")
	}
	return region
}

func PakKey(region string, patterns []string) pyxtea.Key {
	if region == "" {
		log.Println("Auto-detecting pak region (use -region to improve startup delay.)")
		key := pak.MustDetectRegion(patterns, xteaKeys)
		log.Printf("Detected pak region as %s.", strings.ToUpper(ForKey(key)))
		return key
	}
	return Key(region)
}
