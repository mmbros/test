package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// ProgramItem represents the basic element of a program.
type ProgramItem struct {
	n int
	k string
}

func newProgramItem(s string) *ProgramItem {
	var p ProgramItem
	j := strings.LastIndexAny(s, "0123456789")
	if j < 0 {
		p.n = 1
		p.k = s
	} else {
		p.n, _ = strconv.Atoi(s[:j+1])
		p.k = s[j+1:]
	}
	return &p
}

func (pi *ProgramItem) String() string {
	return fmt.Sprintf("%d%s", pi.n, pi.k)
}

// ProgramRow represents a single row of items.
type ProgramRow []*ProgramItem

// Program is a program.
type Program []ProgramRow

func (pr ProgramRow) String() string {
	var a []string
	for _, pi := range pr {
		a = append(a, pi.String())
	}
	return strings.Join(a, " ")
}

// Len return the length of the row.
func (pr ProgramRow) Len() int {
	var n int
	for _, v := range pr {
		n += v.n
	}
	return n
}

// Add adds a new row in a program
func (p *Program) Add(row string) error {
	r := ProgramRow{}

	for _, v := range strings.Fields(row) {
		pi := newProgramItem(v)
		if pi.n == 0 {
			return fmt.Errorf("Invalid program item %q at row #%d", v, len(*p)+1)
		} else if pi.k == "" {
			return fmt.Errorf("Invalid program item %q at row #%d: missing color", v, len(*p)+1)

		}
		r = append(r, pi)
	}
	if len(r) == 0 {
		return fmt.Errorf("Invalid program at row #%d", len(*p)+1)
	}

	*p = append(*p, r)
	return nil
}

// Fprint writes to w a representation of the Program.
func (p Program) Fprint(w io.Writer) {
	for j, r := range p {
		fmt.Fprintf(w, "%d = %s\n", j+1, r.String())
	}
}

// Print prints a representation of the Program.
func (p Program) Print() {
	p.Fprint(os.Stdout)
}

// Size returns the dimensions (Dx, Dy) = (cols, rows) of the image
// generated by the program.
func (p Program) Size() (int, int) {
	var cols int
	for _, v := range p {
		c := v.Len()
		if c > cols {
			cols = c
		}
	}
	return cols, len(p)
}

// CheckColors is ...
func (p Program) CheckColors(mp *Palette) error {
	for rownum, r := range p {
		for _, i := range r {
			if !mp.HasKey(i.k) {
				return fmt.Errorf("Unknown  color %q at row #%d", i.k, rownum+1)
			}
		}
	}
	return nil
}