package domain

type DeepCopier[T any] interface {
	DeepCopy() T
}
