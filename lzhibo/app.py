#!/usr/bin/python3
# -*- coding: utf-8 -*-
# @Created on    : 18/3/20 下午8:17
# @Author  : zpy
# @Software: PyCharm

from flask import Flask
from tasks import Douyu

app = Flask(__name__)

@app.route('/task')
def task():
    print(Douyu().gettasks())
    return ""

if __name__ == '__main__':
    app.run(debug=True)
