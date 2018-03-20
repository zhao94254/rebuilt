#!/usr/bin/python3
# -*- coding: utf-8 -*-
# @Created on    : 18/3/20 下午7:25
# @Author  : zpy
# @Software: PyCharm
import datetime
import time

class Test:

    def __init__(self, b):
        self.a = datetime.datetime.now()
        self.b = b



a = Test(12)
time.sleep(1)
c = Test(321)
time.sleep(1)
b = Test(242)

s = [a,b,c]
print([i.b for i in s])

s = sorted(s, key=lambda x:x.a)[::-1]


print([i.b for i in s])
