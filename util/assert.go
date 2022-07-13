package util

import (
	"log"
	"reflect"
)

func AssertEqual(a, b any) {
	if !reflect.DeepEqual(a, b) {
		log.Fatalf("ASSERTION: %+v != %+v", a, b)
	}
}
