
# 错误总结
## 错误点：io.Reader 被多次读取导致上传失败

1. 错误原因：io.Reader 只能被读取一次，io.ReadAll(file) 读取完数据后，file 变为空，后续 Upload(file, ...) 传入的是空数据流。
   isImage(io.NopCloser(bytes.NewReader(content))) 读取 content 进行图片格式校验，但 file 没有重新赋值，导致 Upload 方法无法读取文件内容。
   解决方案 ：使用 bytes.NewReader(content) 重新创建 io.Reader，保证 Upload 能正确读取数据。
2. 错误原因： 跨域问题，前后端不在同一个域名上，前端访问会先调用option方法，成功后才会调用其他方法，
   r := gin.Default()
   api := r.Group("/api")
   api.Use(middleware.Cors())
   但如果在api下跨域会不通过，需要将跨域设置在根路由下，即
   r.Use(middleware.Cors())

## shorturl/shorturl-server/server/server.go //根据长链接查询数据库，记录是否存在（可以考虑用缓存，长链做key，或者短链做key，但功能主要是添加，查询业务较小，暂不考虑）
