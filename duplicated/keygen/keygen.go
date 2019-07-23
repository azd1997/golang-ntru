package keygen

import (
	"fmt"
	"github.com/azd1997/golang-ntru/duplicated"
)

//多项式环上的数学运算

//整数环Z，整数N>=2，R表示多项式截断环，R=Z[X]/(X^N-1)
//对于任意正整数q，Rq表示模q的多项式截断环，Rq=(Z/qZ)[X]/(X^N-1)
//q为素数时，Rq具有可逆性，可逆性指对于F及其模q逆Fq，有F·Fq=1 mod q

//NTRU算法需要3个整数参数（N,p,q）和4个多项式环（f,g,r,m），其中(f,g,r)根据参数随机生成
//要求gcd(p,q)=1（即二者最大公约数为1），且q远大于p
//NTRU-1998中系数为三元整数(-1,0,1)

//NTRU标准参数
//安全等级	N	q	p	df	dg	dr
//低		107	64	3	15	12	5
//中		167	128	3	61	20	18
//高		263	128	3	50	24	16
//极高		503	256	3	216	72	55
//考虑到常用的这些参数大小，将其保存为uint16类型（0~65523）

//密钥产生
//根据N随机选择两个多项式f,g。为使对模p模q的乘逆存在，f应满足gcd(f,pq)=1
//以Fp、Fq表示相应的模逆，则f·Fp=1 mod p; f·Fq=1 mod q
//计算h: h = q·Fp·g mod p
//私人密钥为一对多项式环(f,Fq)，公开密钥为h

//加密过程
//设有明文m，根据参数dr随机选择多项式r
//使用公钥加密得到密文e： e = (r·h+m) mod p

//解密过程
//使用私钥(f,Fq)解密： a = f·e mod p ; b = a mod q ; m = Fq · b mod q


//以多项式系数数组切片定义
//a = a[0]+a[1]x+a[2]x^2+...+a[N-1]x^(N-1)
//由于系数设置为三元组（-1,0,1）

const (
	Ntru_N = 107
	Ntru_q = 64
	Ntru_p = 3
	Ntru_df = 15
	Ntru_dg = 12
	Ntru_dr = 5
)

type Polynomial [Ntru_N]uint8;

type NtruCipher struct {
	N uint16
	p uint16
	q uint16
	df uint16
	dg uint16
	dr uint16

	f Polynomial
	g Polynomial
	r Polynomial
	fp Polynomial
	fq Polynomial	//模q逆
	h Polynomial
}

func (ntru *NtruCipher) Init(N,p,q,df,dg,dr uint16) error {
	//初始化N,p,q,dg,dg,dr六项整数参数
	ntru.N = N
	ntru.p = p
	ntru.q = q
	ntru.df = df
	ntru.dg = dg
	ntru.dr = dr

	//初始化f/g/r为系数全0
	var i uint16
	for i = 0; i < N; i++ {
		ntru.f[i] = 0
		ntru.g[i] = 0
		ntru.r[i] = 0
		ntru.fp[i] = 0
		ntru.fq[i] = 0
		ntru.h[i] = 0
	}

	return nil
}

func (ntru *NtruCipher) generateRandomPolynomial(d uint64) error {

	d = d +d - 1
	//j := -1


	return nil
}



/************************************/
//Name: Multi
//Inputs:
//Outputs:
//Description: 多项式环的乘法
/************************************/
func Multi(a,b []int, n int) []int {

	//na := len(a)
	//	//nb := len(b)
	//	//
	//	//if na != nb {
	//	//	fmt.Println("两个输入数组需要等长，且长度符合范围！")
	//	//	return nil
	//	//}
	//	//n := na

	var c = make([]int, n)	//指定初始长度的切片
	var i, j int
	//举例
	//a=[1,2,3]，即a=1+2x+3x^2 ; b=[4,5,6]，即b=4+5x+6x^2
	//a multi b = 1*(4+5x+6x^2) + 2x*(4+5x+6x^2) + (3x^2)*(4+5x+6x^2)
	// = (4+5x+6x^2) + (8x+10x^2+12x^3) + (12x^2+15x^3+18x^4)
	// = (4+5x+6x^2) + (8x+10x^2+12) + (12x^2+15+18x)
	// = (4+12+15) + (5+8+18)x + (6+10+12)x^2
	// = 31 + 31x + 28x^2	//这里显然c[0] = a[0]b[0]+a[2]b[1]+a[1]b[2]
	// 当计算c[0]时，i=0, k=1, c[i]=0, for循环{
	// 1、j=0,k=0,c[0]=0+a[0]b[0]=4,k=3
	// 2、j=1,k=2,c[0]=4+a[2]b[1]=19
	// 3、j=2,k=1,c[0]=19+a[1]b[2]=31}
	// 当计算c[1]时，i=1, k=2, c[i]=0, for循环{
	// 1、j=0,k=1,c[0]=0+a[1]b[0]=8
	// 2、j=1,k=0,c[0]=8+a[0]b[1]=13,k=3
	// 3、j=2,k=2,c[0]=13+a[2]b[2]=31}
	// 当计算c[2]时，i=2, k=3, c[i]=0, for循环{
	// 1、j=0,k=2,c[0]=0+a[2]b[0]=12
	// 2、j=1,k=1,c[0]=12+a[1]b[1]=22
	// 3、j=2,k=0,c[0]=22+a[0]b[2]=28}
	for i=0; i<n; i++ {		//i表示c数组的第i位，c[i]也即x^(i-1)项系数
		k := i + 1
		c[i] = 0	//先全部赋0，go数组元素默认为nil
		for j=0; j<n; j++ {
			k--
			c[i] = c[i] + a[k]*b[j]
			if k==0 {
				k = n
			}
		}
	}

	return c
}


/************************************/
//Name: mod
//Inputs:
//Outputs:
//Description: 整型取模
/************************************/
//func mod(a,b int) int {
//	//这里不对输入的a,b做检查，假定其为整型数
//
//	b = math.Abs()
//
//	//c := a/b并向无穷小方向舍入
//	//用a的绝对值填充c
//	c := a
//	if a<0 {
//		c = -a
//	}
//
//	for
//
//	//r := a-c*b
//
//	return r
//}


/************************************/
//Name: Amodp
//Inputs:
//Outputs:
//Description: 多项式的对p求模
/************************************/
func Amodp(a []int, n,p int) []int {
	//通常p设置为3，则p/2=1
	p1 := p/2

	//关于取模与求余  a%b
	//计算过程：1、c=a/b; 2、r=a-c*b
	//求模与求余的区别在于第一步的计算，求余在取c时向0方向舍入；求模取c时则向无穷小方向舍入
	//例如5%2,模和余计算都是1；
	//但-5%2： 1.模 c=-5/2=-3；r=a-c*b=1；2.余 c=-5/2=-2；r=a-c*b=-1
	//一句话描述：求模结果与b符号一致；求余结果与a符号一致


	var i int
	for i=0; i<n; i++ {
		//a[i] = int(math.Mod(float64(a[i]), float64(p)))
		//从名称定义来说，这里是求余计算。很奇怪,这并不是求模定义
		a[i]=a[i]%p

		//由于p=3，其余数一定在[-2,-1,0,1,2],经过下列处理则将值束缚在[-1,0,1]
		// 2=> -1; -2=>1
		if a[i] > p1 {
			a[i] = a[i] - p
		}
		if a[i] < -p1 {
			a[i] = a[i] + p
		}
	}
	return a

}


/************************************/
//Name: BmodXn
//Inputs:
//Outputs:
//Description: 多项式的特殊求模
/************************************/
func BmodXn(b []int, n int) []int {
	//这里提供的n是X^n-1的n
	//例如x^2-1，n=2，则最终得到的模一定形如px+q，也就是说其系数数组长度为n

	//原理：令a=X^n-1，根据带余除法定理总存在q,r<-R，使得b=a*q+r
	//同时有deg(b)>=deg(a);deg(a)>=deg(r)或r=0
	//由于a是一个特殊多项式，根据带余除法可将r=b mod a 运算变换成加法运算
	//r(i)=b(i)+b(i+n)+b(i+2n)+...

	degb := len(b) - 1

	var i int
	var r = make([]int, n)
	for i=0; i<n; i++ {
		k := i
		for k<=degb {
			r[i] += b[k]
			k = k + n
		}
	}

	return r
}

/************************************/
//Name: AB_1modp
//Inputs:
//Outputs:
//Description: 求解a(X)·b(x)=1 mod p； 返回b(x)，也称作a(x)的模逆或乘逆。p要求为素数
/************************************/
func AB_1modp(a []int, n,p,ng int) []int {

	//1.初始化k=0,b=[1,0,...],c=[0,0,...],f=a,g=[-1,1]
	k := 0
	var b = make([]int, n)
	var c = make([]int, n)
	var i int
	for i=0;i<n;i++ {
		b[i] = 0
		c[i] = 0
	}
	b[0] = 1
	f := a
	//g=[-1,0,...,0,1]，取ng为1，则g=[-1,1]
	var g = make([]int, ng+1)
	g[0] = -1
	g[ng] = 1
	for i=1;i<ng;i++ {
		g[i] = 0
	}



	for {
		//2.f[x]=f[x]/x; c[x]=c[x]*x;k=k+1
		//TODO:记得修改f长度n以适应其变化
		//while f(0)=0且阶次大于0
		if f[0]==0 && len(f)-1 !=0 {
			//f[x]=f[x]/x 即系数向量左移一位，记住多项式环乘法规则，是循环左移
			for i=0;i<n-1;i++ {
				f[i] = f[i+1]
			}
			f = f[0:n-2]
			//c[x]=c[x]*x 即系数向量循环右移一位。注意这里是多项式环上乘法
			tmp := c[n-1]
			for i=n-1;i>0;i-- {
				c[i] = c[i-1]
			}
			c[0] = tmp
			//k=k+1
			k++
		}
		//3.
		//if f(x)=+-1，这意味着f长度为1
		if (f[0] == 1 || f[0] == -1) && len(f) == 1 {
			for item := range b {item *= f[0]}
			b = Amodp(b, n, p)
			//k1 = (n-k) mod p
			k1 := (n-k) % p
			if k1<0 {k1 += n}
			//b(X)=(x^k1)*b(x) mod (x^ng-1)
			for i=0;i<k1-1;i++ {
				b = append(b, b[i])
			}	//现在b长度为n+k1
			b = b[k1:]
			b = BmodXn(b, ng)
			return b
		}
		//4.
		//if deg(f)<deg(g)，交换f和g，b和c
		if len(f) < len(g) {
			tmp1 := f
			f = g
			g = tmp1
			tmp1 = b
			b = c
			c = tmp1
		}

		//5.if f(0)==g(0) {f(x)=f(x)-g(x) mod p; b(x)=b(x)-c(x) mod p}
		//TODO:这里默认了ng<n，后面考虑对输入参数做判断检验
		gp := Amodp(g, ng, p)
		cp := Amodp(c, n, p)
		if f[0] == g[0] {
			for i=0;i<ng;i++ {
				f[i] -= gp[i]
			}
			for i=0;i<n;i++ {
				b[i] -= cp[i]
			}
		} else {
			for i=0;i<ng;i++ {
				f[i] += gp[i]
			}
			for i=0;i<n;i++ {
				b[i] += cp[i]
			}
		}

		break
	}


	return b
}

/************************************/
//Name: generatePolynomial
//Inputs: 系数数组长度、系数1的个数、系数-1的个数
//Outputs:
//Description: 随机产生多项式
/************************************/
//func generatePolynomial(N, d1, d_1 int) []int {
//
//	var r = make([]int,N)
//	rand
//
//	return r
//}

/************************************/
//Name: Conv1
//Inputs:
//Outputs:
//Description: 多项式环上的卷积（乘法）,计算a(x)*b(x) mod q
/************************************/
func Conv1(n,q int, a,b []int) []int {
	var i,j int
	var c = make([]int,n)
	for i=0; i<n; i++ {		//i表示c数组的第i位，c[i]也即x^(i-1)项系数
		k := i + 1
		c[i] = 0	//先全部赋0，go数组元素默认为nil
		for j=0; j<n; j++ {
			k--
			c[i] = c[i] + a[k]*b[j]
			if k==0 {
				k = n
			}
		}
		//
		c[i] = c[i] % q
	}

	return c
}

/************************************/
//Name: Conv2
//Inputs:
//Outputs:
//Description: 多项式环上的HW快速卷积（乘法）,计算a(x)*b(x) mod q
/************************************/
func Conv2(n,q int, a,b []int) []int {
	//针对a(x)系数为[0,1]二元组
	hw := 0
	var i,k int
	var c = make([]int,2*n)
	for i=0;i<2*n;i++ {c[i]=0}
	var t = make([]int, n)


	for i=0; i<n; i++ {		//i表示c数组的第i位，c[i]也即x^(i-1)项系数
		if a[i] != 0 {
			t[hw] = i
			hw++
		}
	}

	for i=0;i<=hw;i++ {
		for k = 0; k < n; k++ {
			c[k+t[i]] += b[k]
		}
	}

	for i=0;i<n;i++ {
		c[i] = (c[i] + c[i+n]) % q
	}

	return c
}


/************************************/
//Name: PMultiAx
//Inputs:
//Outputs:
//Description: 计算b(x)=p*a(x)，其中p=x+2.这是NTRU-2001中的.
//				其实完全可以用Multi()实现
/************************************/
func PMultiAx(n int, a []int) []int {
	var b = make([]int,n)
	b[0] = a[n-1]

	var i int
	for i=0;i<n;i++ {
		b[i] = a[i-1]
	}
	for i=0;i<n;i++ {
		b[i] += 2*a[i]
	}

	return b
}

/************************************/
//Name: AmodP2
//Inputs:
//Outputs:
//Description: 计算a(x) mod p	注意这里的p=x+2
/************************************/
func AmodP2(n,q int, a []int) []int {
	n2 := 2*n
	k :=0
	var i int

	//a:=[]int{14,13,9,15,-14,15,16} n:=7 q:=32

	//a[i]为负则加q，为了保证解密结果的正确性
	for i=0;i<n;i++ {
		if a[i]<0 { a[i] += q }
	}

	i = 0
	for {

		//分解 a[i]=u+2v； i=0,a(i)=14,v=7,u=0
		v := a[i]/2
		u := a[i] - 2*v
		//变换a[i]、a[i+1]、a[i+2]、a[i+n]
		//i=0时，经历下边这部分，a=[0,20,16,18,15,16,0]
		a[i] = u
		a[i+1] += v
		a[i+2] += v
		//a[n+i] = a[i]
		if n+i>=len(a) {
			a = append(a,a[i])
		} else {
			a[n+i] = a[i]
		}

		i++		//现在i=1
		//w取a[i]/2和a[i+1]的较小值,w=min{20/2,16}=10
		w := a[i]/2
		if w>a[i+1] { w = a[i+1] }

		a[i] = a[i] -2*w
		a[i+1] = a[i+1] - w
		//a[n+i] = a[i]
		if n+i>=len(a) {
			a = append(a,a[i])
		} else {
			a[n+i] = a[i]
		}
		//a[n+i+1] = a[i+1]
		if n+i+1 >= len(a) {
			a = append(a,a[i+1])
		} else {
			a[n+i+1] = a[i+1]
		}
		//a=[0,0,6,18,15,16,0，0,6]

		fmt.Println(a)

		//符合条件时继续下一次循环，否则退出循环
		if i<n || (a[i] != 0 && a[i] != 1) || ((a[i+1] != 0 && a[i+1] != 1) && i<n2) {
			continue
		}
		fmt.Println(i)
		break
		//退出循环时，a所有项变为（0,1），得到长度为2N的系数数组
		//本例中，i=9时退出循环
	}

	//取a的第N~2N-1位输出
	var m = make([]int,n)
	for k=0;k<n;k++ {
		m[k] = a[n+k]
	}

	//如果循环次数超过2N-1次，需要考虑//a[n+i] = a[i]//a[n+i+1] = a[i+1]
	//的滞后作用，并对m[0]、m[1]重新赋值
	if i>n2 {
		m[0] = a[n2]
		m[1] = a[n2+1]
	}

	return m
}

/************************************/
//Name: Keygen
//Inputs:
//Outputs:
//Description: 生成公私钥
/************************************/
func Keygen() ([]int,[]int,[]int) {
	//N=11,q=32,p=2+x,df=4,dg=5,dr=5
	n := 11
	q := 32
	df := 4
	dg := 5
	//dr := 5

	//f1对应的是算法中的F  F = random_f(n,df)
	f1 := duplicated.Random_f(n, df)

	//f = 1 + p*F
	f := PMultiAx(n, f1)
	f[0] += 1

	// f*fp = 1 mod p 求fp
	//TODO:
	fq := f

	// g = random_f(n,dg)
	g := duplicated.Random_f(n,dg)

	// h = p* g * fq (mod q)
	h := PMultiAx(n, g)
	h = Multi(h, fq, n)
	h = Amodp(h,n,q)

	return f,fq,h
}

/************************************/
//Name: Encrypt
//Inputs:
//Outputs:
//Description: 加密
/************************************/
func Encrypt(m,h []int, n int) []int {
	var e = make([]int, n)

	//1.产生多项式r
	dr := 5
	q := 32
	r := duplicated.Random_f(n,dr)

	//2. r*h
	e = Multi(r, h, n)
	//3.m mod q
	m = Amodp(m, n, q)

	//4.e = r*h+m (mod q)
	var i int
	for i=0;i<n;i++ {
		e[i] = e[i] + m[i]
	}

	return e
}

/************************************/
//Name: Decrypt
//Inputs:
//Outputs:
//Description: 加密
/************************************/
func Decrypt(e,f,fq []int, n int) []int {
	var m = make([]int, n)
	q := 32
	// a = f*e(mod q)
	a := Amodp(e, n, q)
	a = Multi(f, a, n)

	// m = (mod)p

	return m
}