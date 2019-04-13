// Copyright 2018 - 2019 Christian Müller <dev@c-mueller.xyz>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ads

import (
	"bytes"
	gz "compress/gzip"
	"context"
	"fmt"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"io/ioutil"
	"math/rand"
	"net"
	"strings"
	"testing"
)

const benchmarkSize = 2500000

var blockmap BlockMap

var testCases []test.Case

func init() {
	log.Infof("Initializing Benchmarks")
	log.Info("Loading Blocklist Data")
	blistDataCompressed, _ := ioutil.ReadFile("testdata/benchmark_blocklist.gz")
	log.Info("Extracting Blocklist Data")
	blistDataUncompressed, _ := extractGzip(blistDataCompressed)
	log.Info("Formating Blocklist Data")
	blocklist := strings.Split(strings.ReplaceAll(string(blistDataUncompressed), "0.0.0.0 ", ""), "\n")

	log.Info("Building Blocklist")
	blockmap = make(BlockMap, 0)
	for _, v := range blocklist {
		blockmap[v] = true
	}

	log.Infof("Loaded %d Entries", len(blockmap))

	log.Infof("Generating %d testcases for benchmarking", benchmarkSize)
	testCases = make([]test.Case, 0)
	linedata, _ := ioutil.ReadFile("testdata/benchmark_lookup_set")
	benchLines := strings.Split(string(linedata), "\n")
	for i := 0; i < benchmarkSize; i++ {
		qname := benchLines[i%len(benchLines)]
		typeInt := rand.Intn(3)

		if typeInt == 0 {
			tcase := test.Case{
				Qname: qname, Qtype: dns.TypeAAAA,
				Answer: []dns.RR{
					test.AAAA(fmt.Sprintf("%s. 3600	IN	AAAA fe80::9cbd:c3ff:fe28:e133", qname)),
				},
			}
			testCases = append(testCases, tcase)
		} else {
			tcase := test.Case{
				Qname: qname, Qtype: dns.TypeA,
				Answer: []dns.RR{
					test.A(fmt.Sprintf("%s. 3600	IN	A 10.1.33.7", qname)),
				},
			}
			testCases = append(testCases, tcase)
		}
	}
}

func BenchmarkBlockSpeed(b *testing.B) {
	p := initBenchmarkPlugin(b)
	ctx := context.TODO()

	log.Infof("Trycount %d", b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		testCase := testCases[rand.Intn(benchmarkSize)]

		m := testCase.Msg()

		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		_, err := p.ServeDNS(ctx, rec, m)
		if err != nil {
			b.Errorf("Expected no error, got %v\n", err)
			return
		}
	}
}

func extractGzip(data []byte) ([]byte, error) {
	inputBuffer := bytes.NewReader(data)
	compressionReader, err := gz.NewReader(inputBuffer)
	if err != nil {
		return nil, err
	}

	defer compressionReader.Close()

	return ioutil.ReadAll(compressionReader)
}

func initBenchmarkPlugin(t *testing.B) *DNSAdBlock {
	p := DNSAdBlock{
		Next:       nxDomainHandler(),
		blockMap:   blockmap,
		BlockLists: []string{"http://localhost:8080/mylist.txt"},
		RuleSet:    RuleSet{},
		updater:    nil,
		LogBlocks:  false,
		TargetIP:   net.ParseIP("10.1.33.7"),
		TargetIPv6: net.ParseIP("fe80::9cbd:c3ff:fe28:e133"),
	}

	return &p
}
