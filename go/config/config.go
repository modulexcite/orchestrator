/*
   Copyright 2014 Outbrain Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package config

import (
	"encoding/json"
	"os"
	"regexp"

	"gopkg.in/gcfg.v1"

	"github.com/outbrain/golib/log"
)

var (
	envVariableRegexp = regexp.MustCompile("[$][{](.*)[}]")
)

// Configuration makes for orchestrator configuration input, which can be provided by user via JSON formatted file.
// Some of the parameteres have reasonable default values, and some (like database credentials) are
// strictly expected from user.
type Configuration struct {
	Debug                                        bool // set debug mode (similar to --debug option)
	EnableSyslog                                 bool // Should logs be directed (in addition) to syslog daemon?
	ListenAddress                                string
	AgentsServerPort                             string // port orchestrator agents talk back to
	MySQLTopologyUser                            string
	MySQLTopologyPassword                        string // my.cnf style configuration file from where to pick credentials. Expecting `user`, `password` under `[client]` section
	MySQLTopologyCredentialsConfigFile           string
	MySQLTopologySSLPrivateKeyFile               string // Private key file used to authenticate with a Topology mysql instance with TLS
	MySQLTopologySSLCertFile                     string // Certificate PEM file used to authenticate with a Topology mysql instance with TLS
	MySQLTopologySSLCAFile                       string // Certificate Authority PEM file used to authenticate with a Topology mysql instance with TLS
	MySQLTopologySSLSkipVerify                   bool   // If true, do not strictly validate mutual TLS certs for Topology mysql instances
	MySQLTopologyUseMutualTLS                    bool   // Turn on TLS authentication with the Topology MySQL instances
	MySQLTopologyMaxPoolConnections              int    // Max concurrent connections on any topology instance
	DatabaselessMode__experimental               bool   // !!!EXPERIMENTAL!!! Orchestrator will execute without speaking to a backend database; super-standalone mode
	MySQLOrchestratorHost                        string
	MySQLOrchestratorPort                        uint
	MySQLOrchestratorDatabase                    string
	MySQLOrchestratorUser                        string
	MySQLOrchestratorPassword                    string
	MySQLOrchestratorCredentialsConfigFile       string   // my.cnf style configuration file from where to pick credentials. Expecting `user`, `password` under `[client]` section
	MySQLOrchestratorSSLPrivateKeyFile           string   // Private key file used to authenticate with the Orchestrator mysql instance with TLS
	MySQLOrchestratorSSLCertFile                 string   // Certificate PEM file used to authenticate with the Orchestrator mysql instance with TLS
	MySQLOrchestratorSSLCAFile                   string   // Certificate Authority PEM file used to authenticate with the Orchestrator mysql instance with TLS
	MySQLOrchestratorSSLSkipVerify               bool     // If true, do not strictly validate mutual TLS certs for the Orchestrator mysql instances
	MySQLOrchestratorUseMutualTLS                bool     // Turn on TLS authentication with the Orchestrator MySQL instance
	MySQLConnectTimeoutSeconds                   int      // Number of seconds before connection is aborted (driver-side)
	DefaultInstancePort                          int      // In case port was not specified on command line
	SlaveLagQuery                                string   // custom query to check on slave lg (e.g. heartbeat table)
	SlaveStartPostWaitMilliseconds               int      // Time to wait after START SLAVE before re-readong instance (give slave chance to connect to master)
	DiscoverByShowSlaveHosts                     bool     // Attempt SHOW SLAVE HOSTS before PROCESSLIST
	InstancePollSeconds                          uint     // Number of seconds between instance reads
	ReadLongRunningQueries                       bool     // Whether orchestrator should read and record current long running executing queries.
	BinlogFileHistoryDays                        int      // When > 0, amount of days for which orchestrator records per-instance binlog files & sizes
	UnseenInstanceForgetHours                    uint     // Number of hours after which an unseen instance is forgotten
	SnapshotTopologiesIntervalHours              uint     // Interval in hour between snapshot-topologies invocation. Default: 0 (disabled)
	InstanceBulkOperationsWaitTimeoutSeconds     uint     // Time to wait on a single instance when doing bulk (many instances) operation
	ActiveNodeExpireSeconds                      uint     // Maximum time to wait for active node to send keepalive before attempting to take over as active node.
	HostnameResolveMethod                        string   // Method by which to "normalize" hostname ("none"/"default"/"cname")
	MySQLHostnameResolveMethod                   string   // Method by which to "normalize" hostname via MySQL server. ("none"/"@@hostname"/"@@report_host"; default "@@hostname")
	SkipBinlogServerUnresolveCheck               bool     // Skip the double-check that an unresolved hostname resolves back to same hostname for binlog servers
	ExpiryHostnameResolvesMinutes                int      // Number of minutes after which to expire hostname-resolves
	RejectHostnameResolvePattern                 string   // Regexp pattern for resolved hostname that will not be accepted (not cached, not written to db). This is done to avoid storing wrong resolves due to network glitches.
	ReasonableReplicationLagSeconds              int      // Above this value is considered a problem
	ProblemIgnoreHostnameFilters                 []string // Will minimize problem visualization for hostnames matching given regexp filters
	VerifyReplicationFilters                     bool     // Include replication filters check before approving topology refactoring
	MaintenanceOwner                             string   // (Default) name of maintenance owner to use if none provided
	ReasonableMaintenanceReplicationLagSeconds   int      // Above this value move-up and move-below are blocked
	MaintenanceExpireMinutes                     uint     // Minutes after which a maintenance flag is considered stale and is cleared
	MaintenancePurgeDays                         uint     // Days after which maintenance entries are purged from the database
	CandidateInstanceExpireMinutes               uint     // Minutes after which a suggestion to use an instance as a candidate slave (to be preferably promoted on master failover) is expired.
	AuditLogFile                                 string   // Name of log file for audit operations. Disabled when empty.
	AuditToSyslog                                bool     // If true, audit messages are written to syslog
	AuditPageSize                                int
	AuditPurgeDays                               uint   // Days after which audit entries are purged from the database
	RemoveTextFromHostnameDisplay                string // Text to strip off the hostname on cluster/clusters pages
	ReadOnly                                     bool
	AuthenticationMethod                         string            // Type of autherntication to use, if any. "" for none, "basic" for BasicAuth, "multi" for advanced BasicAuth, "proxy" for forwarded credentials via reverse proxy, "token" for token based access
	HTTPAuthUser                                 string            // Username for HTTP Basic authentication (blank disables authentication)
	HTTPAuthPassword                             string            // Password for HTTP Basic authentication
	AuthUserHeader                               string            // HTTP header indicating auth user, when AuthenticationMethod is "proxy"
	PowerAuthUsers                               []string          // On AuthenticationMethod == "proxy", list of users that can make changes. All others are read-only.
	AccessTokenUseExpirySeconds                  uint              // Time by which an issued token must be used
	AccessTokenExpiryMinutes                     uint              // Time after which HTTP access token expires
	ClusterNameToAlias                           map[string]string // map between regex matching cluster name to a human friendly alias
	DetectClusterAliasQuery                      string            // Optional query (executed on topology instance) that returns the alias of a cluster. Query will only be executed on cluster master (though until the topology's master is resovled it may execute on other/all slaves). If provided, must return one row, one column
	DetectClusterDomainQuery                     string            // Optional query (executed on topology instance) that returns the VIP/CNAME/Alias/whatever domain name for the master of this cluster. Query will only be executed on cluster master (though until the topology's master is resovled it may execute on other/all slaves). If provided, must return one row, one column
	DetectInstanceAliasQuery                     string            // Optional query (executed on topology instance) that returns the alias of an instance. If provided, must return one row, one column
	DataCenterPattern                            string            // Regexp pattern with one group, extracting the datacenter name from the hostname
	PhysicalEnvironmentPattern                   string            // Regexp pattern with one group, extracting physical environment info from hostname (e.g. combination of datacenter & prod/dev env)
	DetectDataCenterQuery                        string            // Optional query (executed on topology instance) that returns the data center of an instance. If provided, must return one row, one column. Overrides DataCenterPattern and useful for installments where DC cannot be inferred by hostname
	DetectPhysicalEnvironmentQuery               string            // Optional query (executed on topology instance) that returns the physical environment of an instance. If provided, must return one row, one column. Overrides PhysicalEnvironmentPattern and useful for installments where env cannot be inferred by hostname
	DetectSemiSyncEnforcedQuery                  string            // Optional query (executed on topology instance) to determine whether semi-sync is fully enforced for master writes (async fallback is not allowed under any circumstance). If provided, must return one row, one column, value 0 or 1.
	SupportFuzzyPoolHostnames                    bool              // Should "submit-pool-instances" command be able to pass list of fuzzy instances (fuzzy means non-fqdn, but unique enough to recognize). Defaults 'true', implies more queries on backend db
	PromotionIgnoreHostnameFilters               []string          // Orchestrator will not promote slaves with hostname matching pattern (via -c recovery; for example, avoid promoting dev-dedicated machines)
	ServeAgentsHttp                              bool              // Spawn another HTTP interface dedicated for orcehstrator-agent
	AgentsUseSSL                                 bool              // When "true" orchestrator will listen on agents port with SSL as well as connect to agents via SSL
	AgentsUseMutualTLS                           bool              // When "true" Use mutual TLS for the server to agent communication
	AgentSSLSkipVerify                           bool              // When using SSL for the Agent, should we ignore SSL certification error
	AgentSSLPrivateKeyFile                       string            // Name of Agent SSL private key file, applies only when AgentsUseSSL = true
	AgentSSLCertFile                             string            // Name of Agent SSL certification file, applies only when AgentsUseSSL = true
	AgentSSLCAFile                               string            // Name of the Agent Certificate Authority file, applies only when AgentsUseSSL = true
	AgentSSLValidOUs                             []string          // Valid organizational units when using mutual TLS to communicate with the agents
	UseSSL                                       bool              // Use SSL on the server web port
	UseMutualTLS                                 bool              // When "true" Use mutual TLS for the server's web and API connections
	SSLSkipVerify                                bool              // When using SSL, should we ignore SSL certification error
	SSLPrivateKeyFile                            string            // Name of SSL private key file, applies only when UseSSL = true
	SSLCertFile                                  string            // Name of SSL certification file, applies only when UseSSL = true
	SSLCAFile                                    string            // Name of the Certificate Authority file, applies only when UseSSL = true
	SSLValidOUs                                  []string          // Valid organizational units when using mutual TLS
	StatusEndpoint                               string            // Override the status endpoint.  Defaults to '/api/status'
	StatusSimpleHealth                           bool              // If true, calling the status endpoint will use the simplified health check
	StatusOUVerify                               bool              // If true, try to verify OUs when Mutual TLS is on.  Defaults to false
	HttpTimeoutSeconds                           int               // Number of idle seconds before HTTP GET request times out (when accessing orchestrator-agent)
	AgentPollMinutes                             uint              // Minutes between agent polling
	AgentAutoDiscover                            bool              // If true, instances should automatically discover when an agent is submitted
	UnseenAgentForgetHours                       uint              // Number of hours after which an unseen agent is forgotten
	StaleSeedFailMinutes                         uint              // Number of minutes after which a stale (no progress) seed is considered failed.
	SeedAcceptableBytesDiff                      int64             // Difference in bytes between seed source & target data size that is still considered as successful copy
	PseudoGTIDPattern                            string            // Pattern to look for in binary logs that makes for a unique entry (pseudo GTID). When empty, Pseudo-GTID based refactoring is disabled.
	PseudoGTIDPatternIsFixedSubstring            bool              // If true, then PseudoGTIDPattern is not treated as regular expression but as fixed substring, and can boost search time
	PseudoGTIDMonotonicHint                      string            // subtring in Pseudo-GTID entry which indicates Pseudo-GTID entries are expected to be monotonically increasing
	DetectPseudoGTIDQuery                        string            // Optional query which is used to authoritatively decide whether pseudo gtid is enabled on instance
	PseudoGTIDCoordinatesHistoryHeuristicMinutes int               // Significantly reducing Pseudo-GTID lookup time, this indicates the most recent N minutes binlog position where search for Pseudo-GTID will heuristically begin (there is a fallback on fullscan if unsuccessful)
	BinlogEventsChunkSize                        int               // Chunk size (X) for SHOW BINLOG|RELAYLOG EVENTS LIMIT ?,X statements. Smaller means less locking and mroe work to be done
	BufferBinlogEvents                           bool              // Should we used buffered read on SHOW BINLOG|RELAYLOG EVENTS -- releases the database lock sooner (recommended)
	SkipBinlogEventsContaining                   []string          // When scanning/comparing binlogs for Pseudo-GTID, skip entries containing given texts. These are NOT regular expressions (would consume too much CPU while scanning binlogs), just substrings to find.
	ReduceReplicationAnalysisCount               bool              // When true, replication analysis will only report instances where possibility of handled problems is possible in the first place (e.g. will not report most leaf nodes, that are mostly uninteresting). When false, provides an entry for every known instance
	FailureDetectionPeriodBlockMinutes           int               // The time for which an instance's failure discovery is kept "active", so as to avoid concurrent "discoveries" of the instance's failure; this preceeds any recovery process, if any.
	RecoveryPollSeconds                          int               // Interval between checks for a recovery scenario and initiation of a recovery process
	RecoveryPeriodBlockMinutes                   int               // (supported for backwards compatibility but please use newer `RecoveryPeriodBlockSeconds` instead) The time for which an instance's recovery is kept "active", so as to avoid concurrent recoveries on smae instance as well as flapping
	RecoveryPeriodBlockSeconds                   int               // (overrides `RecoveryPeriodBlockMinutes`) The time for which an instance's recovery is kept "active", so as to avoid concurrent recoveries on smae instance as well as flapping
	RecoveryIgnoreHostnameFilters                []string          // Recovery analysis will completely ignore hosts matching given patterns
	RecoverMasterClusterFilters                  []string          // Only do master recovery on clusters matching these regexp patterns (of course the ".*" pattern matches everything)
	RecoverIntermediateMasterClusterFilters      []string          // Only do IM recovery on clusters matching these regexp patterns (of course the ".*" pattern matches everything)
	OnFailureDetectionProcesses                  []string          // Processes to execute when detecting a failover scenario (before making a decision whether to failover or not). May and should use some of these placeholders: {failureType}, {failureDescription}, {failedHost}, {failureCluster}, {failureClusterAlias}, {failureClusterDomain}, {failedPort}, {successorHost}, {successorPort}, {successorAlias}, {countSlaves}, {slaveHosts}, {isDowntimed}, {autoMasterRecovery}, {autoIntermediateMasterRecovery}
	PreFailoverProcesses                         []string          // Processes to execute before doing a failover (aborting operation should any once of them exits with non-zero code; order of execution undefined). May and should use some of these placeholders: {failureType}, {failureDescription}, {failedHost}, {failureCluster}, {failureClusterAlias}, {failureClusterDomain}, {failedPort}, {successorHost}, {successorPort}, {successorAlias}, {countSlaves}, {slaveHosts}, {isDowntimed}
	PostFailoverProcesses                        []string          // Processes to execute after doing a failover (order of execution undefined). May and should use some of these placeholders: {failureType}, {failureDescription}, {failedHost}, {failureCluster}, {failureClusterAlias}, {failureClusterDomain}, {failedPort}, {successorHost}, {successorPort}, {successorAlias}, {countSlaves}, {slaveHosts}, {isDowntimed}, {isSuccessful}, {lostSlaves}
	PostUnsuccessfulFailoverProcesses            []string          // Processes to execute after a not-completely-successful failover (order of execution undefined). May and should use some of these placeholders: {failureType}, {failureDescription}, {failedHost}, {failureCluster}, {failureClusterAlias}, {failureClusterDomain}, {failedPort}, {successorHost}, {successorPort}, {successorAlias}, {countSlaves}, {slaveHosts}, {isDowntimed}, {isSuccessful}, {lostSlaves}
	PostMasterFailoverProcesses                  []string          // Processes to execute after doing a master failover (order of execution undefined). Uses same placeholders as PostFailoverProcesses
	PostIntermediateMasterFailoverProcesses      []string          // Processes to execute after doing a master failover (order of execution undefined). Uses same placeholders as PostFailoverProcesses
	CoMasterRecoveryMustPromoteOtherCoMaster     bool              // When 'false', anything can get promoted (and candidates are prefered over others). When 'true', orchestrator will promote the other co-master or else fail
	DetachLostSlavesAfterMasterFailover          bool              // Should slaves that are not to be lost in master recovery (i.e. were more up-to-date than promoted slave) be forcibly detached
	ApplyMySQLPromotionAfterMasterFailover       bool              // Should orchestrator take upon itself to apply MySQL master promotion: set read_only=0, detach replication, etc.
	MasterFailoverLostInstancesDowntimeMinutes   uint              // Number of minutes to downtime any server that was lost after a master failover (including failed master & lost slaves). 0 to disable
	MasterFailoverDetachSlaveMasterHost          bool              // Should orchestrator issue a detach-slave-master-host on newly promoted master (this makes sure the new master will not attempt to replicate old master if that comes back to life). Defaults 'false'. Meaningless if ApplyMySQLPromotionAfterMasterFailover is 'true'.
	PostponeSlaveRecoveryOnLagMinutes            uint              // On crash recovery, slaves that are lagging more than given minutes are only resurrected late in the recovery process, after master/IM has been elected and processes executed. Value of 0 disables this feature
	OSCIgnoreHostnameFilters                     []string          // OSC slaves recommendation will ignore slave hostnames matching given patterns
	GraphiteAddr                                 string            // Optional; address of graphite port. If supplied, metrics will be written here
	GraphitePath                                 string            // Prefix for graphite path. May include {hostname} magic placeholder
	GraphiteConvertHostnameDotsToUnderscores     bool              // If true, then hostname's dots are converted to underscores before being used in graphite path
	GraphitePollSeconds                          int               // Graphite writes interval. 0 disables.
}

// ToJSONString will marshal this configuration as JSON
func (this *Configuration) ToJSONString() string {
	b, _ := json.Marshal(this)
	return string(b)
}

func (this *Configuration) GetDiscoveryPollSeconds() uint {
	// Turning `DiscoveryPollSeconds` into hard coded value. I see no reason anymore why this would be configurable.
	// After a couple years working with this I just set it to 1 whereever.
	// Not making this value configurable solves other problems like inter-dependency between this value
	// and others, like `InstancePollSeconds`, like cache timeouts etc.
	return 1
}

// Config is *the* configuration instance, used globally to get configuration data
var Config = newConfiguration()
var readFileNames []string

func newConfiguration() *Configuration {
	return &Configuration{
		Debug:                                        false,
		EnableSyslog:                                 false,
		ListenAddress:                                ":3000",
		AgentsServerPort:                             ":3001",
		StatusEndpoint:                               "/api/status",
		StatusSimpleHealth:                           true,
		StatusOUVerify:                               false,
		MySQLOrchestratorPort:                        3306,
		MySQLTopologyMaxPoolConnections:              3,
		MySQLTopologyUseMutualTLS:                    false,
		DatabaselessMode__experimental:               false,
		MySQLOrchestratorUseMutualTLS:                false,
		MySQLConnectTimeoutSeconds:                   2,
		DefaultInstancePort:                          3306,
		InstancePollSeconds:                          5,
		ReadLongRunningQueries:                       true,
		BinlogFileHistoryDays:                        0,
		UnseenInstanceForgetHours:                    240,
		SnapshotTopologiesIntervalHours:              0,
		SlaveStartPostWaitMilliseconds:               1000,
		DiscoverByShowSlaveHosts:                     false,
		InstanceBulkOperationsWaitTimeoutSeconds:     10,
		ActiveNodeExpireSeconds:                      5,
		HostnameResolveMethod:                        "default",
		MySQLHostnameResolveMethod:                   "@@hostname",
		SkipBinlogServerUnresolveCheck:               true,
		ExpiryHostnameResolvesMinutes:                60,
		RejectHostnameResolvePattern:                 "",
		ReasonableReplicationLagSeconds:              10,
		ProblemIgnoreHostnameFilters:                 []string{},
		VerifyReplicationFilters:                     false,
		MaintenanceOwner:                             "orchestrator",
		ReasonableMaintenanceReplicationLagSeconds:   20,
		MaintenanceExpireMinutes:                     10,
		MaintenancePurgeDays:                         365,
		CandidateInstanceExpireMinutes:               60,
		AuditLogFile:                                 "",
		AuditToSyslog:                                false,
		AuditPageSize:                                20,
		AuditPurgeDays:                               365,
		RemoveTextFromHostnameDisplay:                "",
		ReadOnly:                                     false,
		AuthenticationMethod:                         "basic",
		HTTPAuthUser:                                 "",
		HTTPAuthPassword:                             "",
		AuthUserHeader:                               "X-Forwarded-User",
		PowerAuthUsers:                               []string{"*"},
		AccessTokenUseExpirySeconds:                  60,
		AccessTokenExpiryMinutes:                     1440,
		ClusterNameToAlias:                           make(map[string]string),
		DetectClusterAliasQuery:                      "",
		DetectClusterDomainQuery:                     "",
		DetectInstanceAliasQuery:                     "",
		DataCenterPattern:                            "",
		PhysicalEnvironmentPattern:                   "",
		DetectDataCenterQuery:                        "",
		DetectPhysicalEnvironmentQuery:               "",
		DetectSemiSyncEnforcedQuery:                  "",
		SupportFuzzyPoolHostnames:                    true,
		PromotionIgnoreHostnameFilters:               []string{},
		ServeAgentsHttp:                              false,
		AgentsUseSSL:                                 false,
		AgentsUseMutualTLS:                           false,
		AgentSSLValidOUs:                             []string{},
		AgentSSLSkipVerify:                           false,
		AgentSSLPrivateKeyFile:                       "",
		AgentSSLCertFile:                             "",
		AgentSSLCAFile:                               "",
		UseSSL:                                       false,
		UseMutualTLS:                                 false,
		SSLValidOUs:                                  []string{},
		SSLSkipVerify:                                false,
		SSLPrivateKeyFile:                            "",
		SSLCertFile:                                  "",
		SSLCAFile:                                    "",
		HttpTimeoutSeconds:                           60,
		AgentPollMinutes:                             60,
		AgentAutoDiscover:                            false,
		UnseenAgentForgetHours:                       6,
		StaleSeedFailMinutes:                         60,
		SeedAcceptableBytesDiff:                      8192,
		PseudoGTIDPattern:                            "",
		PseudoGTIDPatternIsFixedSubstring:            false,
		PseudoGTIDMonotonicHint:                      "",
		DetectPseudoGTIDQuery:                        "",
		PseudoGTIDCoordinatesHistoryHeuristicMinutes: 2,
		BinlogEventsChunkSize:                        10000,
		BufferBinlogEvents:                           true,
		SkipBinlogEventsContaining:                   []string{},
		ReduceReplicationAnalysisCount:               true,
		FailureDetectionPeriodBlockMinutes:           60,
		RecoveryPollSeconds:                          10,
		RecoveryPeriodBlockMinutes:                   60,
		RecoveryPeriodBlockSeconds:                   3600,
		RecoveryIgnoreHostnameFilters:                []string{},
		RecoverMasterClusterFilters:                  []string{},
		RecoverIntermediateMasterClusterFilters:      []string{},
		OnFailureDetectionProcesses:                  []string{},
		PreFailoverProcesses:                         []string{},
		PostMasterFailoverProcesses:                  []string{},
		PostIntermediateMasterFailoverProcesses:      []string{},
		PostFailoverProcesses:                        []string{},
		PostUnsuccessfulFailoverProcesses:            []string{},
		CoMasterRecoveryMustPromoteOtherCoMaster:     true,
		DetachLostSlavesAfterMasterFailover:          true,
		ApplyMySQLPromotionAfterMasterFailover:       false,
		MasterFailoverLostInstancesDowntimeMinutes:   0,
		MasterFailoverDetachSlaveMasterHost:          false,
		PostponeSlaveRecoveryOnLagMinutes:            0,
		OSCIgnoreHostnameFilters:                     []string{},
		GraphiteAddr:                                 "",
		GraphitePath:                                 "",
		GraphiteConvertHostnameDotsToUnderscores:     true,
		GraphitePollSeconds:                          60,
	}
}

func postReadAdjustments() {
	if Config.MySQLOrchestratorCredentialsConfigFile != "" {
		mySQLConfig := struct {
			Client struct {
				User     string
				Password string
			}
		}{}
		err := gcfg.ReadFileInto(&mySQLConfig, Config.MySQLOrchestratorCredentialsConfigFile)
		if err != nil {
			log.Fatalf("Failed to parse gcfg data from file: %+v", err)
		} else {
			log.Debugf("Parsed orchestrator credentials from %s", Config.MySQLOrchestratorCredentialsConfigFile)
			Config.MySQLOrchestratorUser = mySQLConfig.Client.User
			Config.MySQLOrchestratorPassword = mySQLConfig.Client.Password
		}
	}
	{
		// We accept password in the form "${SOME_ENV_VARIABLE}" in which case we pull
		// the given variable from os env
		submatch := envVariableRegexp.FindStringSubmatch(Config.MySQLOrchestratorPassword)
		if len(submatch) > 1 {
			Config.MySQLOrchestratorPassword = os.Getenv(submatch[1])
		}
	}
	if Config.MySQLTopologyCredentialsConfigFile != "" {
		mySQLConfig := struct {
			Client struct {
				User     string
				Password string
			}
		}{}
		err := gcfg.ReadFileInto(&mySQLConfig, Config.MySQLTopologyCredentialsConfigFile)
		if err != nil {
			log.Fatalf("Failed to parse gcfg data from file: %+v", err)
		} else {
			log.Debugf("Parsed topology credentials from %s", Config.MySQLTopologyCredentialsConfigFile)
			Config.MySQLTopologyUser = mySQLConfig.Client.User
			Config.MySQLTopologyPassword = mySQLConfig.Client.Password
		}
	}
	{
		// We accept password in the form "${SOME_ENV_VARIABLE}" in which case we pull
		// the given variable from os env
		submatch := envVariableRegexp.FindStringSubmatch(Config.MySQLTopologyPassword)
		if len(submatch) > 1 {
			Config.MySQLTopologyPassword = os.Getenv(submatch[1])
		}
	}

	if Config.RecoveryPeriodBlockSeconds == 0 && Config.RecoveryPeriodBlockMinutes > 0 {
		// RecoveryPeriodBlockSeconds is a newer addition that overrides RecoveryPeriodBlockMinutes
		// The code does not consider RecoveryPeriodBlockMinutes anymore, but RecoveryPeriodBlockMinutes
		// still supported in config file for backwards compatibility
		Config.RecoveryPeriodBlockSeconds = Config.RecoveryPeriodBlockMinutes * 60
	}
}

// read reads configuration from given file, or silently skips if the file does not exist.
// If the file does exist, then it is expected to be in valid JSON format or the function bails out.
func read(fileName string) (*Configuration, error) {
	file, err := os.Open(fileName)
	if err == nil {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(Config)
		if err == nil {
			log.Infof("Read config: %s", fileName)
		} else {
			log.Fatal("Cannot read config file:", fileName, err)
		}
		postReadAdjustments()
	}
	return Config, err
}

// Read reads configuration from zero, either, some or all given files, in order of input.
// A file can override configuration provided in previous file.
func Read(fileNames ...string) *Configuration {
	for _, fileName := range fileNames {
		read(fileName)
	}
	readFileNames = fileNames
	return Config
}

// ForceRead reads configuration from given file name or bails out if it fails
func ForceRead(fileName string) *Configuration {
	_, err := read(fileName)
	if err != nil {
		log.Fatal("Cannot read config file:", fileName, err)
	}
	readFileNames = []string{fileName}
	return Config
}

// Reload re-reads configuration from last used files
func Reload() *Configuration {
	for _, fileName := range readFileNames {
		read(fileName)
	}
	return Config
}
