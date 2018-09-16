# -*- coding:utf-8 -*-
import random
import time
# 冲突检测，定义state元组为皇后的位置，nextX为下一个皇后的横坐标（即所在列）
# 如state[1] = 2表示，皇后的位置处在第二行第三列。
def conflict(state,nextX):
	nextY = len(state)
	for i in range(nextY):
		# 如果下一个皇后位置与每个放置的皇后位置在同一列或者在同一对角线上，则表示冲突
		if abs(state[i] - nextX) in (0,nextY - i):
			return True
	return False
# 生成器递归生成皇后的位置
def queens(num,state=()):
	for pos in range(num):
		if not conflict(state,pos):
			# 递归的出口，产生皇后的位置
			if len(state) == num - 1:
				yield (pos,)
			else:
				# 将当前的皇后位置添加到state中并传递给下一个皇后
				for result in queens(num,state + (pos,)):
					yield (pos,) + result
# 使用棋盘布局直观显示结果，其中Q表示皇后位置
def prettyprint(solution):
	def line(pos,length = len(solution)):
		return '. ' * pos + 'Q ' + '. ' * (length - pos -1)
	for pos in solution:
		print(line(pos))
# 随机的显示一种结果start = time.time()
start = time.time()
end = time.time()
prettyprint(random.choice(list(queens(15))))
print (end-start)


