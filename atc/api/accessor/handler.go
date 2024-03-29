package accessor

import (
	"context"
	"net/http"

	"code.cloudfoundry.org/lager/v3"
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/auditor"
)

//counterfeiter:generate net/http.Handler

const accessorContextKey atc.ContextKey = "accessor"

//counterfeiter:generate . AccessFactory
type AccessFactory interface {
	Create(req *http.Request, role string) (Access, error)
}

func NewHandler(
	logger lager.Logger,
	action string,
	handler http.Handler,
	accessFactory AccessFactory,
	auditor auditor.Auditor,
	customRoles map[string]string,
) http.Handler {
	return &accessorHandler{
		logger:        logger,
		handler:       handler,
		accessFactory: accessFactory,
		action:        action,
		auditor:       auditor,
		customRoles:   customRoles,
	}
}

type accessorHandler struct {
	logger        lager.Logger
	action        string
	handler       http.Handler
	accessFactory AccessFactory
	auditor       auditor.Auditor
	customRoles   map[string]string
}

func (h *accessorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requiredRole := h.customRoles[h.action]
	if requiredRole == "" {
		requiredRole = DefaultRoles[h.action]
	}

	acc, err := h.accessFactory.Create(r, requiredRole)
	if err != nil {
		h.logger.Error("failed-to-construct-accessor", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	claims := acc.Claims()

	ctx := context.WithValue(r.Context(), accessorContextKey, acc)

	h.auditor.Audit(h.action, claims.UserName, r)
	h.handler.ServeHTTP(w, r.WithContext(ctx))
}

func GetAccessor(r *http.Request) Access {
	accessor := r.Context().Value(accessorContextKey)
	if accessor != nil {
		return accessor.(Access)
	}

	return &access{}
}
