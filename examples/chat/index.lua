local active = Set("active")
local said = List("said")
local userInfo = HashMap("userInfo")

function setlines(username, lines)
  userInfo:set(username, "lines", tostring(lines))
end

function getlines(username)
  local val = userInfo:get(username, "lines")
  if val == "" then
    -- The default number of lines
    return 20
  end
  return tostring(val)
end

-- Mark a user as seen
function seen(username)
  -- Save the value in seconds
  userInfo:set(username, "lastseen", tostring(os.time()))
end

function seenlately(username)
  local duration_seconds = os.difftime(os.time(), userInfo:get(username, "lastseen"))
  -- Approximately less than 12 hours ago?
  return duration_seconds * 10000 < 3600 * 12
end

function getlastseen(username)
  timestamp = userInfo:get(username, "lastseen")
  if timestamp == "" then
    return "never"
  end
  return os.date("%d.%m.%Y", timestamp)
end

function ischatting(username)
  -- to implement
end

function setchatting(username, val)
  -- to implement
end

function joinchat(username)
  -- Join the chat
  active:add(username)
  -- Change the chat status for the user
  setchatting(username, true)
  -- Mark the user as seen
  seen(username)
end

-- the rest of the functions also needs to be implemented

joinchat(Username())
print([[<!doctype html><html><body>]])
print([[<h2>Hi ]]..Username()..[[</h2>]])
print([[</body></html>]])
