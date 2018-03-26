#!/usr/bin/env python
# @Author  : pengyun

import redis

redis_client = redis.StrictRedis(host='localhost', port=6379, db=0)

# 斗鱼的相关配置
Config = {
    "douyu":{
        'minnum': 100000,
        'maxlink': 60,
        'taskname': 'douyu|task',
    }
}