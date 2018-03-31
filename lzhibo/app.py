#!/usr/bin/python3
# -*- coding: utf-8 -*-
# @Created on    : 18/3/20 下午8:17
# @Author  : zpy
# @Software: PyCharm

from flask import Flask, jsonify
from tasks import Douyu, redis_client
import json

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
    data = sorted_data("one")
    return jsonify(data)


@app.route("/fiverank")
def frank():
    data = sorted_data("five")
    return jsonify(data)

@app.route("/halfrank")
def hrank():
    data = sorted_data("half")
    return jsonify(data)


@app.route("/directors")
def prank():
    data = pdata()
    return jsonify(data)


def sorted_data(long):
    _keys = redis_client.get("douyu|task").decode()
    data = []
    tmp = {}
    for k in _keys.split('|'):
        tmp["user"] = json.loads(redis_client.get(k))
        hots = redis_client.get("{}|{}".format(long, k))
        if hots is None:
            hots = 0
        tmp["user"]['hots'] = int(hots)
        data.append(tmp)
        tmp = {}
    data = sorted(data, key=lambda x: -x["user"]["hots"])
    return data


def pdata():
    _keys = redis_client.get("douyu|directorys").decode()
    data = []
    tmp = {}
    for k in _keys.split('|'):
        value = redis_client.get(k)
        if value is None:
            continue
        tmp["directory"] = json.loads(value)
        data.append(tmp)
        tmp = {}
    data = sorted(data, key=lambda x: -x["directory"]["hots"])
    return data


if __name__ == '__main__':
    app.run(host="0.0.0.0", debug=True)
