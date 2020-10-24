package cipher

import (
	"errors"
	"github.com/danyelmorales/weak_encryption_lib/pkg/alphabet"
	"github.com/danyelmorales/weak_encryption_lib/pkg/symbol"
	"log"
	"math"
)

// Cipher used
type Cipher struct {
	Key      int64
	Modulus  int64
	Alphabet []symbol.Symbol
}

// Encrypt from symbol
func (c *Cipher) Encrypt(input []symbol.Symbol) []symbol.Symbol {
	log.Printf("··· Encrypting ...")
	output := make([]symbol.Symbol, len(input))
	for _, v := range input {
		var encrypted = v.Space()
		if v == 0 {
			continue
		}
		if !v.IsSpace() {
			a := int64(v) * c.Key
			b := a % (c.Modulus)
			encrypted = symbol.Symbol(b)
		}
		log.Printf("Symbol %d now is %d", v, encrypted)
		output = append(output, encrypted)
	}
	return output
}

func (c *Cipher) strToSym(input string) ([]symbol.Symbol, error) {
	log.Printf("··· getting byte equivalent of input ...")
	output := make([]symbol.Symbol, len(input))
	for _, v := range input {
		curentSymbol := symbol.Symbol(v)
		var mySymb = curentSymbol.Space()
		if !curentSymbol.IsSpace() {
			assocValue, err := curentSymbol.AssocValue(c.Alphabet)
			if err != nil {
				return nil, errors.New("the Symbol doesn't belong to the alphabet")
			}
			mySymb = assocValue
		}

		log.Printf("Symbol %d is now %d", curentSymbol, mySymb)
		output = append(output, mySymb)
	}
	return output, nil
}

func (c *Cipher) symToStr(input []symbol.Symbol) (string, error) {
	log.Printf("··· getting string equivalent of input ...")
	output := ""
	for _, v := range input {
		if v == 0 {
			continue
		}
		if v.IsSpace() {
			output += string(v.Space())
		} else {
			value, err := v.OriginalValue(c.Alphabet)
			if err != nil {
				log.Fatal(err)
			}
			output += string(value)
		}
	}
	return output, nil
}

// EncryptFromStr encrypt from string
func (c *Cipher) EncryptFromStr(input string) string {
	log.Printf("Input: %s", input)
	v, err := c.strToSym(input)
	if err != nil {
		log.Fatal(err)
	}
	encrypted := c.Encrypt(v)
	newValue, errOfConversion := c.symToStr(encrypted)
	if errOfConversion != nil {
		log.Fatal(errOfConversion)
	}
	log.Printf("Encrypted: %s as %s", input, newValue)
	return newValue
}

// DecryptFromStr decrypt
func (c *Cipher) DecryptFromStr(input string) string {
	v, err := c.strToSym(input)
	if err != nil {
		log.Fatal(err)
	}
	decrypted := c.Decrypt(v)
	//decrypted := c.DecryptBruteForce(v)
	newValue, errOfConversion := c.symToStr(decrypted)
	if errOfConversion != nil {
		log.Fatal(errOfConversion)
	}
	log.Printf("Decrypted: %s is now %s", input, newValue)
	return newValue
}

// DecryptBruteForce decrypt brute force
func (c *Cipher) DecryptBruteForce(input []symbol.Symbol) []symbol.Symbol {
	log.Printf("··· Decrypting ...")
	output := make([]symbol.Symbol, len(input))
	for _, v := range input {
		var decrypted = v.Space()
		if v == 0 {
			continue
		}
		if !v.IsSpace() {
		bruteforce:
			for {
				a := int64(v) % c.Key // determine if there is no exact division, do not pay attention
				if a != 0 {
					v += symbol.Symbol(c.Modulus)
					log.Printf("attemping reverse modulus with: %d\n", v)
				} else {
					decrypted = v / symbol.Symbol(c.Key)
					log.Printf("Decrypted: %d\n", decrypted)
					break bruteforce
				}
			}
		}
		output = append(output, decrypted)
	}
	return output
}

//MMInverse Euclidean
func (c *Cipher) Decrypt(input []symbol.Symbol) []symbol.Symbol {
	multiplicativeInverse, err := ExplainExtendedGCD(c.Key, c.Modulus)
	if err != nil {
		panic(err)
	}
	output := make([]symbol.Symbol, len(input))
	for _, v := range input {
		var decrypted = v.Space()
		if v == 0 {
			continue
		}
		if !v.IsSpace() {
			decrypted = v * symbol.Symbol(math.Abs(float64(multiplicativeInverse.MMInverse)))
			decrypted = decrypted % symbol.Symbol(c.Modulus)
			log.Printf("Decrypted: %d\n", decrypted)
		}
		output = append(output, decrypted)
	}
	return output
}

// GoodBadKeys using the current alphabet we're able to determine good keys and bad keys
func (c *Cipher) GoodBadKeys() ([]symbol.Symbol, []symbol.Symbol) {
	badKeysPattern := alphabet.FactorizeNumber(c.Modulus)
	badKeys := make([]symbol.Symbol, 0)
	goodKeys := make([]symbol.Symbol, 0)
	for i := 1; i <= len(c.Alphabet); i++ {

		isBadKey := false
		for prime := range badKeysPattern {
			if (int64(i) % prime) == 0 {
				isBadKey = true
				break
			}
		}
		if isBadKey {
			badKeys = append(badKeys, symbol.Symbol(i))
		} else {
			goodKeys = append(goodKeys, symbol.Symbol(i))
		}
	}
	return badKeys, goodKeys
}

// GetGoodKey  return a good key for the current  cipher suite
func (c *Cipher) GetGoodKey() symbol.Symbol {
	_, good := c.GoodBadKeys()
	return good[len(good)-1]
}

// IsGoodKey true if the key is a good key
func (c *Cipher) IsGoodKey(key symbol.Symbol) bool {
	bad, _ := c.GoodBadKeys()
	for _, bk := range bad {
		if key.Equals(bk) {
			return false
		}
	}
	return true
}

//New creates a new cipher from an alphabet, it resolves the key and the modulus to be used.
func New(alphabet []symbol.Symbol) Cipher {
	cipher := Cipher{
		Key:      0,
		Modulus:  int64(len(alphabet)),
		Alphabet: alphabet,
	}
	cipher.Key = int64(cipher.GetGoodKey())
	return cipher
}

//NewWithKey creates a new cipher from an alphabet, but is required to define a key.
func NewWithKey(alphabet []symbol.Symbol, key symbol.Symbol) Cipher {
	cipher := Cipher{
		Key:      int64(key),
		Modulus:  int64(len(alphabet)),
		Alphabet: alphabet,
	}
	if !cipher.IsGoodKey(key) {
		log.Fatal("insecure key detected!!")
	}
	return cipher
}
