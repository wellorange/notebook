import numpy as np
import string
import time
import os
from collections import defaultdict

def tree(): 
   return defaultdict(tree)    


file_name = 'input3.txt'

# dataset 是从input读取的数据集合
dataset = []
datapoint=[]

for index,line in enumerate(open(file_name)):
       if index<3:
           # 0 矩阵大小
           # 1 警察数量
           # 2小车数量
           temp=int(line.split()[0])
           dataset.append(temp)
       else:
           # 小车坐标
           linee=line.split(",")
           x=int(linee[0])
           y=int(linee[1])
           lin=[x,y]
           datapoint.append(lin)
       
dataset.append(datapoint)
    

# monitor是标杆坐标, 作用就是在vector返回与monitor构成皇后结构的所有坐标
def Policepoint(monitor,vector):
       exclude=[]
       vector2=np.rot90(vector)      #矩阵转置
       monivalue=vector[monitor[0],monitor[1]]
       monitor2 = np.argwhere(vector2 ==monivalue)[0]
       # 获取右对角线
       excu1=vector2.diagonal(offset=(monitor2[1]-monitor2[0]))
       # 获取moniter的左对角线
       excu2=vector.diagonal(offset=(monitor[1]-monitor[0]))
       # 获取moniter的同行
       excu3=vector[monitor[0],:]
       # 获取moniter的同列
       excu4=vector[:,monitor[1]]
       exclude=np.hstack((excu1,excu2,excu3,excu4))
       exclude=np.unique(exclude)
       include=np.setdiff1d(vector, exclude, assume_unique=True)
       return include









# 小车的总次数计算
def Calculate(n,list):
    matrax=np.zeros((n,n))
    for kv in list:
        matrax[kv[0]][kv[1]]+=1
    return matrax

numberM=Calculate(dataset[0],dataset[3])


global toal,ifelse,history
toal=0
w=[]
ifelse=[]

for ii in range(dataset[0]):
  for i in range(dataset[0]):
      w.append(numberM[ii][i])


ww=w.copy()
maxx=[]

for i in range(3):
  inn=ww.index(max(ww))
  maxx.append(inn)
  del ww[inn]


def iflist(list1,list2,list3):
    if  (len(set(list1) & set(list3))!=0):
        return True
    if  (len(set(list2) & set(list3))!=0):
        return True
    return False



def Policepoint2(dataset):
       global toal,maxx
       n=dataset[0]                  #  矩阵常量
       p=dataset[1]
       vector = np.arange(0, n*n, 1) # 正常的矩阵初始化
       vector = vector.reshape(n, n) # 矩阵初始化
       possibly={}
       for i in range(n):
            for ii in range(n):
              possibly[vector[i,ii]]=Policepoint([i,ii],vector)     

       def recall(key,keyvalue,p):
           global toal
           if len(key)==p:
            # n=15的时候打印吧,因为你会感觉程序没有运行 打印才发现程序在运行
            #   print(key,keyvalue)
               temp=0
               for ii in key:
                    temp+=w[ii]
               if temp>toal:
                   toal=temp
               return
           if len(keyvalue)==0:
               return
           for k,v in  enumerate(keyvalue):
                # 去除重复的去除
                if v>key[len(key)-1]: 
                    kk=key.copy()
                    kk.append(v)
                    posi=np.intersect1d(keyvalue, possibly[v])
                    if iflist(kk,posi,maxx): #概率在去一次   
                        recall(kk,posi,p)

       # one a police
       if dataset[1]==1:
         for i in range(n):
             for ii in range(n):
                 if toal<numberM[i][ii]:
                     toal=numberM[i][ii]
       else:
           for index in possibly:
              flag=0
              if dataset[1]==dataset[0]:
                  flag=dataset[0]
              else:
                  falg=(dataset[0]-dataset[1])*dataset[0]

              if index<flag:
                #   print(index)
                  value=possibly[index]
                  kk=[index]
                  recall(kk,value,p)
              else:
                  break
          
          

# test code       
start = time.time()
Policepoint2(dataset)
end = time.time()
print (end-start)
print(toal)
