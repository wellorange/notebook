import numpy as np
import string
import time
import os
from collections import defaultdict

def tree(): 
   return defaultdict(tree)    


file_name = 'case/input1.txt'
# dataset 是从input读取的数据集合
token=("bed","park","la_num","la","my_num","my","list_num","list")
cur_token=token[0]
dataset = []

number=0
numbert=1
temp=[]


data=[]     
exits=[]     # 双方剔除过的
mykey=[]
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
            number=0
            numbert=1
            temp=[]
         else:
              temp.append(line.split()[0])
              exits.append(line.split()[0])
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
            temp=[]
            cur_token="null"
            number=0
            numbert=1
         else:
              temp.append(line.split()[0])
              if line.split()[0][0:5] in mykey:
                  myexits.append(line.split()[0])
              numbert+=1




my=[]
you=[]
our=[]

for line in dataset[7]:
    if line[0:5] not in exits:
       if (line[10:11]=="N" and line[11:12]=="Y" and line[12:13]=="Y") and (line[5:6]!="F" or int(line[6:9])<=17 or line[9:10]!="N"):
          my.append(line)
          

       if (line[5:6]=="F" and int(line[6:9])>17 and line[9:10]=="N") and (line[10:11]!="N" or line[11:12]!="Y" or line[12:13]!="Y") :
          you.append(line)

       if line[5:6]=="F" and int(line[6:9])>17 and line[9:10]=="N" and line[10:11]=="N" and line[11:12]=="Y" and line[12:13]=="Y":
          our.append(line)




parent=0
mymumber=[0,0,0,0,0,0,0]

for ex in myexits:
    tem=ex[13:21]
    mymumber[0]+=int(tem[0])
    mymumber[1]+=int(tem[1])
    mymumber[2]+=int(tem[2])
    mymumber[3]+=int(tem[3])
    mymumber[4]+=int(tem[4])
    mymumber[5]+=int(tem[5])
    mymumber[6]+=int(tem[6])

# print(my)

#print("=====")

 #   print(my[13:21])
        #     1             2       3         4         5           6           7           8
       # print(line[0:5],line[5:6],line[6:9],line[9:10],line[10:11],line[11:12],line[12:13],line[13:21])
    


pp=0

def Calculate(number,list,n):
    global pp

    tem=list[0][13:21]
    number[0]+=int(tem[0])
    number[1]+=int(tem[1])
    number[2]+=int(tem[2])
    number[3]+=int(tem[3])
    number[4]+=int(tem[4])
    number[5]+=int(tem[5])
    number[6]+=int(tem[6])
    n+=1
    if (dataset[0][0]+1) in mymumber:
        print("asdasd")
        print(mymumber)
    else:
       if  n==dataset[0][0]:
           p=number[0]+number[1]+number[2]+number[3]+number[4]+number[5]+number[6]
           if p>pp:
               pp=p
       else:
           Calculate(number,list[1:],n)

print("my")
for m in my:
    print(m)

print("            ")
print("you")
for y in you:
    print(y)
print("            ")
print("our:")
for o in our:
   print(o)
#Calculate(mymumber,my,dataset[4][0])

#print(pp)
