package synchronization

import (
	"github.com/gin-gonic/gin"
)

func Routing(r *gin.Engine) {
	r.GET("/sync", ping)
}
