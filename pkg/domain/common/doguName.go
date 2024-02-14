package common

type QualifiedDoguName struct {
	namespace      DoguNamespace
	simpleDoguName SimpleDoguName
}

type DoguNamespace string
type SimpleDoguName string
