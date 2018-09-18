package raptorQ

import (
	"fmt"
	"math"
	"sync"
)

type TuplS struct {
	a int
	b  int
	d  int
	a1 int
	b1 int
	d1 int
}

type Element struct {
	val int
}

type Degree struct {
	ori   int
	curr  int
	gtone int
}

type Common struct {
	K   int //symbol number of a source block
	I   int // index of K1 in lookup table
	K1  int //symbol number of a extended source block
	J_K int // systematic index of K1
	T   int //symbol size
	X   int //not use
	S   int //LDPC row number
	H   int //HDPC row number
	W   int //LT symbol number 
	P   int //Permanent inactivated symbol number 占时理解为中间符号的数量
	P1  int //smallest prime>=P
	U   int // P - H
	B   int // W - S
	H1  int //ceil(H/2)
   
   // 矩阵满足 A*C=C1
	L        int     //K1+S+H
	N        int     //received symbol number; set to K1 for encoder
	N1       int     //received symbol plus extended one
	M        int     //N+S+H
	Tuples   []TuplS //Tuples list
	tupl_len int
	C1       []Symbol // size M: S+H(zero)+ K(symbols)+(K1-k)(zero)
	C        []Symbol // L intermediate symbols  中间符号
	sources  [][]byte // pointer to original source array 数据源
	sour  [][]byte // pointer to original source array 数据源
	R        []Symbol // repair symbols           修正符号
	A        [][]byte // generator matrix         生成矩阵
	Abak     [][]byte // backup of A              A的备份
	degree   []Degree //number of 1 in row i    
	dgh      int      // degree list head
	isi      []int    //Encoding Symbol ID list
	status   int      /* 1: para inited, 2: source filled 3: intermediate generated 4: repair generated */
}

func (c *Common) Prepare(K, N, T int) bool {
	if c.setup(K, N, T) == false {
		return false
	}
	c.Tuple()
	c.Matrix_GLDPC()
	
	c.Matrix_GHDPC()
	c.Matrix_GLT()
	for i := 0; i < c.M; i++ {
		for j := 0; j < c.L; j++ {
			c.Abak[i][j] = c.A[i][j]
		}
	}
	c.status = 1
 
	return true
}

func (c *Common) Tuple() {
	c.Tuples = make([]TuplS, c.M+40)
	for i := 0; i < c.M+40; i++ {
		c.Tuples[i] = c.Tupl(i)
	}
    
	c.tupl_len = c.M + 40
}

// 矩阵的LDPC行生成,也就是B部分,也就是第一部分
func (c *Common) Matrix_GLDPC() {

	// G_LDPC_1 部分生成, 不为什么文档这么规定的算法
	for i := 0; i < c.B; i++ {

		a := 1 + i/c.S
		b := i % c.S
		c.A[b][i] = 1

		b = (b + a) % c.S
		c.A[b][i] = 1

		b = (b + a) % c.S
		c.A[b][i] = 1
	}

	// i_s部分,不为什么,文档这么规定的
	for i := 0; i < c.S; i++ {
		c.A[i][c.B+i] = 1
	}

	// G_LDPC,2 部分
	for i := 0; i < c.S; i++ {
		a := i % c.P
		b := (i + 1) % c.P
		c.A[i][c.W+a] = 1
		c.A[i][c.W+b] = 1
	}


	return
}

// 矩阵第二,三部分生成,
func (c *Common) Matrix_GHDPC() {

	for j := 0; j < c.K1+c.S-1; j++ {
		i := RandYim(j+1, 6, c.H)
		c.A[c.S+i][j] = 1
		i = (RandYim(j+1, 6, c.H) + RandYim(j+1, 7, c.H-1) + 1) % c.H
		c.A[c.S+i][j] = 1
	}
    
	for i := c.S; i < c.S+c.H; i++ {
		c.A[i][c.K1+c.S-1] = OCT_EXP[i-c.S]
	}

	for i := c.S; i < c.S+c.H; i++ {
		for j := 0; j < c.K1+c.S; j++ {
			tmp := byte(0)
			for k := j; k < c.K1+c.S; k++ {
				if c.A[i][k] != 0 {
					tmp ^= octmul(c.A[i][k], OCT_EXP[k-j])
				}

			}
			c.A[i][j] = tmp
		}
	}
	/* identity part */
	for i := c.S; i < (c.S + c.H); i++ {
		c.A[i][c.K1+i] = 1
	}

	return
}

func (c *Common) Matrix_GLT() {
	
	var tupl TuplS
	i, j := 0, 0
	for i = 0; i < c.N1; i++ {
		if c.isi[i] < c.tupl_len {
			tupl = c.Tuples[c.isi[i]]
		} else {
			tupl = c.Tupl(c.isi[i])
		}
		a := tupl.a
		b := tupl.b
		d := tupl.d
		c.A[c.S+c.H+i][b] = 1

		for j = 1; j < d; j++ {
			b = (b + a) % c.W
			c.A[c.S+c.H+i][b] = 1
		}

		a = tupl.a1
		b = tupl.b1
		d = tupl.d1

		for b >= c.P {
			b = (b + a) % c.P1
		}
		c.A[c.S+c.H+i][c.W+b] = 1

		for j = 1; j < d; j++ {
			b = (b + a) % c.P1
			for b >= c.P {
				b = (b + a) % c.P1
			}
			c.A[c.S+c.H+i][c.W+b] = 1
		}
		
	}

}

func (c *Common) setup(K, N, T int) bool {
	if K < 1 || K > 56403 {
		return false
	}

	if N < K {
		return false
	}
	c.K, c.T, c.N = K, T, N

	I := 0
	for K >= lookup_table[I][0] {
		I++
	}
	c.I = I
	c.K1 = lookup_table[I][0]
	c.J_K = lookup_table[I][1]
	c.S = lookup_table[I][2]
	c.H = lookup_table[I][3]
	c.W = lookup_table[I][4]

	c.N1 = c.K1 - c.K + c.N

	c.L = c.K1 + c.S + c.H
	c.M = c.N1 + c.S + c.H

	c.P = c.L - c.W
	c.U = c.P - c.H
	c.B = c.W - c.S

	c.H1 = int(math.Ceil(float64(c.H) / float64(2.0)))

	// _P1_len是大于P的最小质数,
	_P1_len, flag := c.P, false
	for !flag {
		for i := 2; i <= int(Sqrt(float32(_P1_len))); i++ {
			if _P1_len%i == 0 {
				_P1_len++
				flag = false
				break
			} else {
				flag = true
			}
		}
	}

	c.P1 = _P1_len

	c.C1 = make([]Symbol, c.M)
	for i := 0; i < c.M; i++ {
	    sy:=Symbol{nil, c.T, 0, 0}
		sy.init(c.T)
		c.C1[i] = Symbol{nil, c.T, 0, 0}
	}

	c.C = make([]Symbol, c.L)
	for i := 0; i < c.L; i++ {
		c.C[i] = Symbol{nil, c.T, 0, 0}
		c.C[i].esi = i
	}

	c.A = make([][]byte, c.M)
	for i := 0; i < c.M; i++ {
		c.A[i] = make([]byte, c.L)
	}

	c.Abak = make([][]byte, c.L)
	for i := 0; i < c.M; i++ {
		c.Abak[i] = make([]byte, c.L)
	}

	c.isi = make([]int, c.N1)
	for i := 0; i < c.N1; i++ {
		c.isi[i] = i
	}

	c.degree = make([]Degree, c.M)
	return true
}

// 元祖发生器,参数X就是isi
func (c *Common) Tupl(X int) (tupl TuplS) {
	A := (53591 + c.J_K*997)
	if A%2 == 0 {
		A = A + 1
	}
	B := 10267 * (c.J_K + 1)
	y := (B + X*A)   // (B+X*A)%(1<<32)  y不要不大于4GB结果其实是一样的

	v := RandYim(y, 0, 1<<20)
	//fmt.Println(v)
	tupl.d = Deg(v,c.W)
	tupl.a = 1 + RandYim(y, 1, c.W-1)
	tupl.b = RandYim(y, 2, c.W)

	if tupl.d < 4 {
		tupl.d1 = 2 + RandYim(X, 3, 2)
	} else {
		tupl.d1 = 2
	}
	tupl.a1 = 1 + RandYim(X, 4, c.P1-1)
	tupl.b1 = RandYim(X, 5, c.P1)
	return
}


// 预编码
func (c *Common) Code(source [][]byte, N int, esi []int) bool {
    var i,j int
	
	if N < c.K {
		return false
	}

	// 不是第一次运行,从备份恢复A
	if c.status != 1 {
		for i = 0; i < c.M; i++ {
			for j = 0; j < c.L; j++ {
				c.A[i][j] = c.Abak[i][j]
			}
		}
	}

	c.status = 2

	// esi只有在解码时候才会有
	if esi != nil {
		N1 := N + c.K1 - c.K
		c.A = make([][]byte, N1+c.S+c.H)
		
		for i := 0; i < (N1 + c.S + c.H); i++ {
			c.A[i] = make([]byte, c.L)
		}

		for i := 0; i < (N1 + c.S + c.H); i++ {
			for j := 0; j < c.L; j++ {
				if i < c.S+c.H {
					c.A[i][j] = c.Abak[i][j]
				} else {
					c.A[i][j] = 0
				}

			}

			c.Abak = make([][]byte, N1+c.S+c.H)
			for i := 0; i < (N1 + c.S + c.H); i++ {
				c.Abak[i] = make([]byte, c.L)
			}

			for i := 0; i < (N1 + c.S + c.H); i++ {
				for j := 0; j < c.L; j++ {
					if i < c.S+c.H {
						c.Abak[i][j] = c.A[i][j]
					} else {
						c.Abak[i][j] = 0
					}

				}

				c.degree = make([]Degree, N1+c.S+c.H)

				c.C1 = make([]Symbol, N1+c.S+c.H)
				for i := 0; i < N1+c.S+c.H; i++ {
					c.C1[i] = Symbol{nil, c.T, 0, 0}
				}

				c.M = N1 + c.S + c.H
				c.N = N
				c.N1 = N1

				c.isi = make([]int, c.N1)

				for i := 0; i < c.N; i++ {
					if esi[i] < c.K {
						c.isi[i] = esi[i]
					} else {
						c.isi[i] = esi[i] + c.K1 - c.K
					}
				}

				for i := N; i < c.N1; i++ {
					c.isi[i] = i - N + c.K
				}

				c.Matrix_GLT()

				for i := c.S + c.H; i < c.M; i++ {
					for j := 0; j < c.L; j++ {
						c.Abak[i][j] = c.A[i][j]
					}

				}

			}
		}

	}
	for i := 0; i < c.L; i++ {
		c.C[i].init(c.T)
		c.C[i].esi = i
	}
	
    for i:=0;i<len(c.C1);i++{
		c.C1[i].init(c.T)
	}

	for i := c.S + c.H; i < c.M; i++ {
		if i < c.S+c.H+c.N {
			c.C1[i].init(c.T)
			c.C1[i].fillData(source[i-c.S-c.H], c.T)
		} else { //padding
			c.C1[i].init(c.T)
		}

	}

   c.sources=source
	return true
}



// 中间符号生成
func (c *Common) Generate_intermediates(soe [][]byte) []Symbol {
     var i,k1,k2,j int
	// 计算A^^-1生成中间符号C,使用高斯消元算法
	if c.status != 2 {
		return nil
	}

	
	 var  ss sync.Pool
	 ss.Put(soe)
	cols1 := make([]int, c.L)
	cols2 := make([]int, c.L)
	
	for i := 0; i < c.M; i++ {
		d := int(0)
		gtone := 0
		for j := 0; j < c.L-c.P; j++ {
			if c.A[i][j] != 0 {
				d++
				if c.A[i][j] > 1 {
					gtone++
				}
			}
			
		}
		c.degree[i].ori = d
		c.degree[i].curr = d
        c.degree[i].gtone = gtone
	}


   
	// setup 1
	_I, _U, r, gtone_start := 0, c.P, 0, 0
	for _I+_U < c.L {
       

		var index,  o int
	retry:
		index, o, r = c.M, c.L, c.L
		
		for i = _I; i < c.M; i++ {
			if (gtone_start != 0 || (gtone_start == 0 && c.degree[i].gtone == 0)) && c.degree[i].curr > 0 && c.degree[i].curr <= r {
				
				index = i
				
				if c.degree[i].curr < r || (c.degree[i].curr == r && c.degree[i].ori < o) {
					o = c.degree[i].ori
					r = c.degree[i].curr	
				}
			}
		}
		if index == c.M {
			if gtone_start != 0 {
				goto retry
			} else {
				return nil
			}
		}



		c.swap_row(_I, index)
		
		
	
     
	   k1=0
	   k2=0
		for j = _I; j < c.L-_U; j++ {
			if j < c.L-_U-r+1 {
				if c.A[_I][j] != 0 {
					
					cols1[k1] = j
					k1++
				
				}
			} else {
				if c.A[_I][j] == 0 {
				
					cols2[k2] = j
					k2++
				}
			}
		}





		if k1 != k2+1 {
			return nil

		}
   

		c.swap_col(_I, int(cols1[0]))
    
		for j = 0; j < k2; j++ {
			
			c.swap_col(cols2[j], cols1[j+1])

		}
		

		if c.A[_I][_I] > 1 {
			v := c.A[_I][_I]
			c.C1[_I].div(v)
			for j = _I; j < c.L; j++ {
				c.A[_I][j] = octdiv(c.A[_I][j], v)
			}
		}
      
      
    
		for i = _I + 1; i < c.M; i++ {
			v := c.A[i][_I]
			if v != 0 {
				c.A[i][_I] = 0
				c.degree[i].curr--
				if v > 1 {
					c.degree[i].gtone--
				}
				for j = c.L - _U - (r - 1); j < c.L; j++ {
					oldv := c.A[i][j]
					c.A[i][j] ^= octmul(v, c.A[_I][j])
					if j < c.L-_U {
						if c.A[i][j] > 0 {
							c.degree[i].curr++
							if c.A[i][j] > 1 {
								c.degree[i].gtone++
							}
						}
						if oldv > 0 {
							c.degree[i].curr--
							if oldv > 1 {
								c.degree[i].gtone--
							}
						}
					}
				}
				c.C1[i].muladd(c.C1[_I], v)
				
			}

			for j = c.L - _U - (r - 1); j < c.L-_U; j++ {
				if c.A[i][j] != 0 {
					c.degree[i].curr--
					if c.A[i][j] > 1 {
						c.degree[i].gtone--
					}
				}
			}
		}


		_I++
		_U += r - 1
		
		
	}
	
	c.Test()
//	c.Test()

	if !c.gaussian_elimination(_I, _I) {
		return nil
	}
	// fmt.Println("_I=",_I,"_U=",_U)

	/* step 3 */
	for jj := _I; jj < c.L; jj++ {
		for i := 0; i < _I; i++ {
			v := c.A[i][jj]
			if v != 0 {
				c.A[i][jj] = 0
				c.C1[i].muladd(c.C1[jj], v)
			}
		}
	}
    // c.Test()

	for i := 0; i < c.L; i++ {
		c.C1[i].esi = c.C[i].esi
	}

	for i := 0; i > c.L; i++ {
		s := c.C[c.C1[i].esi]
		c.C[c.C1[i].esi] = c.C1[i]
		c.C1[i] = s
	}
	
	c.status = 3


	return c.C
}

//  高斯消元的实现
func (c *Common) gaussian_elimination(starti, startj int) bool {
	var i, k, q, jj, kk int
	var firstone int

	HI := make([]int, c.L)
	LOW := make([]int, c.M)
	for jj = startj; jj < c.L; jj++ {
		k = 0
		for i := starti; i <= jj-1; i++ {
			if c.A[i][jj] != 0 {
				HI[k] = i
				k++
			}
		}

		kk = 0
		firstone = -1
		for i = jj; i < c.M; i++ {
			if c.A[i][jj] != 0 {
				LOW[kk] = i
				if c.A[i][jj] == 1 && firstone == -1 {
					firstone = kk
				}
				kk++
			}
		}
		if kk == 0 {
			return false
		}

		if firstone > 0 {
			t := LOW[0]
			LOW[0] = LOW[firstone]
			LOW[firstone] = t
		}

		if c.A[LOW[0]][jj] != 1 {
			v := c.A[LOW[0]][jj]

			c.C1[LOW[0]].div(v)
			for q := jj; q < c.L; q++ {
				c.A[LOW[0]][q] = octdiv(c.A[LOW[0]][q], v)
			}
		}

		for i = 1; i < kk; i++ {
			v := c.A[LOW[i]][jj]

			c.C1[LOW[i]].muladd(c.C1[LOW[0]], v)

			for q = jj; q < c.L; q++ {
				c.A[LOW[i]][q] ^= octmul(c.A[LOW[0]][q], v)
			}
		}

		for i = 0; i < k; i++ {
			v := c.A[HI[i]][jj]
			c.C1[HI[i]].muladd(c.C1[LOW[0]], v)

			for q = jj; q < c.L; q++ {
				c.A[HI[i]][q] ^= octmul(c.A[LOW[0]][q], v)
			}
		}

		if LOW[0] != jj {

			temp := c.A[jj]
			c.A[jj] = c.A[LOW[0]]
			c.A[LOW[0]] = temp

			tempo := c.C1[jj]
			c.C1[jj] = c.C1[LOW[0]]
			c.C1[LOW[0]] = tempo
		}
		
	}

	return true
}

// 使用中间符号与编码器生成修正符号
func (c *Common) Generate_repairs(count int) []Symbol {

	if c.status != 3 {
		return nil
	}

	c.R = make([]Symbol, count)

	for i := c.K; i < c.K+count; i++ {
		var tupl TuplS
		isi := i + c.K1 - c.K
		if isi < c.tupl_len {
			tupl = c.Tuples[isi]
		} else {
			tupl = c.Tupl(isi)
		}
		c.R[i-c.K] = c.LTEnc(c.C, tupl)
	}

	c.status = 4
	return c.R
}




func (c *Common) Recover_symbol(x int) Symbol{
	if x >=c.K {
		return Symbol{}
	}
	tu:=c.Tuples[x]
	s:=c.LTEnc(CC,tu)
	return s
}





func (c *Common) LTEnc(C_L []Symbol, tupl TuplS) Symbol {
	a := tupl.a
	b := tupl.b
	d := tupl.d
	s := C_L[b]
	for j := 1; j < d; j++ {
		b = (b + a) % c.W
		s.xxor(C_L[b])
	}
	
	a = tupl.a1
	b = tupl.b1
	d = tupl.d1
	for b >= c.P {
		b = (b + a) % c.P1
	}
	s.xxor(C_L[c.W+b])
	
	for j := 1; j < d; j++ {
		b = (b + a) % c.P1;
		for b >= c.P {
			b = (b + a) % c.P1
		}
		s.xxor(C_L[c.W+b])
	}
	return s
}

func (c *Common) swap_row(i1, i2 int) {
	if i1 == i2 {
		return
	}

	e := c.A[i1]
	c.A[i1] = c.A[i2]
	c.A[i2] = e

	d := c.degree[i1]
	c.degree[i1] = c.degree[i2]
	c.degree[i2] = d

	s := c.C1[i1]
	c.C1[i1] = c.C1[i2]
	c.C1[i2] = s

}

func (c *Common) swap_col(j1, j2 int) {
	if j1 == j2 {
		return
	}

	for i := 0; i < c.M; i++ {
		t := c.A[i][j1]
		c.A[i][j1] = c.A[i][j2]
		c.A[i][j2] = t
	}

	s := c.C[j1]
	c.C[j1] = c.C[j2]
	c.C[j2] = s
}

func (c *Common) Test() {
	fmt.Println("K=", c.K, "  I=", c.I)
	fmt.Println("T=", c.T, "  J_K=", c.J_K)
	fmt.Println("H=", c.H, "  X=", c.X)
	fmt.Println("S=", c.S, "  H1=", c.H1)
	fmt.Println("L=", c.L, "  N1=", c.N1)
	fmt.Println("N=", c.N, "  P1=", c.P1)
	fmt.Println("M=", c.M, "  U=", c.U)
	fmt.Println("K1=", c.K1, "  B=", c.K1)
	fmt.Println("W=", c.W, "   P=", c.P)
	fmt.Println("tupl_len=", c.tupl_len)
	fmt.Println("Tuples:")
	for i := 0; i < c.L; i++ {
		fmt.Println("Tuple", i, "d,a,b=", c.Tuples[i].d, ",", c.Tuples[i].a, ",", c.Tuples[i].b, ",", "d1,a1,b1=", c.Tuples[i].d1, ",", c.Tuples[i].a1, ",", c.Tuples[i].b1)

	}

	fmt.Println("A maxtrix:")
	for i := 0; i < len(c.A); i++ {
		for j := 0; j < len(c.A[i]); j++ {
			fmt.Printf("%4d", c.A[i][j])
		}

		fmt.Println()
	}
	fmt.Println("C1 maxtrix:")
	for i := 0; i < len(c.C1); i++ {
			fmt.Printf("%4d ", c.C1[i],)
		

		fmt.Println()
	}
	fmt.Println("C maxtrix:")
	for i := 0; i < len(c.C); i++ {
		fmt.Printf("%4d", c.C[i])
	fmt.Println()
}  
	fmt.Println("sources maxtrix:")

	for i := 0; i < len(c.sources); i++ {
		for j := 0; j < len(c.sources[i]); j++ {
			fmt.Printf("%4d", c.sources[i][j])
		}
		fmt.Println()
	}


	
	fmt.Println("R maxtrix:")
	for i := 0; i < len(c.R); i++ {
		fmt.Printf("%4d", c.R[i])
	fmt.Println()
}

	fmt.Println("Abak maxtrix:")
	for i := 0; i < len(c.Abak); i++ {
		for j := 0; j < len(c.Abak[i]); j++ {
			//fmt.Printf("%4d", c.Abak[i][j])
		}
		//fmt.Println()
	}

}


func Sqrt(x float32) float32 {
	var xhalf float32 = 0.5 * x // get bits for floating VALUE
	i := math.Float32bits(x)    // gives initial guess y0
	i = 0x5f375a86 - (i >> 1)   // convert bits BACK to float
	x = math.Float32frombits(i) // Newton step， repeating increases accuracy
	x = x * (1.5 - xhalf*x*x)
	x = x * (1.5 - xhalf*x*x)
	x = x * (1.5 - xhalf*x*x)
	return 1 / x
}

/* section 5.7.2 */
func octmul(u, v byte) byte {
	if u == 0 || v == 0 {
		return 0
	}
	if v == 1 {
		return byte(u)
	}
	if u == 1 {
		return byte(v)
	}
	return OCT_EXP[OCT_LOG[u]+OCT_LOG[v]]
}
func octdiv(u, v byte) byte {
	if u == 0 {
		return 0
	}
	if v == 1 {
		return u
	}

	return OCT_EXP[OCT_LOG[u]-OCT_LOG[v]+255]
}

// 伪随机数产生器
func RandYim(y int, i int, m int) int {
	pow28 := math.Pow(2, 8)
	pow216 := math.Pow(2, 16)
	pow224 := math.Pow(2, 24)
	x0 := int(math.Mod(float64(y+i), pow28))
	x1 := int(math.Mod(math.Floor(float64(y)/pow28)+float64(i), pow28))
	x2 := int(math.Mod(math.Floor(float64(y)/pow216)+float64(i), pow28))
	x3 := int(math.Mod(math.Floor(float64(y)/pow224)+float64(i), pow28))
	return int(math.Mod(float64(V0[x0]^V1[x1]^V2[x2]^V3[x3]), float64(m)))
}

// 度产生器
func Deg(v int,w int) int {
	i:=0
	for i=0; i < 30; i++ {
		if F[i] < v {
		}else{
			break
		}
	}

	if i<(w-2){
		return i
	}else {
		return w-2
	}
}


var CC=[]Symbol{
{data:[]byte{119, 212, 130, 126},nbytes:4},
{data:[]byte{244, 26 ,175, 196},nbytes:4,sbn:32630,esi:1},
{data:[]byte{143, 208,  78,  74 }  ,nbytes:4 ,sbn:0 ,esi:2},
 {data:[]byte{65, 111, 161,  98 }  ,nbytes:4 ,sbn:0 ,esi:3},
 {data:[]byte{68 ,132,  93,  35 }  ,nbytes:4 ,sbn:0 ,esi:4},
 {data:[]byte{23,   9,  59,  56 }  ,nbytes:4 ,sbn:0 ,esi:5},
 {data:[]byte{79, 170, 233,  46 }  ,nbytes:4 ,sbn:0 ,esi:6},
{data:[]byte{169, 103, 167, 187 }  ,nbytes:4 ,sbn:0 ,esi:7},
{data:[]byte{205, 158,  90,  71 }  ,nbytes:4 ,sbn:0 ,esi:8},
{data:[]byte{174,  86, 158, 122 }  ,nbytes:4 ,sbn:0 ,esi:9},
 {data:[]byte{74,  76, 113, 103 }  ,nbytes:4 ,sbn:0 ,esi:10},
 {data:[]byte{12, 124, 110, 205 }  ,nbytes:4 ,sbn:0 ,esi:11},
{data:[]byte{101,  73, 124,  251 }  ,nbytes:4 ,sbn:0 ,esi:12},
{data:[]byte{229, 167,  10,   2 }  ,nbytes:4 ,sbn:0 ,esi:13},
{data:[]byte{232,  39,  29,  74 }  ,nbytes:4 ,sbn:0 ,esi:14},
{data:[]byte{225, 190,  58, 223 }  ,nbytes:4 ,sbn:0 ,esi:15},
{data:[]byte{212,  49, 202,  50 }  ,nbytes:4 ,sbn:0 ,esi:16},
{data:[]byte{23,  61, 185, 114 }  ,nbytes:4 ,sbn:0 ,esi:17},
{data:[]byte{219,  97,  63, 198 }  ,nbytes:4 ,sbn:0 ,esi:18},
{data:[]byte{214, 231, 207, 216 }  ,nbytes:4 ,sbn:0 ,esi:19},
{data:[]byte{184 ,129, 233,  18 }  ,nbytes:4 ,sbn:0 ,esi:20},
{data:[]byte{170,  29,  249, 187 }  ,nbytes:4 ,sbn:0 ,esi:21},
{data:[]byte{207,  48, 111,  59 }  ,nbytes:4 ,sbn:0 ,esi:22},
{data:[]byte{241, 242, 200, 218 }  ,nbytes:4 ,sbn:0 ,esi:23},
{data:[]byte{151, 178,  19, 167 }  ,nbytes:4 ,sbn:0 ,esi:24},
{data:[]byte{145,  95,  89,  11 }  ,nbytes:4 ,sbn:0 ,esi:25},
{data:[]byte{ 60 , 85 , 62 ,145 }  ,nbytes:4 ,sbn:0 ,esi:26},
}
