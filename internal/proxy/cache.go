package proxy

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var WhoIsCache = cache.New(30*24*time.Hour, 10*time.Minute)
