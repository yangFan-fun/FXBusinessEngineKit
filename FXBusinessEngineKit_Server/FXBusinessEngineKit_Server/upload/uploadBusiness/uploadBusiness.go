package uploadBusiness

import (
	"FXBusinessEngineKit_Server/configuration"
	"FXBusinessEngineKit_Server/log"
	"FXBusinessEngineKit_Server/upload/uploadModel"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

func UploadWithFile(file *multipart.FileHeader, context *gin.Context, isOutsea bool) (string, string, error) {
	// 存到指定位置
	currentDate := time.Now().Unix()
	fileStr := strconv.FormatInt(currentDate, 10)
	fileName := fileStr + file.Filename
	filePath := path.Join("./upload", fileName)
	se := context.SaveUploadedFile(file, filePath)
	if se != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[文件上传] 保存文件失败 %s", se))
		return "", "", se
	}

	// 上传到对象存储
	localU, err := connectCOS(filePath, isOutsea)
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[文件上传] 获取存储地址失败 %s", err))
		return "", filePath, err
	}

	model := uploadModel.UploadModel{
		Url: localU,
	}

	modelJson, err := json.Marshal(model)
	if err != nil {
		return "", "", err
	}

	return string(modelJson), filePath, nil
}

func connectCOS(file string, isOutsea bool) (string, error) {

	COS := configuration.COSPath

	u, _ := url.Parse(COS)
	b := &cos.BaseURL{
		BucketURL: u,
	}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  configuration.COSSecretId,
			SecretKey: configuration.COSSecretKey,
		},
	})

	key := file

	result, _, err := client.Object.Upload(context.Background(), key, file, nil)
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[COS] 图片上传失败 %s", err))
		return "", err
	}

	localUrl := result.Location

	return localUrl, nil

}
