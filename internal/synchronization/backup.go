package synchronization

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//	@Summary		Ping
//	@Description	Pong
//	@Tags			temp
//	@Produce		json
//	@Success		200	{object}	string	"pong"
//	@Failure		500	{object}	string	"internal server error"
//	@Router			/sync [get]
func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "pong")
}
