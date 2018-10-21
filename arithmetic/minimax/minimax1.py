


file_name = 'case/input25.txt'
# dataset 是从input读取的数据集合
token=("bed","park","la_num","la","my_num","my","list_num","list")
cur_token=token[0]
dataset = []

number=0
numbert=1
temp=[]


data=[]     
exits=[]     # 双方剔除过的
youkey=[]
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
mexi=[]
yoexi=[]
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

    if  line[0:5] in mykey:
        mexi.append(line)
    if  line[0:5] in youkey:
        yoexi.append(line)





# print(my)

#print("=====")

 #   print(my[13:21])
        #     1             2       3         4         5           6           7           8
       # print(line[0:5],line[5:6],line[6:9],line[9:10],line[10:11],line[11:12],line[12:13],line[13:21])
    


print("床位:",dataset[0][0],"  车位:",dataset[1][0])
print("myexit")



myn=0
myrecord=0
youn=0
yourecord=0
for mmy in mexi:
    a=0
    for i in mmy[13:21]:
        a+=int(i)
    print(mmy,a)
print("            ")
print("youexit")
for youu in yoexi:
    a=0
    for i in youu[13:21]:
        a+=int(i)
    print(youu,a)

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

bibi=[]
print("our:")
for o in our:
    a=0
    for i in o[13:21]:
        a+=int(i)
    cc=[]
    cc.append(o)
    cc.append(a)
    bibi.append(cc)
    print(o,a)




def takeSecond(elem):
    return elem[1]
# 指定第二个元素排序
bibi.sort(key=takeSecond,reverse=True)
iff="my"
for ii in bibi:
     if iff=="my":
         mykey.append(ii[0][0:5])
         iff="you"
     else :
         youkey.append(ii[0][0:5])
         iff="my"

print(record)
print(mykey)
print(youkey)
