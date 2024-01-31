package runtimecontext

import (
	"context"
	"syscall/js"
)

type (
	contextKeyTriggerObj struct{}
	contextKeyRuntimeObj struct{}
)

func New(ctx context.Context, triggerObj, runtimeObj js.Value) context.Context {
	ctx = context.WithValue(ctx, contextKeyTriggerObj{}, triggerObj)
	ctx = context.WithValue(ctx, contextKeyRuntimeObj{}, runtimeObj)
	return ctx
}

// MustExtractTriggerObj extracts trigger object from context.
// This function panics when trigger object was not found.
func MustExtractTriggerObj(ctx context.Context) js.Value {
	v, ok := ctx.Value(contextKeyTriggerObj{}).(js.Value)
	if !ok {
		panic("trigger object was not found")
	}
	return v
}

// MustExtractRuntimeObj extracts runtime object from context.
// This function panics when runtime object was not found.
func MustExtractRuntimeObj(ctx context.Context) js.Value {
	v, ok := ctx.Value(contextKeyRuntimeObj{}).(js.Value)
	if !ok {
		panic("runtime object was not found")
	}
	return v
}
