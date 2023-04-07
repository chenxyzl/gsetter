package examplepb

type DirtyCallback func(field int32)

type DirtyTracker interface {
	IsDirty(field int32) bool
	SetDirty(field int32)
	ResetDirty(field int32)
	ResetAllDirty()
	SetOnDirtyCallback(callback DirtyCallback)
}
