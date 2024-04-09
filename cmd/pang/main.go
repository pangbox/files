package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/subcommands"
	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"github.com/pangbox/pangfiles/pak"
	"github.com/pangbox/pangfiles/version"
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

func getRegionKey(regionCode string) pyxtea.Key {
	key, ok := regionToKey[regionCode]
	if !ok {
		log.Fatalf("Invalid region %q (valid regions: us, jp, th, eu, id, kr)", regionCode)
	}
	return key
}

func getKeyRegion(key pyxtea.Key) string {
	region, ok := keyToRegion[key]
	if !ok {
		panic("programming error: unexpected key")
	}
	return region
}

func getPakKey(region string, patterns []string) pyxtea.Key {
	if region == "" {
		log.Println("Auto-detecting pak region (use -region to improve startup delay.)")
		key := pak.MustDetectRegion(patterns, xteaKeys)
		log.Printf("Detected pak region as %s.", strings.ToUpper(getKeyRegion(key)))
		return key
	}
	return getRegionKey(region)
}

type versionCmd struct{}

func (versionCmd) Name() string           { return "version" }
func (versionCmd) Synopsis() string       { return "Prints version information to stdout." }
func (v versionCmd) Usage() string        { return fmt.Sprintf("%s:\n  %s\n", v.Name(), v.Synopsis()) }
func (versionCmd) SetFlags(*flag.FlagSet) {}
func (versionCmd) Execute(context.Context, *flag.FlagSet, ...interface{}) subcommands.ExitStatus {
	versionStr := "v" + version.Release
	if version.GitCommit != "" {
		versionStr += "+" + version.GitCommit
	}
	fmt.Println(versionStr)
	return subcommands.ExitSuccess
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&versionCmd{}, "")
	subcommands.Register(&cmdPakMount{}, "paks")
	subcommands.Register(&cmdPakExtract{}, "paks")
	subcommands.Register(&cmdUpdateListServe{}, "updatelists")
	subcommands.Register(&cmdUpdateListEncrypt{}, "updatelists")
	subcommands.Register(&cmdUpdateListDecrypt{}, "updatelists")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
