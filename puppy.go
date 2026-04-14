package puppy

import (
	"github.com/vaibhavrajchauhan-crx/dog"
)

func Bark() string {
	return "Woff!!"
}

func Barks() string {
	return "Woff!! Woff!! Woff!!"
}

func DocName(name string) string {
	return dog.Dog(name)
}
