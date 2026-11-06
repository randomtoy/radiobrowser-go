package transport

import (
	"runtime/debug"
	"testing"
)

func TestUserAgentFromBuildInfo(t *testing.T) {
	for _, tc := range []struct {
		name string
		info *debug.BuildInfo
		want string
	}{
		{name: "no build info", want: "radiobrowser-go/dev"},
		{
			name: "unrelated modules",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/app", Version: "v9.0.0"},
				Deps: []*debug.Module{{Path: "example.com/other", Version: "v8.0.0"}},
			},
			want: "radiobrowser-go/dev",
		},
		{
			name: "released dependency",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/app", Version: "v9.0.0"},
				Deps: []*debug.Module{
					{Path: "example.com/other", Version: "v8.0.0"},
					{Path: libraryModulePath, Version: "v0.3.1"},
				},
			},
			want: "radiobrowser-go/0.3.1",
		},
		{
			name: "pseudo-version",
			info: &debug.BuildInfo{Deps: []*debug.Module{{
				Path: libraryModulePath, Version: "v0.3.2-0.20261106090000-0123456789ab",
			}}},
			want: "radiobrowser-go/0.3.2-0.20261106090000-0123456789ab",
		},
		{
			name: "versioned main module",
			info: &debug.BuildInfo{Main: debug.Module{Path: libraryModulePath, Version: "v0.4.0-rc.1"}},
			want: "radiobrowser-go/0.4.0-rc.1",
		},
		{
			name: "development main module",
			info: &debug.BuildInfo{Main: debug.Module{Path: libraryModulePath, Version: "(devel)"}},
			want: "radiobrowser-go/dev",
		},
		{
			name: "unversioned module",
			info: &debug.BuildInfo{Deps: []*debug.Module{{Path: libraryModulePath}}},
			want: "radiobrowser-go/dev",
		},
		{
			name: "local replacement",
			info: &debug.BuildInfo{Deps: []*debug.Module{{
				Path: libraryModulePath, Version: "v0.3.1",
				Replace: &debug.Module{Path: "../radiobrowser-go"},
			}}},
			want: "radiobrowser-go/dev",
		},
		{
			name: "versioned replacement",
			info: &debug.BuildInfo{Deps: []*debug.Module{{
				Path: libraryModulePath, Version: "v0.3.1",
				Replace: &debug.Module{Path: "example.com/fork", Version: "v0.5.2"},
			}}},
			want: "radiobrowser-go/0.5.2",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := userAgentFromBuildInfo(tc.info); got != tc.want {
				t.Errorf("User-Agent = %q, want %q", got, tc.want)
			}
		})
	}
}
