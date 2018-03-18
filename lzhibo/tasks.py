#!/usr/bin/env python
# @Author  : pengyun

from config import *
import requests
import json

# directory: link image onlinenum directoryname

# streamer : link image onlinenum streamername




class Douyu:
    """ 获取斗鱼的信息， 将任务分发下去"""

    def __init__(self):
        self.parse_config()
        self.load()
        self.baselink = "http://api.douyutv.com/api/v1/live/{}"

    def parse_config(self):
        for i,j in Config["douyu"].items():
            setattr(self, i, j)

    def load(self):
        """ 这个数据只需要获取一次"""
        if redis_client.get("douyu|basedata") is None:
            data = requests.get('http://open.douyucdn.cn/api/RoomApi/game').json()
            self.base_data = data['data']
            redis_client.set("douyu|basedata", json.dumps(data))
        else:
            self.base_data = json.loads(redis_client.get("douyu|basedata").decode())['data']

    def getonline(self):
        for i in self.base_data:
            res = requests.get(self.baselink.format(i['short_name'])).json()['data']

    def parse_directory(self):
        pass

    def parse_streammer(self):
        pass

    def gettasks(self):
        pass

if __name__ == '__main__':
    Douyu().getonline()
    print(Douyu().base_data)