package cos

import (
	"context"
	"encoding/base64"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"mediahub/pkg/storage"
	"net/http"
	url2 "net/url"
	"path"
	"strings"
)

type cosStorageFactory struct {
	bucketUrl string
	secretId  string
	secretKey string
	cdnDomain string
}

func NewCosStorageFactory(bucketUrl, secretId, secretKey, cdnDomain string) storage.StorageFactory {
	return &cosStorageFactory{
		bucketUrl: bucketUrl,
		secretId:  secretId,
		secretKey: secretKey,
		cdnDomain: cdnDomain,
	}
}

func (f *cosStorageFactory) CreateStorage() storage.Storage {
	return NewCos(f.bucketUrl, f.secretId, f.secretKey, f.cdnDomain)
}

type cosStorage struct {
	bucketUrl string
	secretId  string
	secretKey string
	cdnDomain string
}

func NewCos(bucketUrl, secretId, secretKey, cdnDomain string) storage.Storage {
	return &cosStorage{
		bucketUrl: bucketUrl,
		secretId:  secretId,
		secretKey: secretKey,
		cdnDomain: cdnDomain,
	}
}

func (s *cosStorage) Upload(r io.Reader, md5Digest []byte, dstPath string) (url string, err error) {
	u, err := url2.Parse(s.bucketUrl)
	if err != nil {
		return "", err
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  s.secretId,
			SecretKey: s.secretKey,
		},
	})

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: s.getContentType(dstPath),
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{},
	}
	if len(md5Digest) != 0 {
		opt.ObjectPutHeaderOptions.ContentMD5 = base64.StdEncoding.EncodeToString(md5Digest)
	}
	_, err = client.Object.Put(context.Background(), dstPath, r, opt)
	if err != nil {
		return "", err
	}
	url = s.bucketUrl + dstPath
	if s.cdnDomain != "" {
		url = s.cdnDomain + dstPath
	}
	return url, err
}

func (s *cosStorage) getContentType(dstFile string) string {
	ext := strings.Trim(path.Ext(dstFile), ".")
	if ext == "jpg" {
		ext = "jpeg"
	}
	return "image/" + ext
}
