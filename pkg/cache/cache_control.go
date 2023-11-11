package cache

import (
	"context"
	"strings"

	"github.com/driif/echo-go-starter/pkg/logs"
)

// CacheControlDirective is a cache control directive
type CacheControlDirective uint8

const (
	CacheControlDirectiveNoCache CacheControlDirective = 1 << iota
	CacheControlDirectiveNoStore
)

// HasDirective returns true if the directive is set
func (d CacheControlDirective) HasDirective(dir CacheControlDirective) bool { return d&dir != 0 }

// AddDirective adds the directive
func (d *CacheControlDirective) AddDirective(dir CacheControlDirective) { *d |= dir }

// ClearDirective clears the directive
func (d *CacheControlDirective) ClearDirective(dir CacheControlDirective) { *d &= ^dir }

// ToggleDirective toggles the directive
func (d *CacheControlDirective) ToggleDirective(dir CacheControlDirective) { *d ^= dir }

// String returns the string representation of the cache control directive
func (d CacheControlDirective) String() string {
	res := make([]string, 0)

	if d.HasDirective(CacheControlDirectiveNoCache) {
		res = append(res, "no-cache")
	}
	if d.HasDirective(CacheControlDirectiveNoStore) {
		res = append(res, "no-store")
	}

	return strings.Join(res, "|")
}

// ParseCacheControlDirective parses a cache control directive
func ParseCacheControlDirective(d string) CacheControlDirective {
	parts := strings.Split(d, "=")
	switch strings.ToLower(parts[0]) {
	case "no-cache":
		return CacheControlDirectiveNoCache
	case "no-store":
		return CacheControlDirectiveNoStore
	default:
		return 0
	}
}

// ParseCacheControlHeader parses a cache control header
func ParseCacheControlHeader(val string) CacheControlDirective {
	res := CacheControlDirective(0)

	directives := strings.Split(val, ",")
	for _, dir := range directives {
		res = res | ParseCacheControlDirective(dir)
	}

	return res
}

// CacheControlDirectiveFromContext returns the cache control directive from the context
func CacheControlDirectiveFromContext(ctx context.Context) CacheControlDirective {
	d := ctx.Value(logs.CTXKeyCacheControl)
	if d == nil {
		return CacheControlDirective(0)
	}

	directive, ok := d.(CacheControlDirective)
	if !ok {
		return CacheControlDirective(0)
	}

	return directive
}
