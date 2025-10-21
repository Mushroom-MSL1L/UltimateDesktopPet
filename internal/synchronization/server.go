package synchronization

import (
	"github.com/gin-gonic/gin"
)

func Server(hostaddress, port string) {
	r := gin.Default()

	SwagSetUp(r, hostaddress, port)
	Routing(r)

	r.Run(hostaddress + ":" + port)
}
