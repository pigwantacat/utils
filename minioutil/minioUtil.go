package minioutil

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"mime/multipart"
	"strings"
	"time"
)

// MinioClient minio客户端结构体
type MinioClient struct {
	// minio客户端
	minioClient *minio.Client
	// 存储桶前缀
	urlPrefix string
	// 存储桶区域
	region string
	// 存储桶集合
	bucketsMap map[string]struct{}
}

// NewMinioClient 创建minio客户端
// @param endpoint: 存储桶地址
// @param accessKeyID: 访问密钥
// @param secretAccessKey: 访问密钥
// @param useSSL: 是否使用SSL
// @param urlPrefix: 存储桶前缀
// @return *MinioClient, bool minio客户端,连接是否成功
func NewMinioClient(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool, urlPrefix string, region string) (client *MinioClient, isConnect bool) {
	// 创建minio客户端
	c, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	// 测试连接
	_, err = c.ListBuckets(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	return &MinioClient{minioClient: c, urlPrefix: urlPrefix, region: region, bucketsMap: make(map[string]struct{})}, true
}

// checkBucket 检查存储桶是否存在，如果不存在则创建
// @param bucketName: 存储桶名称
// @return bool 存储桶是否存在
func (client *MinioClient) checkBucket(bucketName string) bool {
	// 检查存储桶是否存在
	exists, err := client.minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		fmt.Println("Error checking bucket:", err.Error())
		return false
	}
	// 如果存储桶存在
	if exists {
		return true
	}
	// 创建存储桶
	return client.createBucket(bucketName)
}

// createBucket 创建存储桶
// @param bucketName: 存储桶名称
// @return bool 是否创建成功
func (client *MinioClient) createBucket(bucketName string) bool {
	// 创建存储桶
	if err := client.minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: client.region}); err != nil {
		fmt.Println("Error creating bucket:", err.Error())
		return false
	}
	return true
}

// buildFileName 构建文件名
// @param objectName: 存储桶名称
// @return string 生成后的文件名
func (client *MinioClient) buildFileName(objectName string) string {
	index := strings.LastIndex(strings.TrimSpace(objectName), ".")
	now := time.Now()
	return fmt.Sprintf("%s/%s_%d%s", now.Format("2006/01/02"), objectName[:index], now.UnixNano(), objectName[index:])
}

// UploadFile 将给定的文件上传到指定的存储桶和对象名，并返回生成的URL。
// @param bucketName: 存储桶名称
// @param objectName: 对象名称
// @param fileHeader: 文件头信息
// @return string 返回上传后的URL
func (client *MinioClient) UploadFile(bucketName string, objectName string, fileHeader *multipart.FileHeader) string {
	// 打开文件
	file, err := fileHeader.Open()
	defer file.Close()
	// 检查存储桶是否存在
	if _, hasBucket := client.bucketsMap[bucketName]; !hasBucket {
		// 检查minio中是否存在该存储桶,如果不存在则创建
		if checkResult := client.checkBucket(bucketName); !checkResult {
			return ""
		}
		// 将存储桶添加到集合中
		client.bucketsMap[bucketName] = struct{}{}
	}
	// 保存的文件名
	fileName := client.buildFileName(objectName)
	// 将文件上传到指定的存储桶和对象名
	_, err = client.minioClient.PutObject(context.Background(), bucketName, fileName, file, fileHeader.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	// 返回生成的URL
	return client.urlPrefix + "/" + bucketName + "/" + fileName
}
