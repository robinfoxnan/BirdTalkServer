package debugapi

import (
	"birdtalk/server/core"
	"github.com/gin-gonic/gin"
	"net/http"
)

// debug?cmd=listusers&type=mem
// debug?cmd=listusers&type=redis
// debug?cmd=listusers&type=db
// https://127.0.0.1:7817/debug?cmd=listusers&type=mem
func DebugHandler(c *gin.Context) {
	// 获取参数
	cmd := c.Query("cmd")
	typ := c.Query("type")

	switch cmd {
	case "listusers":
		handleListUsers(c, typ)
		return
	case "finduser":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown cmd"})
		return
	}

	return
}

func handleListUsers(c *gin.Context, typ string) {
	users, err := core.GetAllUsers(typ)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cmd":   "listusers",
		"type":  typ,
		"count": len(users),
		"users": users,
	})
}
