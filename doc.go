// Package radiobrowser provides a zero-dependency client for the Radio Browser
// JSON API. Functions accept a context for cancellation and a base URL selecting
// an API server. Use WithUserAgent to identify the calling application.
//
// Station discovery, metadata, checks, history and server information are read
// operations. AddStation, VoteForStation and StationClickCounter modify upstream
// data. The latter registers playback and returns the station's stream URL.
//
// Geographic fields are pointers so that missing data remains distinct from zero.
// Paginated options preserve server defaults unless Limit is positive or LimitSet
// is true. Metadata counts accept both numeric and legacy quoted JSON values.
package radiobrowser
