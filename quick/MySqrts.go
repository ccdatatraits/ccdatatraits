
package main

import (
	"fmt"
	"math"
	"math/big"
	/*"golang.org/x/tour/pic"*/
)

const minuscule float64 = 2.220446049250313e-16

func Sqrt(x float64) float64 {
	fmt.Println("Sqrt")
	z := 1.0
	p := z
	initialCheck := true
	for initialCheck || math.Abs(p - z) > minuscule {
		fmt.Println(math.Abs(p - z))
		p = z
		z -= (z*z -x) / (2*z)
		if initialCheck {initialCheck = false}
	}
	fmt.Println()
	return z
}

/*func Pic(dx, dy int) [][]uint8 {
	array := make([][]uint8, dx)
	for i := range array {
		array[i] = make([]uint8, dy)
		for j := range array[i] {
			array[i][j] = uint8(1)
		}
	}
	return array
}*/

func main() {
	/*pic.Show(Pic)*/
	number := float64(2)

	mathSqrt := math.Sqrt(number)
	mySqrt := Sqrt(number)
	myBigSqrt := big_Sqrt(number)

	fmt.Println("math.sqrt: ", mathSqrt)
	fmt.Println("Sqrt: ", mySqrt, " with diff:", math.Abs(mathSqrt - mySqrt))
	fmt.Println("Big Sqrt: ", myBigSqrt, " with diff:", math.Abs(mathSqrt - myBigSqrt))


	// We'll do computations with 200 bits of precision in the mantissa.
	const prec = 200

	// Compute the square root of 2 using Newton's Method. We start with
	// an initial estimate for sqrt(2), and then iterate:
	//     x_{n+1} = 1/2 * ( x_n + (2.0 / x_n) )

	// Since Newton's Method doubles the number of correct digits at each
	// iteration, we need at least log_2(prec) steps.
	steps := int(math.Log2(prec))

	// Initialize values we need for the computation.
	two := new(big.Float).SetPrec(prec).SetInt64(2)
	half := new(big.Float).SetPrec(prec).SetFloat64(0.5)

	// Use 1 as the initial estimate.
	x := new(big.Float).SetPrec(prec).SetInt64(1)

	// We use t as a temporary variable. There's no need to set its precision
	// since big.Float values with unset (== 0) precision automatically assume
	// the largest precision of the arguments when used as the result (receiver)
	// of a big.Float operation.
	t := new(big.Float)

	// Iterate.
	for i := 0; i <= steps; i++ {
		t.Quo(two, x)  // t = 2.0 / x_n
		t.Add(x, t)    // t = x_n + (2.0 / x_n)
		x.Mul(half, t) // x_{n+1} = 0.5 * t
	}

	// We can use the usual fmt.Printf verbs since big.Float implements fmt.Formatter
	fmt.Printf("sqrt(2) = %.50f\n", x)

	// Print the error between 2 and x*x.
	t.Mul(x, x) // t = x*x
	fmt.Printf("error = %e\n", t.Sub(two, t))
}

func big_Sqrt(input float64) float64 {
	fmt.Println("Big Sqrt")
	big_input := new(big.Float).SetFloat64(input)
	const prec = 200
	z := new(big.Float).SetPrec(prec).SetInt64(1)
	p := new(big.Float).Set(z)
	tmp := new(big.Float)
	big_minuscule := new(big.Float).SetFloat64(minuscule)
	two := new(big.Float).SetPrec(prec).SetInt64(2)
	t := new(big.Float)
	initialCheck := true
	for initialCheck || tmp.Abs(tmp.Sub(p, z)).Cmp(big_minuscule) == 1 {
		fmt.Println(tmp.Abs(tmp.Sub(p, z)))
		p.Set(z)
		t.Mul(two, z)
		tmp.Mul(z, z)
		tmp.Sub(tmp, big_input)
		tmp.Quo(tmp, t)
		z.Sub(z, tmp)
		if initialCheck {initialCheck = false}
	}
	fmt.Println()
	_1, _ := z.Float64()
	return _1
}

//sqrt(2) = 1.41421356237309504880168872420969807856967187537695
