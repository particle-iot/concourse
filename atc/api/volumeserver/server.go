package volumeserver

import (
	"code.cloudfoundry.org/lager/v3"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/gc"
)

type Server struct {
	logger     lager.Logger
	repository db.VolumeRepository
	destroyer  gc.Destroyer
}

func NewServer(
	logger lager.Logger,
	volumeRepository db.VolumeRepository,
	destroyer gc.Destroyer,
) *Server {
	return &Server{
		logger:     logger,
		repository: volumeRepository,
		destroyer:  destroyer,
	}
}
