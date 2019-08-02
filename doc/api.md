# API Reference

## web api

- `GET /ping`

  检查服务是否可用，应该直接返回`pong`。

- `GET /student-courses?semester_id=[学期id]`

  获得学生在某一学期选的课。
  
  学生id在JWT中给出。
