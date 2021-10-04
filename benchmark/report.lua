userid = os.getenv("USER_ID")
groupid = os.getenv("GROUP_ID")
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
    "text": "/laporan pinjam 2021-9",
    "chat": {
      "id": %d,
      "type": "group"
    }
  }
}]], userid, groupid)
wrk.headers["Content-Type"] = "application/json"