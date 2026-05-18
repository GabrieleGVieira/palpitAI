package matchsync

import (
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/repositories"
)

const (
	defaultRequestTimeout = 10 * time.Second
	livePollInterval      = 30 * time.Second
	rateLimitGap          = 6 * time.Second
	todayPollInterval     = 5 * time.Minute
	upcomingPollInterval  = time.Hour
	upcomingWindow        = 30 * 24 * time.Hour
)

type datastore = repositories.Datastore

type syncKind string

const (
	syncLive     syncKind = "live"
	syncToday    syncKind = "today"
	syncUpcoming syncKind = "upcoming"
)
