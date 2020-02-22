package core

type subLoopGroup struct {
	loops []Loop
	index int
}

func (g *subLoopGroup) registe(lp Loop) {
	g.loops = append(g.loops, lp)
}

//for put new accept connection in subLoop load balance
func (g *subLoopGroup) next() Loop {
	g.index++
	size := len(g.loops)
	if g.index >= size {
		g.index -= size
	}
	return g.loops[g.index]
}

func (g *subLoopGroup) iterate(reverse bool,f func(Loop) bool) {
	start,end,delta:= 0,len(g.loops)-1,1
	if reverse {
		start,end,delta = end,start,-delta
	}
	for i:=start;i<=end;i +=delta{
		if !f(g.loops[i]) {
			break
		}
	}
}
