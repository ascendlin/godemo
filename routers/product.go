package routers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

//Product结构体
type Product struct {
	Username    string    `json:"username" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Price       int       `json:"price" binding:"gte=0"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type productHandler struct {
	//读写锁
	sync.RWMutex
	//定义一个map,key为string,value为Product对象
	products map[string]Product
}

func newProductHandler() *productHandler {
	//以map方式内存缓存
	return &productHandler{
		products: make(map[string]Product),
	}
}

//创建订单
func (u *productHandler) Create(c *gin.Context) {
	u.Lock()
	defer u.Unlock()
	// 1. 声明接收的变量
	var product Product
	//将request中的body数据，自动按照json格式解析到结构体Product中
	if err := c.ShouldBindJSON(&product); err != nil {
		//gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 2. 参数校验，判断名称是否存在
	if _, ok := u.products[product.Name]; ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %s already exist", product.Name)})
		return
	}
	product.CreatedAt = time.Now()
	// 3. 逻辑处理
	u.products[product.Name] = product
	log.Printf("Register product %s success", product.Name)
	// 4. 返回结果
	c.JSON(http.StatusOK, product)
}

//get查询订单
func (u *productHandler) Get(c *gin.Context) {
	u.Lock()
	defer u.Unlock()
	//判断名称是否存在，存在则返回product，c.Param为获取url参数，如：ip:port/products/iphone13,则获取的name为iphone13
	product, ok := u.products[c.Param("name")]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("can not found product %s", c.Param("name"))})
		return
	}
	c.JSON(http.StatusOK, product)
}

func InitProductRouter(e *gin.Engine) {
	productHandler := newProductHandler()
	// 路由分组
	v1 := e.Group("/v1")
	{
		productv1 := v1.Group("/products")
		{
			// 路由匹配，ip:port/products,为创建订单
			productv1.POST("", productHandler.Create)
			// ip:port/products/订单name,为get查询订单
			productv1.GET(":name", productHandler.Get)
		}
	}
}
