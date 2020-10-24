package cipher

import (
	"errors"
	"log"
)

type ModularBezoutIdentity struct {
	MMInverse          int64
	BezoutIdentityModB int64
	BezoutIdentity     int64
}

func (c *Cipher) IsCoprime() bool {
	return GCD(c.Key, c.Modulus) == 1
}

func GCD(A int64, B int64) int64 {
	if A == 0 {
		return B
	} else if B == 0 {
		return A
	}
	return GCD(B, A%B)
}

func ExplainExtendedGCD(A int64, B int64) (ModularBezoutIdentity, error) {
	gcd, prevX, prevY := ExtendedGCD(A, B)

	// modulus
	AModB := A % B
	BModB := B % B
	prevXModB := Modulus(prevX, B)
	prevYModB := Modulus(prevY, B)

	// regular BezoutIdentity
	Ax := A * prevX
	By := B * prevY
	BezoutIdentity := Ax + By

	// BezoutIdentity with modulus
	AxModB := AModB * prevXModB
	ByModB := BModB * prevYModB
	BezoutIdentityModB := AxModB + ByModB

	// modularMultiplicativeInverse
	modularMultiplicativeInverse := ModularBezoutIdentity{
		MMInverse:          0,
		BezoutIdentityModB: BezoutIdentityModB,
		BezoutIdentity:     BezoutIdentity,
	}
	if gcd != 1 {
		return modularMultiplicativeInverse, errors.New("No coprime numbers:" + string(A) + "" + string(B))
	}
	if AxModB != 0 {
		modularMultiplicativeInverse.MMInverse = prevXModB
	} else if ByModB != 0 {
		modularMultiplicativeInverse.MMInverse = prevYModB
	}

	// explain step by step mathematical operations
	log.Println("-Modular Multiplicative Inverse-")
	log.Printf("ax + by =  %d(%d) + %d(%d) = %d \n\n", A, prevX, B, prevY, BezoutIdentity)
	log.Printf("ax + by mod B =  %d(%d) + %d(%d) mod %d = %d\n\n", AModB, prevXModB, BModB, prevYModB, B, BezoutIdentityModB)
	log.Printf("1/%d = %d\n", A, modularMultiplicativeInverse.MMInverse)
	log.Println("-End Modular Multiplicative Inverse-")

	return modularMultiplicativeInverse, nil
}

func XGCDModB(A int64, B int64) (int64, int64, int64) {
	gcd, prevX, prevY := ExtendedGCD(A, B)
	return gcd, Modulus(prevX, B), Modulus(prevY, B)
}

/**
* More information:
*
* http://anh.cs.luc.edu/331/notes/xgcd.pdf
* https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
* https://www.khanacademy.org/computing/computer-science/cryptography
**/
func ExtendedGCD(A int64, B int64) (int64, int64, int64) {
	var (
		prevX    int64 = 1
		prevY    int64 = 0
		currentX int64 = 0
		currentY int64 = 1
		gcd      int64 = 0
	)
	for ; B != 0 && A != 0; {
		q, r := A/B, A%B
		s, t := prevX-q*currentX, prevY-q*currentY
		prevX, prevY = currentX, currentY
		currentX, currentY = s, t
		A, B = B, r
	}

	if A == 0 {
		gcd = B
	} else {
		gcd = A
	}

	return gcd, prevX, prevY
}
