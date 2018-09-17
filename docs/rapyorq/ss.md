
#include "Helper.h"
#include "Symbol.h"
#include "Tables.h"
#include "stdlib.h"

using namespace std;

//#define DEBUG 1

/* initialize parameters for raptor code 
  arg:
   _K: symbol number in a source block
   _N: received symbol number; for encoder _N = _K, for decoder, let _N = _K during init, later prepare() will fixup
   _T: symbol size in bytes
  return:
    true if initialized successfully, false otherwise
 */
bool Generators::gen(int _K, int _N, int _T) {
	int i,j;

	if (!_0_init(_K, _N, _T))
		return false;
    
	_1_Tuples();
 
	_2_Matrix_GLDPC();

	_3_Matrix_GHDPC();

	_4_Matrix_GLT();

	/* save for reuse, A will be changed in calculating intermediates */
	for( i = 0; i < (M); i++)
	{
		for(j = 0; j < (L); j++)
			Abak[i][j] = A[i][j];
	}
	status = 1;

	//    ToString();

	//while(true);

	return true;
}

/* calculate parameters, allocate data structures */
bool Generators::_0_init(int _K, int _N, int _T) {
	int i,j;

	if (_K < 1 || _K > 56403) {
		cout << "Invalid K, should in [1,56403]\n";
		return false;
	}

	if (_N < _K) {
		cout << "N should not be smaller than K\n";
		return false;
	}

	K = _K;
	T = _T;
	N = _N;

	I=0;
	while (K > lookup_table[I][0])
		I++;

	K1 = lookup_table[I][0];
	J_K = lookup_table[I][1];
	S = lookup_table[I][2];
	H = lookup_table[I][3];
	W = lookup_table[I][4];

	N1 = K1 - K + N;

	L = K1 + S + H;
	M = N1 + S + H;

	P = L - W;
	U = P - H;
	B = W - S;
	
	H1 = (int)ceil(H/2.0);

	//P1 be the smallest prime that is greater than or equal to P
	int _P1_len = P;
	int flag = false;
	while(!flag)
	{		
		for(i = 2; i <= sqrt((double)_P1_len); i++)
		{
			if(_P1_len % i == 0)
			{
				_P1_len++;
				flag = false;
				break;
			}
			else
				flag = true;
		}
	}

	P1 = _P1_len;

	try {

		C1 = new Symbol*[M];
		for (i = 0; i < M; i++)
			C1[i] = new Symbol(T);

		C = new Symbol*[L];
		for (i = 0; i < L; i++) {
			C[i] = new Symbol(T);
			C[i]->esi = i;
			cout<<i<<endl;
		}
		R = NULL;
		A = new Elem*[M];
		for(i = 0; i < (M); i++)
			A[i] = new Elem[L];

		Abak = new Elem*[M];
		for(i = 0; i < (M); i++)
			Abak[i] = new Elem[L];


		for( i = 0; i < (M); i++)
		{
			for(j = 0; j < (L); j++)
				A[i][j].val = 0;
				//memset(&A[i][j], 0, sizeof(Elem));
		}

	

		/* isi[i]=i for encoder;
		   for decoder, isi[i] == the ESI for the ith received symbol for i in[0,N), ==i-N+K for [N,N1) */
		
		isi = new int[N1];
		for (i=0;i<N1;i++)
			isi[i] = i;

		degree = new Degree[M];
	
	} catch ( bad_alloc &e) {
		cout << "Allocation failed!" << e.what() << endl;
		abort();
	}

	return true;
}


/* calculate Tuples. 
   to avoid repeated calculation, we generate most possible ones during init
*/
#define MAX_OVERHEAD (40)
void Generators::_1_Tuples(void) {
	int i;

	try {	
		this->Tuples = new TuplS[M + MAX_OVERHEAD];
	} catch ( bad_alloc &e) {
		cout << "Allocate Tuples[] failed!" << e.what() << endl;
		abort();
	}

	for (i=0; i < M + MAX_OVERHEAD; i++) {
		Tuples[i] = Tupl(i);
	}

	tupl_len = M + MAX_OVERHEAD;

}

/* Generate the LDPC rows of matrix A */ 
void Generators::_2_Matrix_GLDPC() {
	int i;
	int a,b;

	/* G_LDPC,1 */
	for(i = 0; i < B; i++)
	{
		a = 1 + i / S;
		b = i % S;
		A[b][i].val = 1;

		b = (b + a) % S;
		A[b][i].val = 1;

    	b = (b + a) % S;
		A[b][i].val = 1;
	}

	/* identity part */
	for(i = 0; i < S; i++)
		A[i][B + i].val = 1;

	/* G_LDPC,2 */
	for (i = 0; i < S; i++) {
		a = i % P;
		b = (i + 1) % P;
		A[i][W + a].val = 1;
		A[i][W + b].val = 1;
	}
  
	return ;
}

/* Generate G_HDPC */
void Generators::_3_Matrix_GHDPC()
{
	int i,j,k;

	for (j=0; j< K1 + S - 1; j++) {
		i = RandYim(j+1, 6, H);
		A[S+i][j].val = 1;
		i = (RandYim(j+1, 6, H) + RandYim(j+1, 7, H - 1) + 1) % H;
		A[S+i][j].val = 1;
	}

	for (i=S; i < S + H; i++) {
		A[i][K1 + S - 1].val = OCT_EXP[i-S];
	}

	for (i=S; i< S + H; i++) {
		for (j=0; j < K1 + S; j++) {
			unsigned char tmp = 0;
			for (k=j; k < K1 + S; k++) 
				if (A[i][k].val) tmp ^= octmul(A[i][k].val,OCT_EXP[k - j]); 
			A[i][j].val = tmp;
		}
	}

	/* identity part */
	for(i = S; i < (S + H ) ; i++)
		A[i][K1  + i].val = 1;	
	
	return;
}

/* LT rows of maxtrix A */
void Generators::_4_Matrix_GLT()
{
	int i,j;
	int a,b,d;

	TuplS tupl;

	int flag1 = 0;	

	i = 0;
	while(i < N1)
	{ 
		if ( isi[i] < tupl_len) 
			tupl = Tuples[isi[i]];
		else
			tupl = Tupl(isi[i]);

		a = tupl.a;
		b = tupl.b;
		d = tupl.d;

		A[S + H + i][b].val = 1;

		for(j = 1; j < d; j++)
		{
			b = (b + a) % W;
			A[S + H + i][b].val = 1;
		}

		a = tupl.a1;
		b = tupl.b1;
		d = tupl.d1;

		while (b >= P)
			b = (b + a) % P1;
		A[S + H + i][W + b].val = 1;

		for(j = 1; j < d; j++)
		{
			b = (b + a) % P1;
			while(b >= P)
				b = (b + a) % P1;
			A[S + H + i][W + b].val = 1;
		}

		i++;
	}

	return;
}


/* LT encoder
arg:
   _K is K1
   C_L is the intermediate symbols
   tupl is the Tuple for giving ESI x
return:
   pointer to encoded symbol
   caller need to handle the delete of symbol space
   根据中间符号解码
 */
Symbol* Generators::LTEnc(Symbol** C_L ,TuplS tupl) {
	int a = tupl.a;
	int b = tupl.b;
	int d = tupl.d;
	Symbol *s;
    

	//  初始化一个symbol
	s = new Symbol(T);
   

	*s = *C_L[b];
	// s分别于指定的symbol进行异或运算,结果保存到s中
	for (int j=1; j < d; j++) {
		b=(b + a) % W;
		s->xxor(C_L[b]);
	}

	a = tupl.a1;
	b = tupl.b1;
	d = tupl.d1;

	while (b >= P)
			b = (b + a) % P1;
	s->xxor(C_L[W + b]);

	for(int j = 1; j < d; j++)
	{
		b = (b + a) % P1;
		while(b >= P)
			b = (b + a) % P1;
		s->xxor(C_L[W + b]);
	}

	return s;
}






/* Tuple generator */
// 元祖发生器,用于生成元祖的数据参数
const TuplS Generators::Tupl(int X) {
	TuplS tupl;
	
	unsigned int A = (53591 + J_K*997);
	if (A % 2 == 0) A = A + 1;
	unsigned int B = 10267*(J_K + 1);
	unsigned int y = (B + X*A);
	int v = RandYim(y,0,1048576);
	tupl.d = Deg(v);
	tupl.a = 1 + RandYim(y, 1, W - 1);
	tupl.b = RandYim(y, 2, W);

	if (tupl.d < 4) tupl.d1 = 2 + RandYim(X, 3, 2); else tupl.d1 = 2; 
	tupl.a1 = 1 + RandYim(X, 4, P1 - 1);
	tupl.b1 = RandYim(X, 5, P1);
	return tupl;
}

// 随机数参数器,用来产生随机数
const unsigned int Generators::RandYim(unsigned int y, unsigned char i,unsigned int m) {
	return (V0[((y & 0xff) +i) & 0xff] ^ V1[(((y>>8) & 0xff) + i) & 0xff] ^ V2[(((y>>16) & 0xff) + i) & 0xff] ^ V3[(((y>>24) & 0xff) + i) & 0xff]) % m;
}


/* find j, f[j-1] <= v <f[j]
   d = min(j, W - 2)
   v must < 2^^20
 */
// 度发发生器用来生成度参数
const unsigned int Generators::Deg(unsigned int v) {
	int j=0;
	while( v > f[j])
		j++;
	return min(j, W - 2);
}

/* clean up leftovers of last run, fill in new source block; 
   called by the encoder/decoder, it can be called many times
        encoder->init
		while(get more source) {
		   prepare
		   decode
	    }
arg:
   source: source symbol list with Length _N
   _esi: the ESI list for the source list
return:
   success or not
 */
bool Generators::prepare(char **source, int _N, int *_esi) {

	int i,j;

	if (_N < K) {
		cout << "Invalid N in prepare!" << N << endl;
		return false;
	}

	if (status != 1) {

		/* Not the first run, recover A;
		 */
		for( i = 0; i < (M); i++)
		{
			for(j = 0; j < (L); j++)
				A[i][j] = Abak[i][j];
		}
	}

	status = 2;

	try {
		/* only decoder will provide esi(for each source block) */
		if (_esi) {
			int _N1 = _N + K1 - K;

			/* maxtrix LT parts changed, the data struct is not efficent now */
			for (i=0; i < M ;i++)
				delete []A[i];
			delete []A;

			A = new Elem*[_N1 + S + H];
			for(i = 0; i < (_N1 + S + H); i++)
				A[i] = new Elem[L];

			for(i = 0; i < (_N1 + S + H); i++)
			{
				for(j = 0; j < (L); j++)
					if (i < S + H)
						A[i][j] = Abak[i][j];
					else
						A[i][j].val = 0;
			}

			for (i=0; i < M ;i++)
				delete []Abak[i];
			delete []Abak;

			Abak = new Elem*[_N1 + S + H];
			for(i = 0; i < (_N1 + S + H); i++)
				Abak[i] = new Elem[L];

			for(i = 0; i < (_N1 + S + H); i++)
			{
				for(j = 0; j < (L); j++)
					if (i < S + H)
						Abak[i][j] = A[i][j];
					else
						Abak[i][j].val = 0;

			}


			delete[] degree;
			degree = new Degree[_N1 + S + H];
			
			for (i = 0; i < M; i++)
				delete C1[i];
			delete[] C1;

			this->C1 = new Symbol*[_N1 + S + H];
			for (i = 0; i < _N1 + S + H; i++)
				C1[i] = new Symbol(T); 

			M = _N1 + S + H;
			N = _N;
			N1 = _N1;

			delete[] isi;

			isi = new int[N1];

			for (i = 0; i < N; i++) {
				if (_esi[i] < K)
					isi[i] = _esi[i];
				else
					isi[i] = _esi[i] + K1 - K;
			}
			/* K1 - K (N1 - N) padding symbols */
			for (i=N; i < N1; i++)
				isi[i] = i - N + K;

			_4_Matrix_GLT();

			for (i = S + H; i < M; i++)
				for (j=0; j < L; j++)
					Abak[i][j] = A[i][j];

		}
	} catch ( bad_alloc &e) {
		cout << "Allocation in prepare() failed!" << e.what() << endl;
		abort();
	}


	for (i=0; i<L; i++) {
		C[i]->init(T);
		C[i]->esi = i;
	}
   
	for (i=S + H; i< M; i++)
		if (i < S + H + N)
		    // 把源码的一组拆分合并为一个symbol
			C1[i]->fillData(source[i - S - H],T);
		else 
			C1[i]->init(T);

	sources = source;

	return true;
}

void Generators::swap_row(int i1, int i2)
{
	if (i1 == i2) return;

	Elem *e;
	e = A[i1];
	A[i1] = A[i2];
	A[i2] = e;

	Degree d = degree[i1];
	degree[i1] = degree[i2];
	degree[i2] = d;

	
	Symbol *s;
	s = C1[i1];
	C1[i1] = C1[i2];
	C1[i2] = s;

}

void Generators::swap_col(int j1, int j2)
{
	if (j1 == j2) return;
	
	for (int i=0; i<M; i++) {
		int t = A[i][j1].val;
		A[i][j1].val = A[i][j2].val;
		A[i][j2].val = t;
	}

	Symbol *s = C[j1];
	C[j1] = C[j2];
	C[j2] = s;	
}

/* A[i1] = A[i1] ^ A[i2]; C1[i1] = C1[i1] ^ C1[i2]; */
void Generators::xxor(int i1, int i2, int U)
{
	if (i1 == i2) return;

	int d = 0;
	for (int j=i2; j<L; j++) {
		A[i1][j].val ^= A[i2][j].val;
		if (j < U && A[i1][j].val == 1) d++;
	}
	degree[i1].curr = d;
	C1[i1]->xxor(C1[i2]);
}

int Generators::gaussian_elimination(int starti, int startj)
{
	int i, k, q, jj, kk;
	int firstone;

	int* HI  = new int[L];
	int* LOW = new int[M];	
		
	for (jj=startj; jj<L; jj++)
	{
		//PrintMatrix();
		k=0;
		for (i=starti; i<=jj-1; i++)
		{
			if (A[i][jj].val)
			{
				HI[k]=i;
				k++;
			}
		}
			
		kk=0;	
		firstone = -1;
		for (i=jj; i<M; i++)
		{
			if(A[i][jj].val)
			{
				LOW[kk]=i;
				if (A[i][jj].val == 1 && firstone == -1) firstone = kk;
				kk++;
			}
		}
		
		if (kk==0){
			cout << " Encoder: due to unclear reasons the process can not continue" << endl;
			delete[] HI;
			delete[] LOW;
			return 0;
		}
		
		
		if (firstone > 0) {
			int t = LOW[0];
			LOW[0] = LOW[firstone];
			LOW[firstone] = t;
		}
		

		if (A[LOW[0]][jj].val != 1) {
			unsigned char v = A[LOW[0]][jj].val;

			C1[LOW[0]]->div(v);
			for (q=jj; q<L; q++)
				A[LOW[0]][q].val = octdiv(A[LOW[0]][q].val, v);
		}

		for (i=1; i<kk; i++)
		{
			unsigned char v = A[LOW[i]][jj].val;

			C1[LOW[i]]->muladd(C1[LOW[0]],v);
			
			for (q=jj; q<L; q++)
				A[LOW[i]][q].val ^= octmul(A[LOW[0]][q].val, v);
		}
		
		for (i=0; i<k; i++)
		{
			unsigned char v = A[HI[i]][jj].val;

			C1[HI[i]]->muladd(C1[LOW[0]], v);
			
			for (q=jj; q<L; q++ )
				A[HI[i]][q].val ^= octmul(A[LOW[0]][q].val, v);
		}
		
		if (LOW[0] != jj) {
			Elem* temp;
			Symbol *tempo;

			temp = A[jj];
			A[jj] = A[LOW[0]];
			A[LOW[0]] = temp;

			tempo = C1[jj];
			C1[jj] = C1[LOW[0]];
			C1[LOW[0]] = tempo;
		}
		//PrintMatrix();
	}

	return 1;
}

// check correctness of intermediates
void Generators::verify(void)
{
	int i,j;
	Symbol *s,*s1;
	char *p;

	s = new Symbol(T);
	s1 = new Symbol(T);

	for (i = 0; i < M; i++) {
		s->init(T);
		for (j = 0; j < L; j++) {
			if (Abak[i][j].val)
				s->muladd(C[j], Abak[i][j].val);
		}

		if (i < S + H || i >= S + H + N)
			p = (char*)s1->data;
		else
			p = (char*)sources[i - S - H];
		
		if (memcmp(s->data, p, T) != 0) {
			printf("Check fail for line %d,%x vs %x\n", i, *(int*)s->data, *(int*)p);
		}
	}
	
	delete s;
	delete s1;
}


/* core decode algorithm here. Calculate C from equation A * C = C1 
   return the intermediate symbol list C, NULL if decode fail
 */
Symbol ** Generators::generate_intermediates(void)
{
	// ToString();
  /* calculate A^^-1 to get intermediate symbol list C by Gaussian elimination */

	if (status != 2) {
		cout << "Wrong call sequence! Filling the source block before generate intermediates" << endl;
		return NULL;
	}


	int* cols1 = new int[L];
	int* cols2 = new int[L];

	int k1,k2;
	

	/* init degree list */
	int i,j,d,gtone;
	for (i = 0; i < M; i++) 
	{
		d = 0;
		gtone = 0;
		for (j=0; j < L - P; j++)
			if (A[i][j].val) {
				d++;
				if (A[i][j].val > 1) gtone++;
			}
		degree[i].curr = degree[i].ori = d;
		degree[i].gtone = gtone; 
		//	cout<<"M="<<M<<" i="<<i<<" curr="<<degree[i].curr<<"  ori="<<degree[i].ori<<"  gtone="<<degree[i].gtone<<endl;
	}

    //ToString();
	/* step 1 */
	int _I, _U, r;
	int gtone_start = 0;

	_I = 0;
	_U = P;
   

 

	while (_I + _U < L) {


		int index, o;
retry:
		index = M; o = L; r = L;
		for (i = _I; i < M; i ++) {
			if ((gtone_start || (gtone_start==0 && degree[i].gtone==0)) && degree[i].curr > 0 && degree[i].curr <= r) {
				index = i; 
				if (degree[i].curr < r || (degree[i].curr == r && degree[i].ori < o)) {
					o = degree[i].ori;
					r = degree[i].curr;
				}
			}
		}

		if (index == M) {
		    if (gtone_start) goto retry;
			cout << "Cannot find enough rows to decode" << endl;
			PrintMatrix();
			return NULL;
		}

          
// 	  

 // cout<<endl;
      
		swap_row(_I, index);
      

	  
		k1 = k2 = 0;
		//cout<<"_I="<<_I<<"cols1[0]="<<cols1[0]<<endl;
	
		for (j = _I; j < L - _U; j++) {
			if (j < L - _U - r + 1) {
				if (A[_I][j].val != 0 ) {
					cols1[k1++] = j;
				}
			} else {
				if (A[_I][j].val == 0)
					cols2[k2++] = j;
			}
		}
      
		if (k1 != k2 + 1) {
			printf("Assert fail: %d!= %d + 1, _I=%d\n", k1, k2, _I);
			return NULL;
		
		}
           


		swap_col(_I,cols1[0]);
         
		
		for (j=0; j<k2; j++) {
		
			swap_col(cols2[j], cols1[j+1]);
		}
			
		if (A[_I][_I].val > 1) {
			unsigned char v = A[_I][_I].val;

			C1[_I]->div(v);
			for (j=_I; j<L; j++)
				A[_I][j].val = octdiv(A[_I][j].val, v);
		}
     
   
		for (i=_I+1;i<M;i++) {

			unsigned char v = A[i][_I].val;	
			if (v) {					
				A[i][_I].val = 0; 
				degree[i].curr--; 
				if (v > 1) degree[i].gtone--;
				for (j = L - _U - (r - 1); j < L; j++) {
					int oldv = A[i][j].val;
					A[i][j].val ^= octmul(v, A[_I][j].val);
					if (j < L - _U) {
						if (A[i][j].val > 0) {
							degree[i].curr ++;
							if (A[i][j].val > 1) degree[i].gtone ++;
						}
						if (oldv > 0) { 
							degree[i].curr --;
							if (oldv > 1) degree[i].gtone --;
						}
					}
				}
				C1[i]->muladd(C1[_I], v);
			} 
			for (j=L - _U - (r - 1); j < L - _U; j++)
				if (A[i][j].val) {
					degree[i].curr--;
					if (A[i][j].val > 1) degree[i].gtone--;
				}
		}
		_I++;
		_U += r - 1;

		
     
	}




























	delete[] cols1;
	delete[] cols2;



    if (!gaussian_elimination(_I, _I)) return NULL;
   	cout << "_I=" << _I << "_U=" << _U << endl;


	for (int jj=_I; jj<L; jj++)
		for (int i=0; i< _I; i++) {
			unsigned char v = A[i][jj].val;
			if (v) {
				A[i][jj].val = 0;
				C1[i]->muladd(C1[jj],v);
			}
		}
		
	
	for (i=0; i<L; i++) {
		C1[i]->esi = C[i]->esi;
	}

	cout<<"The generation maxtrixC1:" << endl;
	for (int i=0; i < (sizeof C1) ;i++) {
		    C1[i][0];    
	}

	for (i=0; i<L; i++) {
		Symbol *s;
		s = C[C1[i]->esi];
		C[C1[i]->esi] = C1[i];
		C1[i] = s; 
	}

	status = 3;

	return C;
}

/*Generate repair symbols using LT encoder,must be called after C is calculated
  The ESI of generated symbols starts from K
arg:
  count: number of repair symbols to generate
Return:
  the list of request count of repair symbols
 */
Symbol **Generators::generate_repairs(int count)
{
	int i, isi;
	
	if (status != 3) {
		cout << "Wrong call sequence! Generate intermediates before generate repairs" << endl;
		return NULL;
	}

	/* caller needs to free R */
	R = new Symbol*[count];

	/*
	for (i=0; i < count; i++) 
	{
		R[i].init(T);
	}
	*/

	for (i = K; i < K + count; i++) 
	{
		TuplS tupl;
		isi = i + K1 - K;
		if (isi < tupl_len) 
			tupl = Tuples[isi];
		else 
			tupl = Tupl(isi);
		R[i - K] = LTEnc(C, tupl);
	}

	status = 4;

	return R;
}

Symbol *Generators::recover_symbol(int x)
{
	Symbol *s;

	if (x >= K) {
		printf("try to recover non-source symbols!\n");
		return NULL;
	}

	TuplS tupl;
	tupl = Tuples[x];

	s = LTEnc(C, tupl);

	return s;
}

int Generators::getL() {
	return this->L;
}

int Generators::getK() {
	return this->K;
}


//打印矩阵用的
void Generators::PrintMatrix(void)
{
#ifdef DEBUG
	int i;
	cout<<"Press a number  to print the generation matrix:" << endl;
	cin >> i;
	for (i=0; i < M ;i++) {
		for (int j=0; j < L;j++) 
			printf("%2x ", A[i][j].val);
		cout<<endl;
	}
#endif
}





void Generators::ToString() {
	cout << "K=" << K <<"  I="<<I<< endl; 
	cout << "T=" << T <<"  J_K="<<J_K<<endl; 
	cout << "H=" << H <<"  X="<<X<< endl; 
	cout << "S=" << S <<"  H1="<<H1<< endl; 
	cout << "L=" << L <<"  N1="<<N1<< endl; 
	cout << "N=" << N <<"  P1="<<P1<< endl; 
	cout << "M=" << M <<"  U="<<U<< endl; 
	cout << "K1=" << K1 <<"  B="<<B<< endl; 
	cout << "W=" << W <<"  P="<<P<< endl; 
	cout << "Tuple_len=" << tupl_len << endl; 
	cout <<"Tuples:" << endl;
	for (int i=0;i < L; i++) {
		cout << "Tuple " << i << " d,a,b=" << Tuples[i].d << "," << Tuples[i].a << "," << Tuples[i].b << ",";
		cout << " d1,a1,b1=" << Tuples[i].d1 << "," << Tuples[i].a1 << "," << Tuples[i].b1 << endl;

	}

	cout<<"A maxtrix:" << endl;
	for (int i=0; i < M ;i++) {
		for (int j=0; j < L;j++) 
			printf("%3d ", A[i][j].val);
			//cout<<A[i][j].val<<" ";
		cout<<endl;
	}
  	cout<<"C maxtrix:"<< endl;
	for (int i=0; i < L;i++) {
		C[i]->print();
		cout<<endl;
	}
	cout<<"C1 maxtrix:" << endl;
	for (int i=0; i < M ;i++) {
		C1[i]->print();
		cout<<endl;
	}
	 cout<<"Abak maxtrix:" << endl;
	for (int i=0; i < M ;i++) {
		for (int j=0; j < L;j++) 
			printf("%3d ", Abak[i][j].val);
			//cout<<A[i][j].val<<" ";
		cout<<endl;
	}

	cout<<"source:" << endl;
	for (int i=0; i < (sizeof sources) ;i++) {
		for (int ii=0; ii < (sizeof sources[i]) ;ii++) {
			 printf("%c", sources[i][ii]);
		cout<<endl;
	}
	}
    cout<<"R maxtrix:" << endl;
	for (int i=0; i < (sizeof R) ;i++) {
		R[i]->print();
		cout<<endl;
	}

  

    


	
   
		

	
	
}

Generators::Generators() {
	status = 0;
}

// 应该是析构函数释放内存用的
Generators::~Generators() {
	int i;
	for (i=0; i< M; i++)
		delete C1[i];
	delete[] C1;

	for (i=0; i < M ;i++)
		delete []A[i];
	delete []A;
	for (i=0; i < M ;i++)
		delete []Abak[i];
	delete []Abak;
	delete[] Tuples;
	for (i=0; i < L; i++)
		delete C[i];
	delete[] C;
	// caller handle the delete of R
	//delete[] R;
	delete[] isi;
	delete[] degree;
}

/* section 5.7.2 */
unsigned char octmul(unsigned char u, unsigned char v)
{
	if (u == 0 || v == 0) return 0;
	if (v == 1) return u;
	if (u == 1) return v;
	return OCT_EXP[OCT_LOG[u] + OCT_LOG[v]];
}

unsigned char octdiv(unsigned char u, unsigned char v)
{
	if (u == 0) return 0;
	if (v == 1) return u;

	return OCT_EXP[OCT_LOG[u] - OCT_LOG[v] + 255];
}
package raptorQ

import (
	"fmt"
	"math"
)

type TuplS struct {
	d  int
	a  int
	b  int
	d1 int
	a1 int
	b1 int
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
     //c.Test()
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

	c.sources = source
	return true
}



// 中间符号生成
func (c *Common) Generate_intermediates() []Symbol {
	c.Test()
     var i,k1,k2,j int
	// 计算A^^-1生成中间符号C,使用高斯消元算法
	if c.status != 2 {
		return nil
	}
     
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
			//fmt.Println(i,c.degree[i].gtone,c.degree[i].curr)
			if (gtone_start != 0 || (gtone_start == 0 && c.degree[i].gtone == 0)) && c.degree[i].curr > 0 && c.degree[i].curr <= r {
				
				index = i
				
				if c.degree[i].curr < r || (c.degree[i].curr == r && c.degree[i].ori < o) {
					o = c.degree[i].ori
					r = c.degree[i].curr
				//	fmt.Println(r)
				
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
	







  
	

	 	



	if !c.gaussian_elimination(_I, _I) {
		return nil
	}
	fmt.Println("_I=",_I,"_U=",_U)

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



func (c *Common) LTEnc(C_L []Symbol, tupl TuplS) Symbol {
	a := tupl.a
	b := tupl.b
	d := tupl.d
	s := C_L[b]
	for j := 1; j < d; j++ {
		b = (b + a) % c.W
		s.xxor(&C_L[b])
	}
	a = tupl.a1
	b = tupl.b1
	d = tupl.d1
	for b >= c.P {
		b = (b + a) % c.P1
	}
	s.xxor(&C_L[c.W+b])

	for j := 1; j < d; j++ {
		for b >= c.P {
			b = (b + a) % c.P1
		}
		s.xxor(&C_L[c.W+b])
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
			fmt.Printf("%4d", c.Abak[i][j])
		}
		fmt.Println()
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
