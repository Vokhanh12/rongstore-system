package errors

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// ---- YAML schema ----

type catalogYAML struct {
	Version  int          `yaml:"version"`
	Service  string       `yaml:"service"`
	Errors   []entryYAML  `yaml:"errors"`
	Defaults defaultsYAML `yaml:"defaults"`
}

type entryYAML struct {
	Code       string `yaml:"code"`
	Key        string `yaml:"key"`
	HTTPStatus int    `yaml:"http_status"`
	GRPCCode   string `yaml:"grpc_code"` // optional (not used in this loader)
	Message    string `yaml:"message"`
	Severity   string `yaml:"severity"`  // optional
	Retryable  *bool  `yaml:"retryable"` // optional
}

type defaultsYAML struct {
	UnknownDomainKey struct {
		Code       string `yaml:"code"`
		HTTPStatus int    `yaml:"http_status"`
		Message    string `yaml:"message"`
	} `yaml:"unknown_domain_key"`
	InternalFallback struct {
		Code       string `yaml:"code"`
		HTTPStatus int    `yaml:"http_status"`
		Message    string `yaml:"message"`
	} `yaml:"internal_fallback"`
}

// ---- In-memory state ----

type mapping struct {
	Code    string
	Status  int
	Message string
}

var (
	mu             sync.RWMutex
	loadedByDomain = map[string]mapping{} // domain key -> mapping
	defaults       = struct {
		Unknown  mapping
		Internal mapping
	}{
		Unknown:  mapping{Code: "AUTH-VAL-999", Status: 400, Message: "Unknown domain error"},
		Internal: mapping{Code: "CORE-INF-000", Status: 500, Message: "Internal server error"},
	}
	loaded bool
)

// InitFromYAML loads the catalog once at startup.
func InitFromYAML(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read errors.yaml: %w", err)
	}
	var cat catalogYAML
	if err := yaml.Unmarshal(data, &cat); err != nil {
		return fmt.Errorf("parse errors.yaml: %w", err)
	}

	tmp := make(map[string]mapping, len(cat.Errors))
	for _, e := range cat.Errors {
		k := strings.TrimSpace(e.Key)
		if k == "" || e.Code == "" || e.HTTPStatus == 0 {
			continue // skip invalid rows
		}
		tmp[k] = mapping{
			Code:    e.Code,
			Status:  e.HTTPStatus,
			Message: e.Message,
		}
	}

	// set defaults if present
	if cat.Defaults.UnknownDomainKey.Code != "" {
		defaults.Unknown = mapping{
			Code:    cat.Defaults.UnknownDomainKey.Code,
			Status:  cat.Defaults.UnknownDomainKey.HTTPStatus,
			Message: cat.Defaults.UnknownDomainKey.Message,
		}
	}
	if cat.Defaults.InternalFallback.Code != "" {
		defaults.Internal = mapping{
			Code:    cat.Defaults.InternalFallback.Code,
			Status:  cat.Defaults.InternalFallback.HTTPStatus,
			Message: cat.Defaults.InternalFallback.Message,
		}
	}

	mu.Lock()
	loadedByDomain = tmp
	loaded = true
	mu.Unlock()
	return nil
}

// lookupDomain returns mapping for a domain key from YAML (if loaded) or from static map as fallback.
func lookupDomain(domainKey string) (mapping, bool) {
	mu.RLock()
	defer mu.RUnlock()
	if m, ok := loadedByDomain[domainKey]; ok {
		return m, true
	}
	// fallback to static map (mapping.go) if present
	if m, ok := domainToTransport[domainKey]; ok {
		return mapping{Code: m.Code, Status: m.Status, Message: m.Message}, true
	}
	return mapping{}, false
}

func defaultUnknown() mapping {
	mu.RLock()
	defer mu.RUnlock()
	return defaults.Unknown
}
func defaultInternal() mapping {
	mu.RLock()
	defer mu.RUnlock()
	return defaults.Internal
}

func IsCatalogLoaded() bool {
	mu.RLock()
	defer mu.RUnlock()
	return loaded
}
