package NTRU_2001

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

/************************************/
//Name: RandomPoly
//Inputs: 系数数组长度n；系数1的个数d
//Outputs: 多项式系数数组
//Description: 生成系数为(-1,0,1)的随机多项式(NTRU-1998)
/************************************/
func RandomPoly(n,d1,d_1 int) []int {
	//d1表示系数1个数；d_1表示系数-1个数
	var i int
	//初始化poly全0
	var poly = make([]int, n)
	for i=0;i<n;i++ {poly[i]=0}

	//根据时间生成不同的随机源，根据随机源生成不同的随机数
	random := rand.New(rand.NewSource(time.Now().Unix()))

	//系数1的随机生成
	k := 0
	for k<d1 {
		//这个随机范围最好足够大（要比你可能选用的N值大）
		i = int(math.Floor(float64(n * random.Intn(1000) / 1000)))
		if poly[i] == 0 {
			poly[i] = 1
			k++
		}
	}

	//系数-1的随机生成
	k = 0
	for k<d_1 {
		//这个随机范围最好足够大（要比你可能选用的N值大）
		i = int(math.Floor(float64(n * random.Intn(1000) / 1000)))
		if poly[i] == 0 {
			poly[i] = -1
			k++
		}
	}

	return poly
}

/************************************/
//Name: Random_gr
//Inputs: 系数数组长度n；系数1的个数d
//Outputs: 多项式系数数组
//Description: 生成系数为(-1,0,1)的随机多项式(NTRU-1998)
/************************************/
func Random_gr(n,d int) []int {
	//注意，生成的f其1系数个数比-1系数多一个
	d = d + d
	j := -1
	var i int
	var f = make([]int, n)
	for i=0;i<n;i++ {f[i]=0}

	//根据时间生成不同的随机源，根据随机源生成不同的随机数
	random := rand.New(rand.NewSource(time.Now().Unix()))
	k := 0

	for k<d {
		i = int(math.Floor(float64(n * random.Intn(100) / 100)))
		if f[i] == 0 {
			j = -j
			f[i] = j
			k++
		}
	}

	return f
}


/************************************/
//Name: Random_f
//Inputs: 系数数组长度n；系数1的个数d
//Outputs: 多项式系数数组
//Description: 生成系数为(-1,0,1)的随机多项式(NTRU-1998)
/************************************/
func Random_f(n,d int) []int {
	//注意，生成的f其1系数个数比-1系数多一个
	d = d+d-1
	j := -1
	var i int
	var f = make([]int, n)
	for i=0;i<n;i++ {f[i]=0}

	//根据时间生成不同的随机源，根据随机源生成不同的随机数
	random := rand.New(rand.NewSource(time.Now().Unix()))
	k := 0

	for k<d {
		i = int(math.Floor(float64(n * random.Intn(100) / 100)))
		if f[i] == 0 {
			j = -j
			f[i] = j
			k++
		}
	}

	return f
}


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
//Name: ConvHw1
//Inputs:
//Outputs:
//Description: 多项式环上的HW快速卷积（乘法）,计算a(x)*b(x) mod q
/************************************/
func ConvHw1(n,q int, a,b []int) []int {
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
//Name: ConvHw1
//Inputs:
//Outputs:
//Description: 多项式环上的HW快速卷积（乘法）,计算a(x)*b(x) mod q
/************************************/
func ConvHw2(n,q int, a,b []int) []int {
	//针对a(x)系数为[-1,0,1]三元组
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
			//根据a[t[i]]值为1还是-1确定在c中加还是减b[k]
			if a[t[i]] == 1 {
				c[k+t[i]] += b[k]
			} else {
				c[k+t[i]] -= b[k]
			}
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
//Name: DegOfPoly
//Inputs:
//Outputs:
//Description: 求多项式阶次
/************************************/
func DegOfPoly(a []int) int {
	i := len(a) - 1
	for i>=0 {
		if a[i] != 0 {
			break
		}
		i--
	}
	return i
}

/************************************/
//Name: AB_1modp_3
//Inputs:
//Outputs:
//Description: 求解a(X)·b(x)=1 mod p； 返回b(x)，也称作a(x)的模逆或乘逆。p要求为素数
/************************************/
func AB_1modp_3(a []int, n,ng int) []int {

	//1.初始化k=0,b=[1,0,...],c=[0,0,...],f=a,g=[-1,1],p=3
	k := 0
	p := 3
	var b = make([]int, n)
	var c = make([]int, n)
	var i int
	for i=0;i<n;i++ {
		b[i] = 0
		c[i] = 0
	}
	b[0] = 1
	f := a
	//g=[-1,0,...,0,1,0,0,...]，取ng为1，则g=[-1,1,0,0,...]，默认ng<=n
	var g = make([]int, n)
	g[0] = -1
	g[ng] = 1
	for i=0;i<n;i++ {
		switch i {
		case 0:
			g[i] = -1
		case ng:
			g[i] = 1
		default:
			g[i] = 0
		}
	}

	for {
		//2.f[x]=f[x]/x; c[x]=c[x]*x;k=k+1
		//while f(0)=0且阶次大于0
		if f[0]==0 && DegOfPoly(f)>0 {
			//f[x]=f[x]/x 即系数向量左移一位，记住多项式环乘法规则，是循环左移，缺位补0
			for i=0;i<n-1;i++ {
				f[i] = f[i+1]
			}
			f[n-1] = 0
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
		//if f(x)=+-1，
		//bool_f := true
		//for i=1;i<n;i++ {
		//	if f[i] != 0 {
		//		bool_f = false
		//	}
		//}
		//if (f[0] == 1 || f[0] == -1) && (bool_f == true) {
		if DegOfPoly(f) == 0 {
			for i = 0; i < n; i++ {
				b[i] *= f[0]
			}
			b_copy := b
			for i = 0; i < k; i++ {
				b[i] = b_copy[n-k+i]
			}
			for i = k; i < n; i++ {
				b[i] = b_copy[i-k]
			}
			b = BmodXn(b, n)
		}

		//4.
		//if deg(f)<deg(g)，交换f和g，b和c
		if DegOfPoly(f) < DegOfPoly(g) {
				//tmp1 := f
				//f = g
				//g = tmp1
				//tmp1 = b
				//b = c
				//c = tmp1
				f,g = g,f
				b,c = c,b
			}

		//5.if f(0)==g(0) {f(x)=f(x)-g(x) mod p; b(x)=b(x)-c(x) mod p}
		if f[0] == g[0] {
			for i=0;i<n;i++ {
				f[i] -= g[i]
				f[i] %= p
				b[i] -= c[i]
				b[i] %= p
			}
		} else {
			for i=0;i<n;i++ {
				f[i] += g[i]
				f[i] %= p
				b[i] += c[i]
				b[i] %= p
			}
		}
		break
	}


	return b
}


func ModPrime(x,prime int) (ret int) {
	ret = x % prime
	if ret<0 {
		ret += prime
	}
	return
}


/************************************/
//Name: AB_1modq_2N
//Inputs:
//Outputs:
//Description: 求解a(X)·b(x)=1 mod p； 返回b(x)，也称作a(x)的模逆或乘逆。p要求为素数
/************************************/
func AB_1modq_2N(a []int, n,ng,q int) []int {
	//1.初始化k=0,b=[1,0,...],c=[0,0,...],f=a,g=[-1,1],q为2的指数幂
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
	//g=[-1,0,...,0,1,0,0,...]，取ng为1，则g=[-1,1,0,0,...]，默认ng<=n
	var g = make([]int, n)
	g[0] = -1
	g[ng] = 1
	for i=0;i<n;i++ {
		switch i {
		case 0:
			g[i] = -1
		case ng:
			g[i] = 1
		default:
			g[i] = 0
		}
	}

	for {
		//2.f[x]=f[x]/x; c[x]=c[x]*x;k=k+1
		//while f(0)=0且阶次大于0
		if f[0]==0 && DegOfPoly(f)>0 {
			//f[x]=f[x]/x 即系数向量左移一位，记住多项式环乘法规则，是循环左移，缺位补0
			for i=0;i<n-1;i++ {
				f[i] = f[i+1]
			}
			f[n-1] = 0
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
		//if f(x)=+-1，
		//bool_f := true
		//for i=1;i<n;i++ {
		//	if f[i] != 0 {
		//		bool_f = false
		//	}
		//}
		//if (f[0] == 1 || f[0] == -1) && (bool_f == true) {
		if DegOfPoly(f) == 0 {
			for i = 0; i < n; i++ {
				b[i] *= f[0]
			}
			b_copy := b
			for i = 0; i < k; i++ {
				b[i] = b_copy[n-k+i]
			}
			for i = k; i < n; i++ {
				b[i] = b_copy[i-k]
			}
			b = BmodXn(b, n)
		}

		//4.
		//if deg(f)<deg(g)，交换f和g，b和c
		if DegOfPoly(f) < DegOfPoly(g) {
			tmp1 := f
			f = g
			g = tmp1
			tmp1 = b
			b = c
			c = tmp1
		}

		//5.
		for i=0;i<n;i++ {
				f[i] += g[i]
				f[i] %= 2
				b[i] += c[i]
				b[i] %= 2
		}


		v := 2
		if v<q {
			v = 2*v
			for i=0;i<n;i++ {

			}
		}
		break
	}


	return b
}

/************************************/
//Name: AB_1modp
//Inputs:
//Outputs:
//Description: 求解a(X)·b(x)=1 mod p； 返回b(x)，也称作a(x)的模逆或乘逆。p要求为素数
/************************************/
//func AB_1modp(a []int, n,p int) []int {
//
//	//1.初始化k=0,b=[1,0,...],c=[0,0,...],f=a,g=[-1,1]
//	k := 0
//	var b = make([]int, n+1)
//	var c = make([]int, n+1)
//
//	for i:=0;i<=n;i++ {
//		b[i] = 0
//		c[i] = 0
//	}
//	b[0] = 1
//
//	var f = make([]int, n+1)
//	for i:=0;i<n;i++ {
//		f[i] = a[i]
//	}
//	//g=[-1,0,...,0,1]，取ng为1，则g=[-1,1]
//	var g = make([]int, n+1)
//	g[0] = -1
//	g[n] = 1
//	for i:=1;i<n;i++ {
//		g[i] = 0
//	}
//	deg_g := n
//
//
//
//	for {
//		//2.f[x]=f[x]/x; c[x]=c[x]*x;k=k+1
//		//while f(0)=0且阶次大于0
//		if f[0]==0 && len(f)-1 !=0 {
//			//f[x]=f[x]/x 即系数向量左移一位，记住多项式环乘法规则，是循环左移
//			f0 := f[0]
//			for i:=0;i<n-1;i++ {
//				f[i] = f[i+1]
//			}
//			f[] = f0
//			//c[x]=c[x]*x 即系数向量循环右移一位。注意这里是多项式环上乘法
//			tmp := c[n-1]
//			for i=n-1;i>0;i-- {
//				c[i] = c[i-1]
//			}
//			c[0] = tmp
//			//k=k+1
//			k++
//		}
//		//3.
//		//if f(x)=+-1，这意味着f长度为1
//		if (f[0] == 1 || f[0] == -1) && len(f) == 1 {
//			for item := range b {item *= f[0]}
//			b = Amodp(b, n, p)
//			//k1 = (n-k) mod p
//			k1 := (n-k) % p
//			if k1<0 {k1 += n}
//			//b(X)=(x^k1)*b(x) mod (x^ng-1)
//			for i=0;i<k1-1;i++ {
//				b = append(b, b[i])
//			}	//现在b长度为n+k1
//			b = b[k1:]
//			b = BmodXn(b, ng)
//			return b
//		}
//		//4.
//		//if deg(f)<deg(g)，交换f和g，b和c
//		if len(f) < len(g) {
//			tmp1 := f
//			f = g
//			g = tmp1
//			tmp1 = b
//			b = c
//			c = tmp1
//		}
//
//		//5.if f(0)==g(0) {f(x)=f(x)-g(x) mod p; b(x)=b(x)-c(x) mod p}
//		//TODO:这里默认了ng<n，后面考虑对输入参数做判断检验
//		gp := Amodp(g, ng, p)
//		cp := Amodp(c, n, p)
//		if f[0] == g[0] {
//			for i=0;i<ng;i++ {
//				f[i] -= gp[i]
//			}
//			for i=0;i<n;i++ {
//				b[i] -= cp[i]
//			}
//		} else {
//			for i=0;i<ng;i++ {
//				f[i] += gp[i]
//			}
//			for i=0;i<n;i++ {
//				b[i] += cp[i]
//			}
//		}
//
//		break
//	}
//
//
//	return b
//}

/************************************/
//Name: AB_1modp
//Inputs:
//Outputs:
//Description: 求解a(X)·b(x)=1 mod p； 返回b(x)，也称作a(x)的模逆或乘逆。p要求为素数
/************************************/
func Invert(a []int, n,prime int) []int {

	//1.初始化k=0,b=[1,0,...],c=[0,0,...],f=a,g=[-1,1]
	k := 0
	var b = make([]int, n+1)
	var c = make([]int, n+1)
	for i:=0;i<=n;i++ {
		b[i] = 0
		c[i] = 0
	}
	b[0] = 1

	var f = make([]int, n+1)
	for i:=0;i<n;i++ {
		f[i] = a[i]
	}
	//g=[-1,0,...,0,1]，取ng为1，则g=[-1,1]
	var g = make([]int, n+1)
	g[0] = -1
	g[n] = 1
	for i:=1;i<n;i++ {
		g[i] = 0
	}
	deg_g := n
	deg_f := DegOfPoly(f)

	fmt.Println("Initialized")
	fmt.Println(f)

	for {
		//2.f[x]=f[x]/x; c[x]=c[x]*x;k=k+1
		//while f(0)=0且阶次大于0
		for f[0]==0 && deg_f>0 {
			//f[x]=f[x]/x 即系数向量左移一位，记住多项式环乘法规则，是循环左移
			deg_f--	//注意这里必须使用这个变量，不能采用DegOfPoly,原因在于这里对f/x的处理方式
			f0 := f[0]
			for i:=0;i<n-1;i++ {
				f[i] = f[i+1]
			}
			f[n] = f0
			//c[x]=c[x]*x 即系数向量循环右移一位。注意这里是多项式环上乘法
			tmp := c[n]
			for i:=n;i>0;i-- {
				c[i] = c[i-1]
			}
			c[0] = tmp
			//k=k+1
			k++
			fmt.Println("Step1:", f)
		}
		//3.
		//if f(x)=+-1，这意味着f阶次为1,但是在这里f仍为n+1长度
		if deg_f == 0 && f[0] != 0 {
			//for item := range b {item *= f[0]}
			//b = Amodp(b, n, prime)
			////k1 = (n-k) mod p
			//k1 := (n-k) % prime
			//if k1<0 {k1 += n}
			////b(X)=(x^k1)*b(x) mod (x^ng-1)
			//for i:=0;i<k1-1;i++ {
			//	b = append(b, b[i])
			//}	//现在b长度为n+k1
			//b = b[k1:]
			//b = BmodXn(b, n)
			//return b

			//Make sure there is a solution, return nil if a is not invertible.
			// b(X) = f[0]inv * b(X) mod p
			// return X^(N-k) * b

			shift := n - k
			shift %= n
			if shift < n {
				shift += n
			}
			ret := make([]int, n)
			for i := range ret {
				ret[(i+shift)%n] = ModPrime(f[0]*b[i], prime)
			}
			fmt.Println("Step2:", ret)
			return ret
		}
		//4.
		//if deg(f)<deg(g)，交换f和g，b和c
		if deg_f < deg_g {
			f,g = g,f
			b,c = c, b
			deg_f, deg_g = deg_g, deg_f
		}

		//5.if f(0)==g(0) {f(x)=f(x)-g(x) mod p; b(x)=b(x)-c(x) mod p}

		gp := Amodp(g, n+1, prime)
		cp := Amodp(c, n+1, prime)
		if f[0] == g[0] {
			for i:=0;i<=n;i++ {
				f[i] -= gp[i]
			}
			for i:=0;i<=n;i++ {
				b[i] -= cp[i]
			}
		} else {
			for i:=0;i<=n;i++ {
				f[i] += gp[i]
			}
			for i:=0;i<=n;i++ {
				b[i] += cp[i]
			}
		}


	}

}

/************************************/
//Name: Keygen
//Inputs:
//Outputs:
//Description: 生成公私钥
/************************************/
//func Keygen() ([]int,[]int,[]int) {
//	//N=11,q=32,p=2+x,df=4,dg=5,dr=5
//	n := 11
//	q := 32
//	df := 4
//	dg := 5
//	//dr := 5
//
//	//f1对应的是算法中的F  F = random_f(n,df)
//	f1 := Random_f(n, df)
//
//	//f = 1 + p*F
//	f := PMultiAx(n, f1)
//	f[0] += 1
//
//	// f*fp = 1 mod p 求fp
//	//TODO:
//	fq := f
//
//	// g = random_f(n,dg)
//	g := Random_f(n,dg)
//
//	// h = p* g * fq (mod q)
//	h := PMultiAx(n, g)
//	h = Multi(h, fq, n)
//	h = Amodp(h,n,q)
//
//	return f,fq,h
//}

/************************************/
//Name: Encrypt
//Inputs:
//Outputs:
//Description: 加密
/************************************/
//func Encrypt(m,h []int, n int) []int {
//	var e = make([]int, n)
//
//	//1.产生多项式r
//	dr := 5
//	q := 32
//	r := Random_f(n,dr)
//
//	//2. r*h
//	e = Multi(r, h, n)
//	//3.m mod q
//	m = Amodp(m, n, q)
//
//	//4.e = r*h+m (mod q)
//	var i int
//	for i=0;i<n;i++ {
//		e[i] = e[i] + m[i]
//	}
//
//	return e
//}

/************************************/
//Name: Decrypt
//Inputs:
//Outputs:
//Description: 加密
/************************************/
//func Decrypt(e,f,fq []int, n int) []int {
//	var m = make([]int, n)
//	q := 32
//	// a = f*e(mod q)
//	a := Amodp(e, n, q)
//	a = Multi(f, a, n)
//
//	// m = (mod)p
//
//	return m
//}


//密钥生成
//输入参数N,p,q,df,dg,dr
//随机生成g,r,以及私钥f
//检查逆元fp,fq是否存在？不存在则重新生成私钥f   f * fp = 1 mod p ; f * fq = 1 mod q
//若存在则计算公钥h  h = g * fq (mod q)  h= p·fq * g (mod q)	//g/fq也应该保密
//文件存储f,fp,h

//加密
//从文件取出公钥h
//将明文消息转为二进制明文M （[]int） ，再转为三进制m（针对p=3）
//加密操作： e = r*h + m(mod q)

//解密
//从文件取出私钥对(f,fp)
//解密操作： 1.a = f * e(mod q) ;
// 2.对多项式a系数作模q处理，调整至(-q/2,q/2)区间内
// 3.d = fp * a (mod p)

//区间调整与解密失败
// a=f*e= f*(r*h+m(mod q)) = f*r*h + f*m (mod q)
// = f*r* (p·fq*g) + f*m (mod q) = p·r*g+f*m (mod q)
//不作模q处理，令 t = p·r*g+f*m
// d = fp * a (mod p) = fp * (p·r*g+f*m (mod q)) (mod p)
// = p·fp*r*g + fp*f*m (mod q)(mod p)
// =

//参数{N,p,q,df,dg,dr}选取不恰当的话会造成3种可能
//1.多项式环上找不到存在模p逆或者模q逆的小多项式f，使得密钥对生成困难
//2.系统解密限制失败和间距失败的概率很高，以致系统不可用
//3.安全性大幅降低，稍加运算即可破解密文

//推荐使用的标准参数
//安全等级	N	q	p	df	dg	dr
//低		107	64	3	15	12	5
//中		167	128	3	61	20	18
//高		263	128	3	50	24	16
//极高		503	256	3	216	72	55