package proxy

import (
	"github.com/likexian/whois"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"guarde/pkg/parser"
	"guarde/pkg/utils"
	"strings"
)

func (config *Configuration) IsAllowed(ip string) bool {
	allow := true
	if cached, ok := WhoIsCache.Get(ip); ok {
		allow = cached.(bool)
	} else {
		log.Info().Str("ip", ip).Msg("Requesting WHOIS properties.")
		scan, err := whois.Whois(ip)
		if err == nil {
			properties := parser.WhoIs(scan)
			for _, ruleset := range config.Ruleset {
				for key, value := range ruleset {
					key := strings.ToLower(key)
					property, ok := properties[key]
					if !ok {
						log.Warn().Str("property", key).Str("ip", ip).Msg("Failed to get ruleset property.")
						continue
					}
					if utils.HasPrefixStr(value, "%") {
						if !strings.Contains(property, value[1:]) {
							allow = false
							break
						}
					}
					if property != value {
						allow = false
						break
					}
				}
			}
		} else {
			log.Warn().Str("ip", ip).Msg("Failed to get WHOIS properties.")
		}
		WhoIsCache.Set(ip, allow, cache.DefaultExpiration)
	}
	return allow
}
