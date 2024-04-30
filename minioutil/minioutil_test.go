package minioutil

import (
	"fmt"
	"testing"
)

func TestMinioClient(t *testing.T) {
	client, isConnect := NewMinioClient("127.0.0.1:9000", "admin", "admin123", false, "", "beijing")
	if !isConnect {
		return
	}
	r := gin.Default()
	r.POST("/minio_upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		path := client.UploadFile("mytest", file.Filename, file)
		fmt.Println(path)
	})
	//默认为监听8080端口
	err := r.Run(":8000")
	if err != nil {
		return
	}
}
