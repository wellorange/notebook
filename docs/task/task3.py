import numpy as np
import string
import time
import itertools

    
file_name = 'input1.txt'

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
    #    result=[]
    #    for index in include:
    #         point= np.argwhere(vector==index)[0]
    #         result.append([point[0],point[1]])
       return include




def Policepoint2(dataset):
       n=dataset[0]                  #  矩阵常量
       p=dataset[1]
       vector = np.arange(0, n*n, 1) # 正常的矩阵初始化
       vector = vector.reshape(n, n) # 矩阵初始化
       possibly={}
       possiable=[]
       for i in range(n):
            for ii in range(n):
              possibly[vector[i,ii]]=Policepoint([i,ii],vector)
      

    #    key=0
    #    value=possibly[0]
    #    print(value)
    #    k=[]
    #    k.append(key)

    #    def xuxu(key,value):
           

    #    for p in range(p-1):
    #        xuxu()
    # #    for i in itertools.combinations(key, p-1):
    # #        pss




    #        posi=possibly[i[0]]
    #        po=[]
    #        for ii in i:
    #           po.append(ii)
    #           posi=list(set(posi).intersection(set(possibly[ii])))
    #        if len(posi)!=0: 
    #          possiable.append([po,posi])
       
    #    for ii in possiable:
    #        print(ii)

       if p==2:
          for posi in possibly:
             possiable.append([[posi],possibly[posi]])
          return possiable
       #p==3 时候
       for index in possibly:
            p1=index
            for p2 in possibly[index]:
                posi=np.intersect1d(possibly[p1], possibly[p2])
                possiable.append([[p1,p2],posi])
   
       #p>3才会执行
       for indx in possiable:
           pass
        #    print(indx)

       if p>3:
        for police in range(p-3):
             sothat=[]
             for index in possiable:
                 if len(index[1])==0:
                     continue
                 sub1=index[0]
                 for p2 in index[1]:
                     su=sub1.copy()
                     su.append(p2)
                     posi=possibly[su[0]]
                     for indx in su:
                         posi=np.intersect1d(possibly[indx], posi)
                     sothat.append([su,posi])
                     sothat2=[]
                     #print(su)
                     for i,ii in enumerate(sothat):
                        if len(ii[1])!=0:
                           sothat2.append(sothat[i])   
             possiable=sothat2
       # 坐标转换          
       export=[]
       for po in possiable:
           pp1=[]
           pp2=[]
           for po1 in po[0]:
               point= np.argwhere(vector==po1)[0]
               pp1.append([point[0],point[1]])
           for po2 in po[1]:
               point= np.argwhere(vector==po2)[0]
               pp2.append([point[0],point[1]])
           export.append([pp1,pp2])


       for ii in export:
           pass
            # print(ii)
       return export






# 小车的总次数计算
def Calculate(n,list):
    matrax=np.zeros((n,n))
    for kv in list:
        matrax[kv[0]][kv[1]]+=1
    return matrax
numberM=Calculate(dataset[0],dataset[3])

def Wrappolice(dataset):
     total=0
     t1=0
     t2=0
     n=dataset[0]
     possibly=[]   # 这里是另外的警察的位置[0,1,2] 0,1是第一个警察的位置,2是第二个警察的位置

     # 一个警察的时候
     if dataset[1]==1:
         for i in range(n):
             for ii in range(n):
                 if total<numberM[i][ii]:
                     total=numberM[i][ii]
     else:
       #当不是一个警察的时候
      pos=Policepoint2(dataset)
      
      for p in pos:
          t1=0
          t2=0
          c1=[]
          c2=[]
          for p1 in p[0]:
              c1=p[0]
              t1+=numberM[p1[0]][p1[1]]
          for p2 in p[1]:
              c2=p2
              t2=numberM[p2[0]][p2[1]]
              if total<(t1+t2):
                #   print(c1,c2)
                  total=t1+t2
     return total  

     



start = time.time()
total=Wrappolice(dataset)
end = time.time()

print (end-start)
print(total)
