local thread_counter = 0

function setup(thread)
    thread:set("id", thread_counter)
    thread_counter = thread_counter + 1
end

function init(args)
    wrk.method = "POST"
    wrk.headers["Content-Type"] = "application/json"
    req_counter = 0
end

function request()
    req_counter = req_counter + 1
    local body = string.format('{"url":"https://example.com/item/%d/%d"}', id, req_counter)
    return wrk.format("POST", "/", nil, body)
end
