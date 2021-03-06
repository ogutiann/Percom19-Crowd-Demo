package zcrypto

import (
	"crypto/rand"
	"log"
	"math/big"
)

const (
	N = "78026704611133664070808216822389351558128055220211968857081009471051751212773"
	G = "78026704611133664070808216822389351558128055220211968857081009471051751212774"
	Lambda = "78026704611133664070808216822389351557567911749966982131663303912667682518368"
	Mu = "27507391580316718640261341452623916400024156305339914611702685519047060407763"
)

var PubKey *PublicKey
var PriKey *PrivateKey

type PublicKey struct {
	N, G *big.Int
	n2 *big.Int
}

type Cypher struct {
	C *big.Int
}

type PrivateKey struct {
	PublicKey
	Lambda *big.Int
	Mu *big.Int
}

func (this *PublicKey) GetNSquare() *big.Int {
	if this.n2 !=nil {
		return this.n2
	}
	this.n2 = new(big.Int).Mul(this.N,this.N)
	return this.n2
}

func LCM(x, y *big.Int) *big.Int {
	return new(big.Int).Mul(new(big.Int).Div(x, new(big.Int).GCD(nil, nil, x, y)), y)
}

func minusOne(x *big.Int) *big.Int {
	return new(big.Int).Add(x, big.NewInt(-1))
}

func computeMu(g, lambda, n *big.Int) *big.Int {
	n2 := new(big.Int).Mul(n, n)
	u := new(big.Int).Exp(g, lambda, n2)
	return new(big.Int).ModInverse(L(u, n), n)
}

func L(u, n *big.Int) *big.Int {
	t := new(big.Int).Add(u, big.NewInt(-1))
	return new(big.Int).Div(t, n)
}

func NewPallier(bits int) (*PublicKey,*PrivateKey,error) {
	random:= rand.Reader
	var p,q *big.Int
	var errChan = make(chan error,1)
	go func(){
		var err error
		p, err = rand.Prime(random,bits)
		errChan <-err
	}()

	q, err:= rand.Prime(random,bits)
	if err!=nil {
		return nil,nil,err
	}

	if err = <-errChan; err!=nil {
		return nil,nil,err
	}

	log.Println("length p:", len(p.Bytes()))
	log.Println("length q:", len(q.Bytes()))

	n := new(big.Int).Mul(p, q)
	lambda := new(big.Int).Mul(minusOne(p), minusOne(q))
	g := new(big.Int).Add(n, big.NewInt(1))
	mu := new(big.Int).ModInverse(lambda, n)

	pub:= &PublicKey{
		N:n,
		G:g,
	}

	pri:= &PrivateKey{
		PublicKey: PublicKey {
			N: new(big.Int).Set(n),
			G: new(big.Int).Set(g),
		},
		Lambda:lambda,
		Mu:mu,
	}

	return pub,pri,nil
}

func (this *PublicKey)Encrypt(m *big.Int) (*Cypher,error) {
	r,err:= RandomNumberInGroup(this.N)
	if err != nil {
		return nil, err
	}
	nSquare := this.GetNSquare()

	var message *big.Int
	if m.Cmp(big.NewInt(0))==-1 {
		message = new(big.Int).ModInverse(new(big.Int).Neg(m),this.N)
	} else {
		message = m
	}

	gm := new(big.Int).Exp(this.G, message, nSquare)
	rn := new(big.Int).Exp(r, this.N, nSquare)
	return &Cypher{
		C:new(big.Int).Mod(new(big.Int).Mul(rn, gm), nSquare),
	}, nil
}

func (this *PrivateKey) Decrypt(cypher *Cypher) (*big.Int) {
	tmp := new(big.Int).Exp(cypher.C, this.Lambda, this.GetNSquare())
	msg := new(big.Int).Mod(new(big.Int).Mul(L(tmp, this.N), this.Mu), this.N)
	return msg
}

func init() {
	_n , _ := new(big.Int).SetString(N,10)
	_g , _ := new(big.Int).SetString(G,10)
	_l,_:= new(big.Int).SetString(Lambda,10)
	_m,_:= new(big.Int).SetString(Mu,10)
	PubKey = &PublicKey{
		N: _n,
		G: _g,
	}
	PriKey = &PrivateKey {
		PublicKey:PublicKey{
			N:_n,
			G:_g,
		},
		Lambda:_l,
		Mu: _m,
	}
}