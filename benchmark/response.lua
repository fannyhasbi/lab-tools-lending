userid = os.getenv("USER_ID")
groupid = os.getenv("GROUP_ID")
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "peminjaman_testing_bot"

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
    "text": "/tanggapi pinjam 46",
    "chat": {
      "id": %d,
      "type": "group"
    }
  }
}]], userid, groupid)
