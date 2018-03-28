#!/usr/bin/python3
# -*- coding: utf-8 -*-
# @Created on    : 18/3/20 下午8:17
# @Author  : zpy
# @Software: PyCharm

from flask import Flask, jsonify
from tasks import Douyu, redis_client


# Redis key
# douyu|basedata   douyu|directorys(short_name)   douyu|task(roomid)

# all shortname

# all user info


app = Flask(__name__)

@app.route('/task')
def task():
    res = Douyu().gettasks()
    return res

@app.route('/g/<string:key>')
def data(key):
    return redis_client.get(key)

@app.route("/keys")
def keys():
    return str(redis_client.keys())

@app.route("/onerank")
def rank():
    data = {}
    _keys = [i.decode() for i in redis_client.keys() if i.decode().startswith("one|")]
    for k in _keys:
        data[k] = redis_client.get(k)
    return jsonify(data)

if __name__ == '__main__':
    app.run(host="127.0.0.1", debug=True)
