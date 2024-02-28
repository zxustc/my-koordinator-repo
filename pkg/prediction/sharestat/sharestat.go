package sharestate

import (
	"k8s.io/apimachinery/pkg/types"
)

type shareState struct {
}

type ShareState interface {
	UpdateShareState()
	DeleteShareState()
	GetShareState()
}

func UpdateShareState(types.NamespacedName) {

}

func DeleteShareState(types.NamespacedName) {

}

func GetShareState() {

}
