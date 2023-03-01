package server

import (
	"net/http"

	"github.com/shivendra195/supplyChainManagement/utils"
	"github.com/sirupsen/logrus"
)

type healthResponse struct {
	Available bool `json:"up"`
}

// HealthCheck  godoc
// @Summary 	Health check endpoint
// @Tags 		health
// @Accept 		json
// @Produce 	json
// @Success     200 {object} healthResponse
// @Router 		/health [get]
func (srv *Server) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	logrus.WithField("test-key", "testing").WithField("test-key-2", "testing-2").Info("testing health route")
	utils.EncodeJSON200Body(w, healthResponse{
		Available: true})
}
