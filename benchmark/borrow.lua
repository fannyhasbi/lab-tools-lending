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
        "text": "/pinjam 1",
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
        "text": "1",
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
        "text": "30",
        "chat": {
          "id": %d,
          "type": "private"
        }
      }
    }]], userid, userid)

  elseif (counter == 3)
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
        "text": "test reason",
        "chat": {
          "id": %d,
          "type": "private"
        }
      }
    }]], userid, userid)

  elseif(counter == 4)
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
        "text": "no",
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
