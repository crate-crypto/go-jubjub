package curve

import fq "github.com/decentralisedkev/go-jubjub/internal/Fq"

type CompletedPoint struct {
	u, v, z, t fq.FieldQ
}

func (cp *CompletedPoint) AddExtended(q, r ExtendedPoint) *CompletedPoint {

	var a, b, c, d, t, d2 fq.FieldQ
	d2.SetD2()

	a.Sub(q.v, q.u)
	t.Sub(r.v, r.u)
	a.Mul(a, t)
	b.Add(q.u, q.v)
	t.Add(r.u, r.v)
	b.Mul(b, t)
	c.Mul(q.t1, r.t1)
	c.Mul(c, q.t2)
	c.Mul(c, r.t2)
	c.Mul(c, d2)
	d.Mul(q.z, r.z)
	d.Double()
	cp.u.Sub(b, a)
	cp.t.Sub(d, c)
	cp.z.Add(d, c)
	cp.v.Add(b, a)
	return cp
}
