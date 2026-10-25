package middleware

import (
	"net/url"
	"strings"
)

// OriginConfig provides the exact browser origins permitted by CORS.
type OriginConfig interface {
	GetAllowedOrigins() []string
}

func allowedOriginSet(origins []string) map[string]struct{} {
	allowed := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		if normalized := normalizeOrigin(origin); normalized != "" {
			allowed[normalized] = struct{}{}
		}
	}
	return allowed
}

func normalizeOrigin(value string) string {
	parsed, err := url.Parse(strings.TrimSpace(value))
	if err != nil {
		return ""
	}
	scheme := strings.ToLower(parsed.Scheme)
	if parsed.Host == "" || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" || (scheme != "http" && scheme != "https") {
		return ""
	}
	return scheme + "://" + strings.ToLower(parsed.Host)
}
