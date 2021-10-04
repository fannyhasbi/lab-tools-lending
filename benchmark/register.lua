userid = os.getenv("USER_ID")
wrk.method = "POST"
wrk.body   = string.format([[{
  "message": {
    "message_id": 1,
    "from": {
      "id": %d,
      "first_name": "Test",
      "last_name": "TestLast",
      "username": "test"
    },
    "text": "/registrasi",
    "chat": {
      "id": %d,
      "type": "private"
    }
  }
}]], userid, userid)
wrk.headers["Content-Type"] = "application/json"