package apiserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wenwenxiong/go-ipam-client/pkg/client/goipam"
)

type handler struct {
	client goipam.Client
}

type Subnet struct {
	Cidr string `json:"cidr"`
}
type IP struct {
	Cidr string `json:"cidr"`
	Ip   string `json:"ip"`
}

func newHandler(c goipam.Client) *handler {
	return &handler{client: c}
}

func (h handler) GetPrefix(c *gin.Context) {
	var subnet Subnet
	err := c.ShouldBindJSON(&subnet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	result := h.client.GetPrefix(subnet.Cidr)
	c.JSON(http.StatusOK, gin.H{
		"cidr": result.Prefix.Cidr,
	})
}

func (h handler) CreatePrefix(c *gin.Context) {
	var subnet Subnet
	err := c.ShouldBindJSON(&subnet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	result := h.client.CreatePrefix(subnet.Cidr)
	c.JSON(http.StatusOK, gin.H{
		"cidr": result.Prefix.Cidr,
	})
}

func (h handler) DeletePrefix(c *gin.Context) {
	var subnet Subnet
	err := c.ShouldBindJSON(&subnet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	result := h.client.DeletePrefix(subnet.Cidr)
	c.JSON(http.StatusOK, gin.H{
		"cidr":    result.Prefix.Cidr,
		"message": "delete cidr success",
	})

}

func (h handler) AcquireIP(c *gin.Context) {
	var ip IP
	err := c.ShouldBindJSON(&ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	result := h.client.AcquireIP(ip.Cidr, ip.Ip)
	c.JSON(http.StatusOK, gin.H{
		"cidr": ip.Cidr,
		"ip":   result.Ip.Ip,
	})
}

func (h handler) ReleaseIP(c *gin.Context) {
	var ip IP
	err := c.ShouldBindJSON(&ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	result := h.client.ReleaseIP(ip.Cidr, ip.Ip)
	c.JSON(http.StatusOK, gin.H{
		"cidr":    ip.Cidr,
		"ip":      result.Ip.Ip,
		"message": "release ip success ",
	})
}
