// (c) Chi Vinh Le <cvl@chinet.info> – 13.06.2015

package gopkg

type Discovery interface {
	GetRepository(path string) (string, error)
}
