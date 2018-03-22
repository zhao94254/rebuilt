#!/usr/bin/python3
# -*- coding: utf-8 -*-
# @Created on    : 18/3/20 下午8:17
# @Author  : zpy
# @Software: PyCharm

from flask import Flask
from tasks import Douyu, redis_client


# Redis key
# douyu|basedata   douyu|directorys(short_name)   douyu|task(roomid)

# all shortname

# all user info


app = Flask(__name__)

@app.route('/task')
def task():
    print(Douyu().gettasks())
    return ""

@app.route('/<string:key>')
def data(key):
    return redis_client.get(key)



if __name__ == '__main__':
    app.run(debug=True)
