package radiobrowser

// ServerStats represents /json/stats. Click counters cover the named time window.
type ServerStats struct {
	SupportedVersion int    `json:"supported_version"`
	SoftwareVersion  string `json:"software_version"`
	Status           string `json:"status"`
	Stations         int    `json:"stations"`
	StationsBroken   int    `json:"stations_broken"`
	Tags             int    `json:"tags"`
	ClicksLastHour   int    `json:"clicks_last_hour"`
	ClicksLastDay    int    `json:"clicks_last_day"`
	Languages        int    `json:"languages"`
	Countries        int    `json:"countries"`
}

// ServerMirror is an API mirror returned by /json/servers. A hostname can appear
// multiple times with different IPv4 or IPv6 addresses.
type ServerMirror struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
}

// ServerConfig represents /json/config. All timeout and interval fields ending
// in Seconds are measured in seconds; CacheTTL is also measured in seconds.
type ServerConfig struct {
	CheckEnabled                             bool     `json:"check_enabled"`
	PrometheusExporterEnabled                bool     `json:"prometheus_exporter_enabled"`
	PullServers                              []string `json:"pull_servers"`
	TCPTimeoutSeconds                        int      `json:"tcp_timeout_seconds"`
	BrokenStationsNeverWorkingTimeoutSeconds int      `json:"broken_stations_never_working_timeout_seconds"`
	BrokenStationsTimeoutSeconds             int      `json:"broken_stations_timeout_seconds"`
	ChecksTimeoutSeconds                     int      `json:"checks_timeout_seconds"`
	ClickValidTimeoutSeconds                 int      `json:"click_valid_timeout_seconds"`
	ClicksTimeoutSeconds                     int      `json:"clicks_timeout_seconds"`
	MirrorPullIntervalSeconds                int      `json:"mirror_pull_interval_seconds"`
	UpdateCachesIntervalSeconds              int      `json:"update_caches_interval_seconds"`
	ServerName                               string   `json:"server_name"`
	ServerLocation                           string   `json:"server_location"`
	ServerCountryCode                        string   `json:"server_country_code"`
	CheckRetries                             int      `json:"check_retries"`
	CheckBatchSize                           int      `json:"check_batchsize"`
	CheckPauseSeconds                        int      `json:"check_pause_seconds"`
	APIThreads                               int      `json:"api_threads"`
	CacheType                                string   `json:"cache_type"`
	CacheTTL                                 int      `json:"cache_ttl"`
	LanguageReplaceFilepath                  string   `json:"language_replace_filepath"`
	LanguageToCodeFilepath                   string   `json:"language_to_code_filepath"`
}
