package core

type subLoopGroup struct {
	loops []Loop
	index int
}

func (g *subLoopGroup) registe(lp Loop){
	g.loops=append(g.loops,lp)
}

//for put new accept connection in subLoop load balance
func (g *subLoopGroup) next()Loop{
	g.index ++
	size := len(g.loops)
	if g.index == size{
		g.index -= size
		//index 0 is mainLoop , so we should start with index 1
		g.index ++
	}
	return g.loops[g.index]
}

func (g *subLoopGroup) iterate(f func(Loop)bool){
	for _,lp:=range g.loops{
		if !f(lp){
			break
		}
	}
}
