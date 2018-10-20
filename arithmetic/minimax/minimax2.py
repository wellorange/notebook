


testdata=[[[[3,17],[2,12]],[[15],[25,0]]],[[[2,5],[3]],[]]]



def minimax(data,ty,depth,alpha,beta):
    dat=data
    t=ty
    dep=depth
    al=alpha
    be=beta
    if dep>=4:
      return dat
    if ty=="max":
        maxScore=-1000
        for ii in dat:
            d=ii
            childScore=minimax(d,changetype(t),dep+1,al,be)
            if childScore> maxScore:
                maxScore=childScore
                al=maxScore
            if al>=be:
                break
        return maxScore
    elif ty=="min":
        minScore=1000
        for ii in dat:
            d=ii
            childScore=minimax(d,changetype(t),dep+1,al,be)
            if childScore < minScore:
                minScore=childScore
                be=minScore
            if al>=be:
                break

        return minScore



def changetype(ty):
   if ty=="max":
       return "min"
   else:
       return "max"

print(minimax(testdata,"max",0,0,0))
