package utils

import (
	"bufio"
	"bytes"
)

type ColumnScanner interface {
	Next() bool
	Line() []string
}

type columnScanner struct {
	scanners []*bufio.Scanner
	next     []string
}

func (c *columnScanner) Next() bool {
	var hasNext bool
	for i, s := range c.scanners {
		if s == nil {
			continue
		}
		if !s.Scan() {
			c.scanners[i] = nil
			continue
		}
		hasNext = true
	}
	if !hasNext {
		c.scanners = nil
		return false
	}
	return true
}

func (c columnScanner) Line() []string {
	if len(c.scanners) == 0 {
		return nil
	}
	ss := make([]string, len(c.scanners))
	for i, scn := range c.scanners {
		if scn == nil {
			continue
		}
		ss[i] = scn.Text()
	}
	return ss
}

func NewColumnScanner(values []string) ColumnScanner {
	cs := &columnScanner{}
	for _, v := range values {
		cs.scanners = append(cs.scanners, bufio.NewScanner(bytes.NewBufferString(v)))
	}
	return cs
}
