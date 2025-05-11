package main

import (
	"github.com/trysourcetool/onprem-portal/internal/config"
	"github.com/trysourcetool/onprem-portal/internal/logger"
)

func init() {
	config.Init()
	logger.Init()
}

func main() {}
