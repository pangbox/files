package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"github.com/pangbox/pangfiles/encoding/litexml"
	"github.com/pangbox/pangfiles/updatelist"
	"github.com/pangbox/pangfiles/util"
)

type cacheentry struct {
	modTime time.Time
	fSize   int64
	fInfo   updatelist.FileInfo
}

type server struct {
	key   pyxtea.Key
	dir   string
	cache map[string]cacheentry
	mutex sync.RWMutex
}

func (s *server) calcEntry(wg *sync.WaitGroup, entry *updatelist.FileInfo, f os.FileInfo) {
	defer wg.Done()
	var err error

	name := f.Name()
	*entry, err = updatelist.MakeFileInfo(s.dir, "", f, f.Size())

	if err != nil {
		log.Printf("Error calculating entry for %s: %s", name, err)
		entry.Filename = name
	} else {
		log.Printf("Successfully calculated entry for %s", name)

		s.mutex.Lock()
		defer s.mutex.Unlock()

		s.cache[name] = cacheentry{
			modTime: f.ModTime(),
			fSize:   f.Size(),
			fInfo:   *entry,
		}
	}
}

func (s *server) updateList(rw io.Writer) {
	start := time.Now()

	files, err := ioutil.ReadDir(s.dir)
	if err != nil {
		panic(err)
	}

	doc := updatelist.Document{}
	doc.Info.Version = "1.0"
	doc.Info.Encoding = "euc-kr"
	doc.Info.Standalone = "yes"
	doc.PatchVer = "FakeVer"
	doc.PatchNum = 9999
	doc.UpdateListVer = "20090331"

	hit, miss := 0, 0

	var wg sync.WaitGroup
	doc.UpdateFiles.Files = make([]updatelist.FileInfo, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()

		s.mutex.RLock()
		cache, ok := s.cache[name]
		s.mutex.RUnlock()

		if ok && cache.modTime == f.ModTime() && cache.fSize == f.Size() {
			// Cache hit
			hit++
			doc.UpdateFiles.Files = append(doc.UpdateFiles.Files, cache.fInfo)
			doc.UpdateFiles.Count++
		} else {
			// Cache miss, calculate concurrently.
			miss++
			doc.UpdateFiles.Files = append(doc.UpdateFiles.Files, updatelist.FileInfo{})
			doc.UpdateFiles.Count++
			entry := &doc.UpdateFiles.Files[len(doc.UpdateFiles.Files)-1]
			wg.Add(1)
			go s.calcEntry(&wg, entry, f)
		}
	}

	wg.Wait()

	data, err := litexml.Marshal(doc)
	if err != nil {
		panic(err)
	}

	pyxtea.EncipherStream(s.key, util.NullInputPadder{Reader: bytes.NewReader(data)}, rw)

	log.Printf("Updatelist served in %s (cache hits: %d, misses: %d)", time.Since(start), hit, miss)
}

func (s *server) extracontents(w io.Writer) {
	w.Write([]byte{
		0x3c, 0x3f, 0x78, 0x6d, 0x6c, 0x20, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
		0x6e, 0x3d, 0x22, 0x31, 0x2e, 0x30, 0x22, 0x20, 0x73, 0x74, 0x61, 0x6e,
		0x64, 0x61, 0x6c, 0x6f, 0x6e, 0x65, 0x3d, 0x22, 0x79, 0x65, 0x73, 0x22,
		0x20, 0x3f, 0x3e, 0x0a, 0x0a, 0x3c, 0x65, 0x78, 0x74, 0x72, 0x61, 0x63,
		0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x3e, 0x0a, 0x0a, 0x09, 0x3c,
		0x74, 0x68, 0x65, 0x6d, 0x65, 0x73, 0x3e, 0x0a, 0x09, 0x09, 0x3c, 0x70,
		0x61, 0x6e, 0x67, 0x79, 0x61, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c,
		0x74, 0x20, 0x73, 0x72, 0x63, 0x3d, 0x22, 0x70, 0x61, 0x6e, 0x67, 0x79,
		0x61, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x2e, 0x78, 0x6d,
		0x6c, 0x22, 0x20, 0x75, 0x72, 0x6c, 0x3d, 0x22, 0x68, 0x74, 0x74, 0x70,
		0x3a, 0x2f, 0x2f, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31,
		0x3a, 0x38, 0x30, 0x38, 0x30, 0x2f, 0x53, 0x34, 0x5f, 0x50, 0x61, 0x74,
		0x63, 0x68, 0x2f, 0x65, 0x78, 0x74, 0x72, 0x61, 0x63, 0x6f, 0x6e, 0x74,
		0x65, 0x6e, 0x74, 0x73, 0x2f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
		0x2f, 0x22, 0x2f, 0x3e, 0x0a, 0x09, 0x09, 0x3c, 0x21, 0x2d, 0x2d, 0x20,
		0x3c, 0x70, 0x61, 0x6e, 0x67, 0x79, 0x61, 0x5f, 0x64, 0x65, 0x66, 0x61,
		0x75, 0x6c, 0x74, 0x20, 0x73, 0x72, 0x63, 0x3d, 0x22, 0x70, 0x61, 0x6e,
		0x67, 0x79, 0x61, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x2e,
		0x78, 0x6d, 0x6c, 0x22, 0x20, 0x75, 0x72, 0x6c, 0x3d, 0x22, 0x68, 0x74,
		0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x73, 0x75, 0x70, 0x65, 0x72, 0x73, 0x73,
		0x2e, 0x73, 0x79, 0x74, 0x65, 0x73, 0x2e, 0x6e, 0x65, 0x74, 0x3a, 0x38,
		0x30, 0x38, 0x30, 0x2f, 0x53, 0x34, 0x5f, 0x50, 0x61, 0x74, 0x63, 0x68,
		0x2f, 0x65, 0x78, 0x74, 0x72, 0x61, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
		0x74, 0x73, 0x2f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x22,
		0x2f, 0x3e, 0x09, 0x20, 0x2d, 0x2d, 0x3e, 0x09, 0x0a, 0x09, 0x3c, 0x2f,
		0x74, 0x68, 0x65, 0x6d, 0x65, 0x73, 0x3e, 0x0a, 0x09, 0x0a, 0x3c, 0x2f,
		0x65, 0x78, 0x74, 0x72, 0x61, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
		0x73, 0x3e, 0x0a,
	})
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	log.Printf("%s %s", r.Method, r.URL)
	if strings.Contains(strings.ToLower(r.URL.Path), "updatelist") {
		s.updateList(w)
	} else if strings.Contains(strings.ToLower(r.URL.Path), "extracontents") {
		s.extracontents(w)
	}
}

func main() {
	var key pyxtea.Key
	region := flag.String("region", "us", "Region to use (us, jp, th, eu, id, kr)")
	listen := flag.String("listen", ":8080", "Address to listen on.")
	flag.Parse()

	switch *region {
	case "us":
		key = pyxtea.KeyUS
	case "jp":
		key = pyxtea.KeyJP
	case "th":
		key = pyxtea.KeyTH
	case "eu":
		key = pyxtea.KeyEU
	case "id":
		key = pyxtea.KeyID
	case "kr":
		key = pyxtea.KeyKR
	default:
		log.Fatalf("invalid region %q (valid regions: us, jp, th, eu, id, kr)", *region)
	}

	if flag.NArg() < 1 {
		log.Fatalln("Please provide a command. (valid commands: serve)")
	}

	switch flag.Arg(0) {
	case "serve":
		if flag.NArg() < 2 {
			log.Fatalln("Serve requires 1 argument (path to game folder)")
		}
		s := server{
			key:   key,
			dir:   flag.Arg(1),
			cache: map[string]cacheentry{},
		}
		http.ListenAndServe(*listen, &s)
	}
}
