package healthv1

import (
	"net/http"

	"github.com/nagorneva/nagorneva_reviews/internal/pkg/httpx"
)

// Get returns the process health status for load balancers and local checks.
func Get(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
