FROM golang:1.12-alpine as builder
RUN apk add git
COPY . /go/src/shuStudentCourse
ENV GO111MODULE on
WORKDIR /go/src/shuStudentCourse
RUN go get && go build

FROM alpine
MAINTAINER longfangsong@icloud.com
COPY --from=builder /go/src/shuStudentCourse/shuStudentCourse /
WORKDIR /
CMD ./shuStudentCourse
ENV PORT 8000
EXPOSE 8000