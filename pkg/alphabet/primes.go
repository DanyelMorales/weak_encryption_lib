package alphabet

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"
)

//LoadPrimes Prime loader
func LoadPrimes(callback chan int64) {

	r := csv.NewReader(strings.NewReader(PrimeNumbers))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		for value := range record {
			prime, errorC := strconv.Atoi(record[value])
			if errorC != nil {
				log.Fatal("Prime reader: Not a prime number", errorC)
			}
			callback <- int64(prime)
		}

	}
	callback <- -1
}

func factorizeWithPrime(n int64, prime int64, counter int64) (int64, int64) {
	a, b := n/prime, n%prime

	if b > 0 {
		return n, counter
	}
	counter++
	return factorizeWithPrime(a, prime, counter)
}

// FactorizeNumber  factorize a number
func FactorizeNumber(n int64) map[int64]int64 {
	cb := make(chan int64)
	primes := make(map[int64]int64, 0)

	go LoadPrimes(cb)

	for {
		prime := <-cb
		if prime == -1 {
			break
		}

		n, times := factorizeWithPrime(n, prime, 0)

		if times > 0 {
			primes[prime] = times
		}

		if n == 1 {
			close(cb)
			break
		}
	}
	return primes
}
