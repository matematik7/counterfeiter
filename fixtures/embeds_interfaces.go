package fixtures

import (
	"net/http"

	"github.com/matematik7/counterfeiter/fixtures/another_package"
)

type EmbedsInterfaces interface {
	http.Handler
	another_package.AnotherInterface
	InterfaceToEmbed

	DoThings()
}

type InterfaceToEmbed interface {
	EmbeddedMethod() string
}
