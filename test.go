package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/azd1997/golang-ntru/ntru_crypto"
	"github.com/azd1997/golang-ntru/ntru_utils/params"
	"log"
)

func main() {
	fmt.Println("测试初始化！")

	keypair, err := ntru_crypto.GenerateKey(rand.Reader, params.EES1171EP1)
	if err != nil {
		log.Fatal(err)
	}

	blen := keypair.Params.MaxMsgLenBytes
	plaintext := make([]byte, blen)
	if _, err = rand.Reader.Read(plaintext); err != nil {
		log.Fatal(err)
	}

	ciphertext, err := ntru_crypto.Encrypt(rand.Reader, &keypair.PublicKey, plaintext)
	if err != nil {
		log.Fatal(err)
	}

	plaintext2, err := ntru_crypto.Decrypt(keypair, ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	if bytes.Compare(plaintext, plaintext2) != 0 {
		log.Fatal("plaintext != plaintext2")
	}

	fmt.Println("测试成功：plaintext = plaintext2")






















	//var ntru keygen.NtruCipher
	//err := ntru.Init(11, 111, 1, 2, 3, 3)
	//if err != nil {
	//	fmt.Println("error happend!")
	//}
	//fmt.Println(ntru.N)

	//a := []int{2,1,3,5}
	//b := []int{1,2,1,3}
	//c := keygen.Multi(a,b,4)
	//fmt.Println(c)

	//a := []int{12,-14,-8,8,-1}
	//	//p := 3
	//	//r := keygen.Amodp(a,5,p)
	//	//fmt.Println(r)
	//	//b:= []int{1,4,-3,2,1}
	//	//r1 := keygen.BmodXn(b, 2)
	//	//fmt.Println(r1)

	//a := []int{-2,1,-3,0,2,0,3}
	//n := 7
	//p := 5
	//ng := 1
	//b := keygen.AB_1modp(a,n,p,ng)
	//fmt.Println(b)
	//
	//fmt.Println(keygen.Random_f(9,3))
	//fmt.Println(keygen.Random_f(9,3))
	//fmt.Println(keygen.Random_f(9,3))
	//
	//r := rand.New(rand.NewSource(99))
	//fmt.Println(r.Int())
	//fmt.Println(r.Int())
	//fmt.Println(r.Int())
	//fmt.Printf("%b\n", rand.Int())
	//
	//fmt.Println(int(math.Floor(float64(n*r.Intn(100) / 100))))
	//fmt.Println(int(math.Floor(float64(n*r.Intn(100) / 100))))
	//fmt.Println(int(math.Floor(float64(n*r.Intn(100) / 100))))

	//a:=[]int{14,13,9,15,-14,15,16}
	//n:=7
	//q:=32
	//fmt.Println(keygen.AmodP2(n,q,a))

	//fmt.Println(NTRU_2001.Random_gr(7,2))
	//fmt.Println(NTRU_2001.RandomPoly(9,3,5))

	//fmt.Println(NTRU_2001.DegOfPoly(a))

	//f := []int{1,1,1,0,-1,0,1,0,0,1,-1}
	////g := []int{-1,0,1,1,0,1,0,0,-1,0,-1}
	//fmt.Println(duplicated.Invert(f,11,3))

}
