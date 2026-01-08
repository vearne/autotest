-- Utility functions for testing

function formatTimestamp(ts)
    return os.date("%Y-%m-%d %H:%M:%S", ts)
end

function validateEmail(email)
    return string.match(email, "[^@]+@[^@]+") ~= nil
end

function generateUUID()
    local template ='xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    return string.gsub(template, '[xy]', function (c)
        local v = (c == 'x') and math.random(0, 0xf) or math.random(8, 0xb)
        return string.format('%x', v)
    end)
end

function deepEquals(t1, t2)
    local type1, type2 = type(t1), type(t2)
    if type1 ~= type2 then return false end
    if type1 ~= 'table' then return t1 == t2 end

    for k, v in pairs(t1) do
        if not deepEquals(v, t2[k]) then return false end
    end
    for k, v in pairs(t2) do
        if not deepEquals(v, t1[k]) then return false end
    end
    return true
end
