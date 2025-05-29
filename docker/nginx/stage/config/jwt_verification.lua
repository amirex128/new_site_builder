local jwt = require "resty.jwt"
ngx.req.set_header("X-Gateway-Authorization", "f50405b0f722f4b68c176066b145b0a4f130df550aafc1c8d7aacda062268fd4")
local jwt_secret = "3e44e4aaa0f22573dbfa5008bd868165e9190432316dc15467732d8037f026af05211477f62e9f16b98c1e2beace0b1489ed9aaf9077bf3f061600c7acdecfa2"
local jwt_issuer = "LionHeartIssuer"
local jwt_audience = "LionHeartAudience"
local jwt_validation_skip = 0
local request_uri = ngx.var.request_uri

if request_uri == "/" or request_uri:match("^/[^/]+/Free/") then
    jwt_validation_skip = 1
end

if jwt_validation_skip == 1 then
    return
end

local token = ngx.var.http_authorization
if token and token:find("Bearer") then
    token = token:sub(8)
else
    ngx.exit(ngx.HTTP_UNAUTHORIZED)
end

local jwt_obj = jwt:verify(jwt_secret, token)
if not jwt_obj["verified"] then
    ngx.exit(ngx.HTTP_UNAUTHORIZED)
end

local payload = jwt_obj["payload"]
if payload["iss"] ~= jwt_issuer or payload["aud"] ~= jwt_audience then
    ngx.exit(ngx.HTTP_UNAUTHORIZED)
end

ngx.req.set_header("X-Type", payload["type"])
ngx.req.set_header("X-Roles", payload["roles"])

if payload["user_id"] then
    ngx.req.set_header("X-User-Id", payload["user_id"])
    ngx.req.set_header("X-Site-Ids", payload["site_ids"])
elseif payload["customer_id"] then
    ngx.req.set_header("X-Customer-Id", payload["customer_id"])
end