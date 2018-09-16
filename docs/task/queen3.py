#-*- coding:utf-8 -*-
import time
def find(row, ld, rd):
   
    global n, upperlim,count
    if row == upperlim: #当row == upperlim时，所有列已有棋子，即产生一个解
        count += 1
    else:
        # row ld rd 或运算会返回所有不能放的位置 取反就是能放的位置 没有能放的位置就是皇后放完了
        pos = upperlim & (~(row | ld | rd)) # 所以返回的是能放的位置
        while pos: #有位置可放
            p = pos & (~pos + 1) #取出pos中最右边那个1也就是取最右边那个可以放的位置
            pos = pos - p #将取出的1（合法的位置）从pos中去除
            find(row | p, (ld | p) << 1, (rd | p) >> 1) #递归查找


def main():
    global n,upperlim,count
    n = 4
    upperlim = (1 << n) - 1 #表示N个1
    count = 0 
    find(0, 0 ,0) 
    print(count)
    print(upperlim)

start = time.time()
main()
end = time.time()
print (end-start)