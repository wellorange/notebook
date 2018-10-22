import numpy as np
import string
import time
import os
from collections import defaultdict

def tree(): 
   return defaultdict(tree)    


file_name = 'case/input25.txt'
# dataset 是从input读取的数据集合
token=("bed","park","la_num","la","my_num","my","list_num","list")
cur_token=token[0]
dataset = []

number=0
numbert=1
temp=[]

  
exits=[]     # 双方剔除过的
mykey=[]
youkey=[]
youexits=[]  # 你剔除过的
myexits=[]   # 我剔除过的
for index,line in enumerate(open(file_name)):
    if cur_token==token[0]:
        temp.append(int(line))
        dataset.append(temp)
        temp=[]
        cur_token=token[1]
    elif cur_token==token[1]:
        temp.append(int(line))
        dataset.append(temp)
        temp=[]
        cur_token=token[2]
    elif cur_token==token[2]:
        temp.append(int(line))
        dataset.append(temp)
        temp=[]
        number=int(line)
        cur_token=token[3]
    elif cur_token==token[3]:
         if numbert==number:
            cur_token=token[4]
            temp.append(line.split()[0])
            dataset.append(temp)
            exits.append(line.split()[0])
            youkey.append(line.split()[0])   
            number=0
            numbert=1
            temp=[]
         else:
              temp.append(line.split()[0])
              exits.append(line.split()[0])
              youkey.append(line.split()[0])
              numbert+=1
    elif cur_token==token[4]:
        temp.append(int(line))
        dataset.append(temp)
        temp=[]
        number=int(line)
        cur_token=token[5]
    elif cur_token==token[5]:
         if numbert==number:
            cur_token=token[6]
            number=0
            numbert=1
            temp.append(line.split()[0])
            exits.append(line.split()[0])
            mykey.append(line.split()[0])
            dataset.append(temp)
            temp=[]
         else:
              temp.append(line.split()[0])
              exits.append(line.split()[0])
              mykey.append(line.split()[0])
              numbert+=1
    elif cur_token==token[6]:
          temp.append(int(line))
          dataset.append(temp)
          temp=[]
          number=int(line)
          cur_token=token[7]
    elif cur_token==token[7]:
         if numbert==number:
            temp.append(line.split()[0])
            dataset.append(temp)
            if line.split()[0][0:5] in mykey:
                  myexits.append(line.split()[0])
            if line.split()[0][0:5] in youkey:
                  youexits.append(line.split()[0])

            temp=[]
            cur_token="null"
            number=0
            numbert=1
         else:
              temp.append(line.split()[0])

              if line.split()[0][0:5] in mykey:
                  myexits.append(line.split()[0])
              if line.split()[0][0:5] in youkey:
                  youexits.append(line.split()[0])
              numbert+=1












my=[]
you=[]
our=[]
record={}

for line in dataset[7]:
    a=0
    for i in line[13:21]:
           a+=int(i)
    record[line[0:5]]=a
    if line[0:5] not in exits:
       if (line[10:11]=="N" and line[11:12]=="Y" and line[12:13]=="Y") and (line[5:6]!="F" or int(line[6:9])<=17 or line[9:10]!="N"):
          my.append(line)
          

       if (line[5:6]=="F" and int(line[6:9])>17 and line[9:10]=="N") and (line[10:11]!="N" or line[11:12]!="Y" or line[12:13]!="Y") :
          you.append(line)

       if line[5:6]=="F" and int(line[6:9])>17 and line[9:10]=="N" and line[10:11]=="N" and line[11:12]=="Y" and line[12:13]=="Y":
          our.append(line)



youexitnumber=len(youexits)
youexittotal=0
for yii in youexits:
    for i in yii[13:21]:
        youexittotal+=int(i)


myexitesnumber=len(myexits)
myexittotal=0
for mii in myexits:
    for i in mii[13:21]:
        myexittotal+=int(i)
haha=[]

youbit=[]
mybit=[]


def Calculateyou(liss,number,n,myrecord,yourecored):
    if len(liss)==1:
        yourecored.append(liss[0][0:5])
        cd=yourecored.copy()
        youbit.append(cd)
        yourecored.pop()
        return liss
    for inde,lii in enumerate(liss):
        subblis=liss.copy()
        subblis.remove(lii)
        yourecored.append(lii[0:5])
        if n==10:
          if len(yourecored)>=3:
             yourecored=yourecored[0:1]+yourecored[len(yourecored)-1:]
        Calculatemy(subblis,number,n,myrecord,yourecored)


def Calculatemy(lis,number,n,myrecord,yourecord):
    n-=1
    for index,li in enumerate(lis):
        if len(lis)==6:
            myrecord=[]
        myrecord.append(li[0:5])
        if n==9:
           cc=myrecord.copy()
           mybit.append(cc)
           myrecord.pop()
        elif n==10:
            if len(myrecord)>=3:
                 myrecord=myrecord[0:1]+myrecord[len(myrecord)-1:]
        sublis=lis.copy()
        sublis.remove(li)
        Calculateyou(sublis,number,n,myrecord,yourecord)
Calculatemy(our+my,myexittotal,dataset[1][0]-myexitesnumber,[],[])

# 上面出现的数字是针对case25的,case25跑同之后再替换掉  作用是根据递归的深度来获取路径



myn=0
mycr=[]
youn=0
youcr=[]



print("                         ")
print("                         ")
print("                         ")
print("                         ")
print("                         ")
print("                         ")
print("                         ")
print("myexits")
for mi in myexits:
    a=0
    for i in mi[13:21]:
        a+=int(i)
    mycr.append(mi[0:5])
    print(mi,a)

print("            ")
print("youexits")
for yi in youexits:
    a=0
    for i in yi[13:21]:
        a+=int(i)
    youcr.append(yi[0:5])
    print(yi,a)

print("            ")


print("my")
for m in my:
    a=0
    for i in m[13:21]:
        a+=int(i)
    print(m,a)

print("            ")
print("you")
for y in you:
    a=0
    for i in y[13:21]:
        a+=int(i)
    print(y,a)
print("            ")
print("our:")
for o in our:
    a=0
    for i in o[13:21]:
        a+=int(i)
    print(o,a)


print(record)
for index,yyi in enumerate(mybit):
    temp=myexittotal
    for zi in yyi:
        temp+=int(record[zi])

    if temp>myn:
        myn=temp
        mycr=mykey+yyi
        youcr=youkey+youbit[index]

for ni in youcr:
    youn+=int(record[ni])
print("mycore:",myn,"youcore:",youn)
print("myrecord:",mycr,"yourecord:",youcr)

