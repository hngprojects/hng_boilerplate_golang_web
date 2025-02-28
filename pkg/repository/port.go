package repository

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func ResolvePortParsing(port string, logger *utility.Logger) string {
	if _, err := strconv.Atoi(port); err != nil {
		u, err := url.Parse(port)
		if err != nil {
			utility.LogAndPrint(logger, fmt.Sprintf("parsing url %v to get port failed with: %v", port, err))
			panic(err)
		}

		detectedPort := u.Port()
		if detectedPort == "" {
			utility.LogAndPrint(logger, fmt.Sprintf("detecting port from url %v failed with: %v", port, err))
			panic(err)
		}
		port = detectedPort
	}
	return port
}
