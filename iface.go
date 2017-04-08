package pbars

type ProgressReceiver interface {
	Update(progress, length int64)
}
