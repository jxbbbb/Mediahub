package Controller

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mediahub/pkg/config"
	"mediahub/pkg/log"
	"mediahub/pkg/storage"
	"mediahub/pkg/zerror"
	"mediahub/services"
	"mediahub/services/shorturl"
	proto2 "mediahub/services/shorturl/proto"
	"net/http"
	"path"
)

type Controller struct {
	sf     storage.StorageFactory
	log    log.ILogger
	config *config.Config
}

func NewController(sf storage.StorageFactory, log log.ILogger, config *config.Config) *Controller {
	return &Controller{sf: sf, log: log, config: config}
}

func (c *Controller) Upload(ctx *gin.Context) {
	userId := ctx.GetInt64("user_id")
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		c.log.Error(zerror.NewByErr(err))
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.log.Error(zerror.NewByErr(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		c.log.Error(zerror.NewByErr(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	//图片格式校验
	if !isImage(io.NopCloser(bytes.NewReader(content))) {
		err = zerror.NewByMsg("仅支持jpg，png，gif")
		c.log.Error(zerror.NewByErr(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	//生成唯一文件名，避免重复上传
	md5Digest := calMD5Digest(content)
	filename := fmt.Sprintf("%x%s", md5Digest, path.Ext(fileHeader.Filename))
	filePath := "/public/" + filename
	if userId != 0 {
		filePath = fmt.Sprintf("/%d/%s", userId, filename)
	}

	s := c.sf.CreateStorage()
	url, err := s.Upload(bytes.NewReader(content), md5Digest, filePath)
	if err != nil {
		c.log.Error(zerror.NewByErr(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	shortPool := shorturl.NewShortUrlClientPool()
	clientConn := shortPool.Get()
	defer shortPool.Put(clientConn)
	//生成短链接

	client := proto2.NewShortUrlClient(clientConn)
	in := &proto2.Url{
		Url:    url,
		UserID: userId,
	}
	outGoingCtx := context.Background()
	outGoingCtx = services.AppendBearerToken(outGoingCtx, c.config.DependOn.ShortUrl.AccessToken)
	outUrl, err := client.GetShortUrl(outGoingCtx, in)
	if err != nil {
		c.log.Error(zerror.NewByErr(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": outUrl.Url,
	})
	return

}

func isImage(r io.Reader) bool {
	_, _, err := image.Decode(r)
	if err != nil {
		return false
	}
	return true
}

func calMD5Digest(msg []byte) []byte {
	m := md5.New()
	m.Write(msg)
	bs := m.Sum(nil)
	return bs
}

//错误总结
//错误点：io.Reader 被多次读取导致上传失败
//
//错误原因
//io.Reader 只能被读取一次，io.ReadAll(file) 读取完数据后，file 变为空，后续 Upload(file, ...) 传入的是空数据流。
//isImage(io.NopCloser(bytes.NewReader(content))) 读取 content 进行图片格式校验，但 file 没有重新赋值，导致 Upload 方法无法读取文件内容。
//解决方案
//方法 1（推荐）：使用 bytes.NewReader(content) 重新创建 io.Reader，保证 Upload 能正确读取数据。
//方法 2：使用 io.TeeReader 让 isImage 和 Upload 共享 file，避免 file 被提前读空。
