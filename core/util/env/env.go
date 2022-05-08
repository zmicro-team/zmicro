//go:generate stringer -type=Env -linecomment

package env

import "os"

const EnvZmicro = "ZMICRO_ENV"

type Env int8

const (
	None Env = iota
	Develop
	Testing
	Staging
	Product
)

var current = Develop

func Set(s string) {
	switch s {
	case Develop.String():
		current = Develop
	case Testing.String():
		current = Testing
	case Staging.String():
		current = Staging
	case Product.String():
		current = Product
	default:
		current = None
	}
}

func Get() Env {
	if current == None {
		env := os.Getenv(EnvZmicro)
		if env == "" {
			current = Develop
		} else {
			Set(env)
			if current == None {
				current = Develop
			}
		}
	}

	return current
}

func SetDevelop() {
	current = Develop
}

func SetTesting() {
	current = Testing
}

func SetStaging() {
	current = Staging
}

func SetProduct() {
	current = Product
}

func IsDevelop() bool {
	return Get() == Develop
}

func IsTesting() bool {
	return Get() == Testing
}

func IsStaging() bool {
	return Get() == Staging
}

func IsProduct() bool {
	return Get() == Product
}
