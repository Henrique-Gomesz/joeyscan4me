package runner

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Henrique-Gomesz/JoeyScan4Me/pkg/logging"

	"github.com/projectdiscovery/goflags"
)

type Options struct {
	Domain            string
	Workdir           string
	Server            bool
	Profile           string
	SummaryJSON       bool
	SummaryFile       string
	SubfinderThreads  int
	SubfinderTimeout  int
	SubfinderMaxTime  int
	HttpxThreads      int
	HttpxRateLimit    int
	HttpxTimeout      int
	HttpxPorts        string
	KatanaDepth       int
	KatanaTimeout     int
	KatanaConcurrency int
	KatanaParallelism int
	KatanaRateLimit   int
}

type profilePreset struct {
	SubfinderThreads  int
	SubfinderTimeout  int
	SubfinderMaxTime  int
	HttpxThreads      int
	HttpxRateLimit    int
	HttpxTimeout      int
	HttpxPorts        string
	KatanaDepth       int
	KatanaTimeout     int
	KatanaConcurrency int
	KatanaParallelism int
	KatanaRateLimit   int
}

const bugBountyDefaultPorts = "66,80,81,82,83,84,85,86,87,88,89,90,280,300,443,444,445,457,591,593,832,981,1010,1080,1099,1100,1220,1234,1241,1311,1313,1337,1352,1433,1434,1443,1521,1533,1581,1719,1720,1723,1755,1830,1900,1944,2052,2053,2082,2083,2086,2087,2095,2096,2181,2222,2301,2375,2376,2480,2525,2718,3000,3001,3002,3003,3004,3005,3006,3007,3008,3009,3010,3011,3030,3050,3070,3111,3128,3168,3200,3260,3269,3300,3306,3333,3380,3389,3401,3443,3500,3690,3700,3780,3785,3790,4000,4001,4002,4003,4004,4005,4006,4040,4044,4063,4080,4100,4194,4200,4201,4242,4243,4280,4300,4321,4369,4443,4444,4445,4500,4567,4568,4680,4711,4712,4747,4848,4869,4993,5000,5001,5002,5003,5004,5005,5006,5007,5008,5009,5010,5050,5060,5061,5080,5087,5100,5101,5104,5108,5173,5190,5200,5222,5269,5280,5281,5357,5400,5432,5443,5500,5555,5556,5600,5601,5602,5671,5672,5678,5800,5801,5802,5900,5984,5985,5986,6000,6001,6002,6003,6060,6061,6080,6082,6083,6084,6085,6086,6161,6180,6188,6200,6346,6347,6379,6443,6543,6789,6888,7000,7001,7002,7070,7071,7080,7081,7396,7443,7474,7496,7547,7700,7777,7778,7779,7800,8000,8001,8002,8003,8004,8005,8006,8007,8008,8009,8010,8011,8012,8013,8014,8015,8020,8028,8030,8040,8042,8060,8069,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8095,8099,8100,8118,8123,8139,8140,8161,8172,8180,8181,8182,8200,8222,8243,8280,8281,8290,8300,8333,8384,8400,8401,8403,8443,8444,8445,8480,8484,8500,8529,8530,8531,8545,8554,8555,8580,8585,8590,8591,8600,8610,8649,8686,8718,8765,8766,8777,8787,8800,8812,8834,8838,8880,8881,8882,8883,8887,8888,8889,8899,8983,8989,8990,9000,9001,9002,9003,9004,9009,9010,9011,9012,9014,9015,9020,9021,9022,9023,9024,9025,9026,9027,9028,9029,9030,9043,9050,9051,9060,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9095,9096,9097,9098,9099,9100,9111,9115,9151,9191,9200,9217,9229,9295,9296,9300,9333,9392,9418,9443,9444,9445,9485,9500,9502,9503,9555,9580,9600,9800,9860,9865,9870,9871,9875,9876,9877,9943,9944,9981,9988,9990,9991,9992,9993,9997,9998,9999,10000,10001,10002,10003,10004,10005,10008,10009,10010,10025,10080,10081,10082,10083,10088,10100,10180,10200,10243,10250,10255,10443,10554,10616,10617,10621,11000,11001,11111,11211,11234,11333,11371,11500,12000,12043,12046,12201,12222,12345,12443,13443,14000,14147,15000,15671,15672,15674,16000,16080,16992,17000,17070,17988,18080,18081,18082,18083,18084,18085,18086,18087,18088,18089,18090,18091,18092,18093,18094,18095,18096,18100,18200,18888,19000,19080,19150,19300,19315,19888,19999,20000,20010,20011,20012,20013,20014,20015,20016,20017,20018,20019,20020,20080,20443,20720,20880,22222,23333,23472,23791,24007,24430,25000,25025,25565,27000,27017,27018,27019,28017,28080,28443,30000,30001,30002,30003,30004,30005,30006,30007,30008,30009,30010,30821,31001,31337,32000,32400,32768,32769,33300,34205,35601,35729,36000,37777,38080,38888,40000,40001,40443,41080,43110,44300,44443,48080,49152,49153,50000,50001,50002,50003,50030,50060,50070,50075,50080,50090,50100,51000,51413,53413,55000,55555,60000,60001,60002,60003,60004,60005,60080,61613,61616,63330,64210,64738,65000,65432,65500,65535"

func ParseOptions() *Options {
	opt := &Options{}
	flagSet := goflags.NewFlagSet()

	flagSet.StringVar(&opt.Domain, "d", "", "domain to scan (e.g. example.com)")
	flagSet.StringVar(&opt.Workdir, "w", "./", "working directory for output files, defaults to current directory")
	flagSet.BoolVar(&opt.Server, "server", false, "start gowitness server at the end of scan to view screenshots")
	flagSet.StringVar(&opt.Profile, "profile", "balanced", "scan profile: balanced, stealth, aggressive, bugbounty")
	flagSet.BoolVar(&opt.SummaryJSON, "summary-json", true, "write scan summary as JSON file")
	flagSet.StringVar(&opt.SummaryFile, "summary-file", "scan_summary.json", "scan summary JSON output filename")
	flagSet.IntVar(&opt.SubfinderThreads, "subfinder-threads", 10, "number of subfinder threads")
	flagSet.IntVar(&opt.SubfinderTimeout, "subfinder-timeout", 30, "subfinder timeout in seconds")
	flagSet.IntVar(&opt.SubfinderMaxTime, "subfinder-max-time", 10, "subfinder max enumeration time in minutes")
	flagSet.IntVar(&opt.HttpxThreads, "httpx-threads", 50, "number of httpx threads")
	flagSet.IntVar(&opt.HttpxRateLimit, "httpx-rate-limit", 150, "httpx requests per second")
	flagSet.IntVar(&opt.HttpxTimeout, "httpx-timeout", 10, "httpx timeout in seconds")
	flagSet.StringVar(&opt.HttpxPorts, "httpx-ports", "", "custom httpx ports list (nmap-style, e.g. 80,443,8080 or http:80,https:8443)")
	flagSet.IntVar(&opt.KatanaDepth, "katana-depth", 3, "maximum katana crawl depth")
	flagSet.IntVar(&opt.KatanaTimeout, "katana-timeout", 10, "katana timeout in seconds")
	flagSet.IntVar(&opt.KatanaConcurrency, "katana-concurrency", 100, "number of katana concurrent crawling goroutines")
	flagSet.IntVar(&opt.KatanaParallelism, "katana-parallelism", 100, "number of katana URLs processing goroutines")
	flagSet.IntVar(&opt.KatanaRateLimit, "katana-rate-limit", 150, "katana requests per second")

	if err := flagSet.Parse(); err != nil {
		logging.LogError("Error parsing flags:", err)
		os.Exit(1)
	}

	setFlags := parseSetFlags(os.Args[1:])

	validateDomain(opt)
	applyProfile(opt, setFlags)
	validateTunables(opt)
	validateSummaryOptions(opt)

	return opt
}

func parseSetFlags(args []string) map[string]bool {
	setFlags := make(map[string]bool)

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			continue
		}

		clean := strings.TrimLeft(arg, "-")
		if clean == "" {
			continue
		}

		if idx := strings.Index(clean, "="); idx >= 0 {
			clean = clean[:idx]
		}

		setFlags[clean] = true
	}

	return setFlags
}

func applyProfile(opt *Options, setFlags map[string]bool) {
	preset, ok := getProfilePreset(strings.ToLower(strings.TrimSpace(opt.Profile)))
	if !ok {
		logging.LogError("Invalid profile. Use one of: balanced, stealth, aggressive, bugbounty", nil)
		os.Exit(1)
	}

	opt.Profile = strings.ToLower(strings.TrimSpace(opt.Profile))

	if !setFlags["subfinder-threads"] {
		opt.SubfinderThreads = preset.SubfinderThreads
	}
	if !setFlags["subfinder-timeout"] {
		opt.SubfinderTimeout = preset.SubfinderTimeout
	}
	if !setFlags["subfinder-max-time"] {
		opt.SubfinderMaxTime = preset.SubfinderMaxTime
	}
	if !setFlags["httpx-threads"] {
		opt.HttpxThreads = preset.HttpxThreads
	}
	if !setFlags["httpx-rate-limit"] {
		opt.HttpxRateLimit = preset.HttpxRateLimit
	}
	if !setFlags["httpx-timeout"] {
		opt.HttpxTimeout = preset.HttpxTimeout
	}
	if !setFlags["httpx-ports"] && strings.TrimSpace(preset.HttpxPorts) != "" {
		opt.HttpxPorts = preset.HttpxPorts
	}
	if !setFlags["katana-depth"] {
		opt.KatanaDepth = preset.KatanaDepth
	}
	if !setFlags["katana-timeout"] {
		opt.KatanaTimeout = preset.KatanaTimeout
	}
	if !setFlags["katana-concurrency"] {
		opt.KatanaConcurrency = preset.KatanaConcurrency
	}
	if !setFlags["katana-parallelism"] {
		opt.KatanaParallelism = preset.KatanaParallelism
	}
	if !setFlags["katana-rate-limit"] {
		opt.KatanaRateLimit = preset.KatanaRateLimit
	}
}

func getProfilePreset(profile string) (profilePreset, bool) {
	switch profile {
	case "balanced":
		return profilePreset{
			SubfinderThreads:  10,
			SubfinderTimeout:  30,
			SubfinderMaxTime:  10,
			HttpxThreads:      50,
			HttpxRateLimit:    150,
			HttpxTimeout:      10,
			HttpxPorts:        "",
			KatanaDepth:       3,
			KatanaTimeout:     10,
			KatanaConcurrency: 100,
			KatanaParallelism: 100,
			KatanaRateLimit:   150,
		}, true
	case "stealth":
		return profilePreset{
			SubfinderThreads:  5,
			SubfinderTimeout:  45,
			SubfinderMaxTime:  20,
			HttpxThreads:      20,
			HttpxRateLimit:    60,
			HttpxTimeout:      15,
			HttpxPorts:        "",
			KatanaDepth:       2,
			KatanaTimeout:     15,
			KatanaConcurrency: 40,
			KatanaParallelism: 40,
			KatanaRateLimit:   60,
		}, true
	case "aggressive":
		return profilePreset{
			SubfinderThreads:  30,
			SubfinderTimeout:  25,
			SubfinderMaxTime:  15,
			HttpxThreads:      500,
			HttpxRateLimit:    500,
			HttpxTimeout:      8,
			HttpxPorts:        "66,80,81,82,83,84,85,86,87,88,89,90,443,444,445,3000,3001,3002,8080,8081,8082,8443,8888,9000,9443,10000,20000,30000,50000,65535",
			KatanaDepth:       4,
			KatanaTimeout:     10,
			KatanaConcurrency: 150,
			KatanaParallelism: 150,
			KatanaRateLimit:   300,
		}, true
	case "bugbounty":
		return profilePreset{
			SubfinderThreads:  25,
			SubfinderTimeout:  30,
			SubfinderMaxTime:  20,
			HttpxThreads:      500,
			HttpxRateLimit:    500,
			HttpxTimeout:      10,
			HttpxPorts:        bugBountyDefaultPorts,
			KatanaDepth:       4,
			KatanaTimeout:     12,
			KatanaConcurrency: 150,
			KatanaParallelism: 150,
			KatanaRateLimit:   250,
		}, true
	default:
		return profilePreset{}, false
	}
}

func validateDomain(opt *Options) {
	rawDomain := strings.TrimSpace(opt.Domain)
	if rawDomain == "" {
		logging.LogError("Domain is required. Use -d flag to specify a domain (e.g., -d example.com)", nil)
		os.Exit(1)
	}

	cleanDomain := strings.TrimPrefix(rawDomain, "https://")
	cleanDomain = strings.TrimPrefix(cleanDomain, "http://")
	cleanDomain = strings.TrimSuffix(cleanDomain, "/")

	if host, port, err := net.SplitHostPort(cleanDomain); err == nil {
		if _, convErr := strconv.Atoi(port); convErr != nil {
			logging.LogError("Invalid domain. Port must be numeric", convErr)
			os.Exit(1)
		}
		cleanDomain = host
	}

	cleanDomain = strings.TrimSuffix(cleanDomain, ".")

	if cleanDomain == "" || strings.ContainsAny(cleanDomain, " /?#") {
		logging.LogError("Invalid domain format. Use a valid host like example.com", nil)
		os.Exit(1)
	}

	opt.Domain = strings.ToLower(cleanDomain)
}

func validateTunables(opt *Options) {
	values := []struct {
		name  string
		value int
	}{
		{name: "subfinder-threads", value: opt.SubfinderThreads},
		{name: "subfinder-timeout", value: opt.SubfinderTimeout},
		{name: "subfinder-max-time", value: opt.SubfinderMaxTime},
		{name: "httpx-threads", value: opt.HttpxThreads},
		{name: "httpx-rate-limit", value: opt.HttpxRateLimit},
		{name: "httpx-timeout", value: opt.HttpxTimeout},
		{name: "katana-depth", value: opt.KatanaDepth},
		{name: "katana-timeout", value: opt.KatanaTimeout},
		{name: "katana-concurrency", value: opt.KatanaConcurrency},
		{name: "katana-parallelism", value: opt.KatanaParallelism},
		{name: "katana-rate-limit", value: opt.KatanaRateLimit},
	}

	for _, item := range values {
		if item.value <= 0 {
			logging.LogError(fmt.Sprintf("Invalid value for --%s. It must be greater than 0", item.name), nil)
			os.Exit(1)
		}
	}
}

func validateSummaryOptions(opt *Options) {
	if strings.TrimSpace(opt.SummaryFile) == "" {
		logging.LogError("Invalid value for --summary-file. It cannot be empty", nil)
		os.Exit(1)
	}
}
