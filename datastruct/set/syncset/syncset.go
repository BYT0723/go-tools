package syncset

import (
	"sync"
)

type SyncSet[T comparable] sync.Map
