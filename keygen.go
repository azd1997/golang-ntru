package golang_ntru

//多项式环上的数学运算

//整数环Z，整数N>=2，R表示多项式截断环，R=Z[X]/(X^N-1)
//对于任意正整数q，Rq表示模q的多项式截断环，Rq=(Z/qZ)[X]/(X^N-1)
//q为素数时，Rq具有可逆性，可逆性指对于F及其模q逆Fq，有F·Fq=1 mod q

//NTRU算法需要3个整数参数（N,p,q）和4个多项式环（f,g,r,m）
//要求gcd(p,q)=1（即二者最大公约数为1），且p远大于q

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


/************************************/
//Name: Multi
//Inputs:
//Outputs:
//Description: 多项式环的乘法
/************************************/
func Multi(N, d, f uint64) []uint64 {

}



/************************************/
//Name:
//Inputs:
//Outputs:
//Description: 随机产生多项式
/************************************/
func generatePolynomial(N, d, f uint64) []uint64 {

}

