package transport

import (
	"context"
	"runtime/debug"
	"strings"
)

const libraryModulePath = "github.com/randomtoy/radiobrowser-go"

var defaultUserAgent = func() string {
	info, _ := debug.ReadBuildInfo()
	return userAgentFromBuildInfo(info)
}()

func DefaultUserAgent() string {
	return defaultUserAgent
}

func userAgentFromBuildInfo(info *debug.BuildInfo) string {
	const prefix = "radiobrowser-go/"
	if info == nil {
		return prefix + "dev"
	}

	var module *debug.Module
	if info.Main.Path == libraryModulePath {
		module = &info.Main
	} else {
		for _, dep := range info.Deps {
			if dep.Path == libraryModulePath {
				module = dep
				break
			}
		}
	}
	if module == nil {
		return prefix + "dev"
	}
	if module.Replace != nil {
		module = module.Replace
	}
	if module.Version == "" || module.Version == "(devel)" {
		return prefix + "dev"
	}
	return prefix + strings.TrimPrefix(module.Version, "v")
}

type userAgentKey struct{}

func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, userAgentKey{}, strings.TrimSpace(userAgent))
}

func userAgent(ctx context.Context) string {
	if value, ok := ctx.Value(userAgentKey{}).(string); ok && value != "" {
		return value
	}
	return defaultUserAgent
}
