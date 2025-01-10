package cfg

import (
	"os"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

type (
	ChangeHandler func(fsnotify.Event)
	ChangeMatcher func(fsnotify.Event) bool
)

func Restart(matchers ...ChangeMatcher) ChangeHandler {
	return func(e fsnotify.Event) {
		for _, m := range matchers {
			if !m(e) {
				return
			}
		}
		if err := syscall.Exec(os.Args[0], os.Args, os.Environ()); err != nil {
			panic(err)
		}
	}
}

func Reload(target any, matchers ...ChangeMatcher) ChangeHandler {
	return func(e fsnotify.Event) {
		for _, m := range matchers {
			if !m(e) {
				return
			}
		}
		if err := Unmarshal(target); err != nil {
			panic(err)
		}
	}
}

func ReloadKey(key string, target any, matchers ...ChangeMatcher) ChangeHandler {
	return func(e fsnotify.Event) {
		for _, m := range matchers {
			if !m(e) {
				return
			}
		}
		if err := UnmarshalKey(key, target); err != nil {
			panic(err)
		}
	}
}
