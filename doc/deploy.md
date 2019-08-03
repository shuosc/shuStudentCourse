# 部署
本项目已经打包成[docker镜像](https://hub.docker.com/r/shuosc/shu-student-course/)。
## 支持服务
### postgresql数据库
migration文件位于本repo的 [migration](https://github.com/shuosc/shuStudentCourse/tree/master/migration) 目录中。

建议使用 [golang-migrate](https://github.com/golang-migrate/migrate) 来进行 migrate。
```shell
migrate -source github://[你的Github用户名]:[你的Github Access Token]@shuosc/shuStudentCourse/migration -database [你的postgrsql数据库url] up
```

## 服务本身
### 环境变量
- `PORT`: 服务端口
- `DB_ADDRESS`: 数据库url
- `JWT_SECRET`: jwt密钥
- `PROXY_ADDRESS`: 访问学校选课网站代理服务地址
- `COURSE_INFO_URL`: 课程信息服务地址
- `COURSE_SELECTION_ADDRESS_URL`: 选课网站地址信息服务地址

### k8s
在k8s下使用如下yaml，假设
- `JWT_SECRET`由k8s secret给出
- 数据库服务器在`shu-student-course-postgres-svc`
- 代理服务地址在`shu-course-proxy-svc`
- 课程信息地址在`shu-course-svc/course`
- 选课网站地址信息服务地址在`shu-course-svc/course-selection-url`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: shu-course
spec:
  selector:
    matchLabels:
      run: shu-student-course
  replicas: 1
  template:
    metadata:
      labels:
        run: shu-student-course
    spec:
      containers:
      - name: shu-student-course
        image: shuosc/shu-student-course
        env:
        - name: PORT
          value: "8000"
        - name: DB_ADDRESS
          value: "postgres://shuosc@shu-student-course-postgres-svc:5432/shu-student-course?sslmode=disable"
        - name: PROXY_ADDRESS
          value: "http://shu-course-proxy-svc"
        - name: COURSE_INFO_URL
          value: "http://shu-course-svc/course"
        - name: COURSE_SELECTION_ADDRESS_URL
          value: "http://shu-course-svc/course-selection-url"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: shuosc-secret
              key: JWT_SECRET
        ports:
        - containerPort: 8000
---
apiVersion: v1
kind: Service
metadata:
  name: shu-student-course-svc
spec:
  selector:
     run: shu-student-course
  ports:
  - protocol: TCP
    port: 8000
    targetPort: 8000
```