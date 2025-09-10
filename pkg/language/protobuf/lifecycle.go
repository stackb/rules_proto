package protobuf

import (
	"context"
)

// Before implements part of the language.LifecycleManager interface.
func (pl *protobufLang) Before(context.Context) {
}

// DoneGeneratingRules implements part of the language.LifecycleManager interface.
func (pl *protobufLang) DoneGeneratingRules() {
}

// AfterResolvingDeps implements part of the language.LifecycleManager interface.
func (pl *protobufLang) AfterResolvingDeps(context.Context) {
}
