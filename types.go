package radiobrowser

// Station represents a radio station returned by Radio Browser API.
type Station struct {
	ChangeUUID  string `json:"changeuuid"`
	StationUUID string `json:"stationuuid"`
	ServerUUID  string `json:"serveruuid"`

	Name        string `json:"name"`
	URL         string `json:"url"`
	URLResolved string `json:"url_resolved"`
	Homepage    string `json:"homepage"`
	Favicon     string `json:"favicon"`

	Tags string `json:"tags"`
	// Deprecated: use CountryCode. The API generates Country from the code.
	Country       string `json:"country"`
	CountryCode   string `json:"countrycode"`
	ISO3166_2     string `json:"iso_3166_2"`
	State         string `json:"state"`
	Language      string `json:"language"`
	LanguageCodes string `json:"languagecodes"`

	Votes             int    `json:"votes"`
	LastChangeTime    string `json:"lastchangetime"`
	LastChangeTimeISO string `json:"lastchangetime_iso8601"`

	Codec   string `json:"codec"`
	Bitrate int    `json:"bitrate"`

	// hls is returned by the API as 0 or 1.
	HLS int `json:"hls"`

	// lastcheckok is returned by the API as 0 or 1.
	LastCheckOK int `json:"lastcheckok"`

	LastCheckTime    string `json:"lastchecktime"`
	LastCheckTimeISO string `json:"lastchecktime_iso8601"`

	LastCheckOKTime string `json:"lastcheckoktime"`
	LastCheckOKISO  string `json:"lastcheckoktime_iso8601"`

	LastLocalCheckTime    string `json:"lastlocalchecktime"`
	LastLocalCheckTimeISO string `json:"lastlocalchecktime_iso8601"`

	ClickTimestamp    string `json:"clicktimestamp"`
	ClickTimestampISO string `json:"clicktimestamp_iso8601"`
	ClickCount        int    `json:"clickcount"`
	ClickTrend        int    `json:"clicktrend"`

	SSLError int `json:"ssl_error"`

	// Coordinates and distance are nil when the API has no corresponding data.
	GeoLat      *float64 `json:"geo_lat"`
	GeoLong     *float64 `json:"geo_long"`
	GeoDistance *float64 `json:"geo_distance"`

	HasExtendedInfo bool `json:"has_extended_info"`
}

// StationCheck represents a single availability check result for a station.
type StationCheck struct {
	CheckUUID   string `json:"checkuuid"`
	StationUUID string `json:"stationuuid"`
	Source      string `json:"source"`
	Codec       string `json:"codec"`
	Bitrate     int    `json:"bitrate"`

	// hls is returned by the API as 0 or 1.
	HLS int `json:"hls"`

	// ok is returned by the API as 0 or 1.
	OK int `json:"ok"`

	TimestampISO8601 string `json:"timestamp_iso8601"`
	Timestamp        string `json:"timestamp"`
	URLCache         string `json:"urlcache"`

	// metainfo_overrides_database is returned by the API as 0 or 1.
	MetaInfoOverridesDatabase int `json:"metainfo_overrides_database"`

	// public is nullable and returned as 0 or 1 when present.
	Public *int `json:"public"`

	Name                   *string `json:"name"`
	Description            *string `json:"description"`
	Tags                   *string `json:"tags"`
	CountryCode            *string `json:"countrycode"`
	CountrySubdivisionCode *string `json:"countrysubdivisioncode"`
	Homepage               *string `json:"homepage"`
	Favicon                *string `json:"favicon"`
	LoadBalancer           *string `json:"loadbalancer"`
	ServerSoftware         *string `json:"server_software"`
	Sampling               *int    `json:"sampling"`
	TimingMS               int     `json:"timing_ms"`
	LanguageCodes          *string `json:"languagecodes"`

	// ssl_error is returned by the API as 0 or 1.
	SSLError int `json:"ssl_error"`

	GeoLat  *float64 `json:"geo_lat"`
	GeoLong *float64 `json:"geo_long"`
}

// StationCheckStep represents a single step performed during a station stream check.
type StationCheckStep struct {
	StepUUID        string  `json:"stepuuid"`
	ParentStepUUID  *string `json:"parent_stepuuid"`
	CheckUUID       string  `json:"checkuuid"`
	StationUUID     string  `json:"stationuuid"`
	URL             string  `json:"url"`
	URLType         *string `json:"urltype"`
	Error           *string `json:"error"`
	CreationISO8601 string  `json:"creation_iso8601"`
}

// StreamingServer represents an audio streaming server (for example Icecast),
// returned by /json/streamingservers. It is not a Radio Browser API mirror.
type StreamingServer struct {
	UUID      string  `json:"uuid"`
	URL       string  `json:"url"`
	StatusURL *string `json:"statusurl"`
	Error     *string `json:"error"`
	Admin     string  `json:"admin"`
	Location  string  `json:"location"`
	Software  string  `json:"software"`
}

// Country represents an item from /json/countries.
type Country struct {
	Name         string `json:"name"`
	ISO3166_1    string `json:"iso_3166_1"`
	StationCount int    `json:"stationcount"`
}

// Codec represents an item from /json/codecs.
type Codec struct {
	Name         string `json:"name"`
	StationCount int    `json:"stationcount"`
}

// State represents an item from /json/states.
type State struct {
	Name         string `json:"name"`
	Country      string `json:"country"`
	StationCount int    `json:"stationcount"`
}

// Language represents an item from /json/languages.
type Language struct {
	Name         string `json:"name"`
	ISO639       string `json:"iso_639,omitempty"`
	StationCount int    `json:"stationcount"`
}

// Tag represents an item from /json/tags.
type Tag struct {
	Name         string `json:"name"`
	StationCount int    `json:"stationcount"`
}

// StationClick represents an item from /json/clicks.
type StationClick struct {
	StationUUID           string `json:"stationuuid"`
	ClickUUID             string `json:"clickuuid"`
	ClickTimestampISO8601 string `json:"clicktimestamp_iso8601"`
	ClickTimestamp        string `json:"clicktimestamp"`
}

// StationClickCounterResponse represents response from /json/url/{stationuuid}.
type StationClickCounterResponse struct {
	OK          bool   `json:"ok"`
	Message     string `json:"message"`
	StationUUID string `json:"stationuuid"`
	Name        string `json:"name"`
	URL         string `json:"url"`
}

// VoteResponse represents response from /json/vote/{stationuuid}.
type VoteResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AddStationResponse represents response from /json/add.
type AddStationResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	UUID    string `json:"uuid"`
}
