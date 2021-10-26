
--get all of the jobs with an expired "score"...
local val = redis.call('zrangebyscore',KEYS[1],'-inf',ARGV[1])
-- If wh have values in the array,we will remove them from the first queue
-- and add them onto the destination queue in chunks of 100,which moves
-- all of the appropriate jobs onto the destination queue very safely.
if (next(val) ~= nil) then
    redis.call('zremrangebyrank',KEYS[1],0,#val - 1)

    for i = 1, #val,100 do
        redis.call('rpush',KEYS[2],unpack(val,i,math.min(i+99,#val)))
    end
end

return val