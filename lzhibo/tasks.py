#!/usr/bin/env python
# @Author  : pengyun

from config import *
import requests
import json

# directory: link image onlinenum directoryname

# streamer : link image onlinenum streamername

# config

# Config = {
#     "douyu":{
#         'minnum': 100000,
#         'maxlink': 100,
#         'taskname': 'douyu|task',
#     }
# }


class Douyu:
    """ 获取斗鱼的信息， 将任务分发下去"""

    def __init__(self):
        self.directory_info = {}
        self.streammer_info = {}
        self.baselink = "http://api.douyutv.com/api/v1/live/{}"
        self.load()
        self.parse_config()



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
            self.parse_directory()

    def getonline(self):
        """ 解析用户数据"""
        for i in self.base_data:
            res = requests.get(self.baselink.format(i['short_name'])).json()['data']
            self.directory_info[i['short_name']]['num'] = sum([i['online'] for i in res])
            for j in res:
                if j['online'] > self.minnum: # 这里通过config来配置
                    self.streammer_info[j['room_id']] = {
                        'roomid': j['room_id'],
                        'online': j['online'],
                        'nickname': j['nickname'],
                        'fans': j['fans'],
                        'image': j['avatar_mid'],
                        'pindao': i
                    }
                yield j['room_id']

    def parse_directory(self):
        """ 解析频道相关"""
        for i in self.base_data:
            self.directory_info[i['short_name']] = {
                'img': i['game_icon'],
                'gname': i['game_name'],
                'link': i['game_url'],
                'num':0,
            }

    def parse_streammer(self):
        pass

    def gettasks(self):
        redis_client.delete(self.taskname)
        res = []
        for i, j in enumerate(self.getonline()):
            if i >= self.maxlink:
                break
            else:
                res.append(j)
        res = '|'.join(res)
        redis_client.set(self.taskname, res)
        return redis_client.get(self.taskname) # decode 的时候通过分割 | 来实现






if __name__ == '__main__':
    d = Douyu()
    print(d.gettasks())