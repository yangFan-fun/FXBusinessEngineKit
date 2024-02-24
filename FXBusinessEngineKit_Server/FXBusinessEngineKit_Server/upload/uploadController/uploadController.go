package uploadController

import (
	"FXBusinessEngineKit_Server/log"
	"FXBusinessEngineKit_Server/tool"
	"FXBusinessEngineKit_Server/upload/uploadBusiness"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func UploadFace(context *gin.Context) {

	file, err := context.FormFile("file")
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[文件上传] 文件上传获取失败 %s", err))
		return
	}

	localU, path, err := uploadBusiness.UploadWithFile(file, context, false)
	if err != nil {
		tool.ResponseFail(context, -1, fmt.Sprintf("文件上传失败 %s", err), nil)
		return
	}

	tool.ResponseSuccess(context, fmt.Sprintf("%s", localU))

	defer func() {
		_ = os.Remove(path)
	}()
}

func UploadFaceOutsea(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[文件上传] 文件上传获取失败 %s", err))
		return
	}

	localU, path, err := uploadBusiness.UploadWithFile(file, context, true)
	if err != nil {
		tool.ResponseFail(context, -1, fmt.Sprintf("文件上传失败 %s", err), nil)
		return
	}

	tool.ResponseSuccess(context, fmt.Sprintf("%s", localU))

	defer func() {
		_ = os.Remove(path)
	}()
}

func UploadVideo(context *gin.Context) {

	file, err := context.FormFile("file")
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[文件上传] 文件上传获取失败 %s", err))
		return
	}

	localU, path, err := uploadBusiness.UploadWithFile(file, context, false)
	if err != nil {
		tool.ResponseFail(context, -1, fmt.Sprintf("文件上传失败 %s", err), nil)
		return
	}

	tool.ResponseSuccess(context, fmt.Sprintf("%s", localU))

	defer func() {
		_ = os.Remove(path)
	}()
}
