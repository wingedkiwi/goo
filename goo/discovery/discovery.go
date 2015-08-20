// (c) Chi Vinh Le <cvl@chinet.info> – 13.06.2015

package discovery

type Discovery interface {
	GetRepository(path string) (string, error)
}
