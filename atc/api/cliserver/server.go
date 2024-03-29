package cliserver

import "code.cloudfoundry.org/lager/v3"

type Server struct {
	logger          lager.Logger
	cliDownloadsDir string
}

func NewServer(logger lager.Logger, cliDownloadsDir string) *Server {
	return &Server{
		logger:          logger,
		cliDownloadsDir: cliDownloadsDir,
	}
}
