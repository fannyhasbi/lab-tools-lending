userid = os.getenv("USER_ID")
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

counter = 0

request = function()
  local str = ""

  if (counter == 0) 
  then
    str = string.format([[{
      "message": {
        "message_id": 1,
        "from": {
          "id": %d,
          "first_name": "Test",
          "last_name": "TestLast",
          "username": "test"
        },
        "text": "/cek",
        "chat": {
          "id": %d,
          "type": "private"
        }
      }
    }]], userid, userid)

  elseif (counter == 1)
  then
    str = string.format([[{
      "message": {
        "message_id": 1,
        "from": {
          "id": %d,
          "first_name": "Test",
          "last_name": "TestLast",
          "username": "test"
        },
        "text": "/cek 1",
        "chat": {
          "id": %d,
          "type": "private"
        }
      }
    }]], userid, userid)

  elseif (counter == 2)
  then
    str = string.format([[{
      "message": {
        "message_id": 1,
        "from": {
          "id": %d,
          "first_name": "Test",
          "last_name": "TestLast",
          "username": "test"
        },
        "text": "/cek 1 foto",
        "chat": {
          "id": %d,
          "type": "private"
        }
      }
    }]], userid, userid)
    counter = -1
  end

  counter = counter + 1
  return wrk.format(nil, nil, nil, str)
end
