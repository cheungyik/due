package redis

// 释放锁
const releaseScript = `
	local val = redis.call('GET', KEYS[1])

	if val == '' then
		return {'OK'}
	end

	if val ~= ARGV[1] then
		return {'NO'}
	end

	redis.call('DEL', KEYS[1])

	return {'OK'}
`

// 续租锁
const renewalScript = `
	
`
