package buffer

struct divide struct {
	first	ByteBuffer
	last	ByteBuffer
	s	uint64
}

type sByteBuffer struct {
	r	[]byte
	next	ByteBuffer
}

type standAllocator struct {
	divs	[]divide
	min	uint64
	psize	uint8
	max	uint64
	l	sync.Mutex
}

func (s *standAllocator) Init(){
	// create alloc index
	s.divs = make([]divide,s.psize)
	for i, := range s.divs {
		s.divs[i].s = 2 << i
		// at init process, don't need create threads.
		s.backAlloc(i)
	}
	s.min = 2;
	s.max = 2 << p.size-1;
}

func (s *standAllocator) Alloc(length uint64) ByteBuffer {
	if(length > s.max) {
		return nil
	}
	return s.doAlloc(length);
}

func (s *standAllocator) PoolSize() uint64 {
	return s.psize
}

func (s *standAllocator) relases(b ByteBuffer) {
	bb ,ok:= b.(*sByteBuffer)
	if !ok {
		panic("release error!")
	}
	bb.next = s.[bb.index].first
	s.first = bb
}

func (s *standAllocator) doAlloc(length uint64) ByteBuffer {
	index := postition(length)
	if s.first == s.last {
		go s.backAlloc(index)
	}
	s.l.Lock()
	v := s.first
	s.first = v.next
	s.l.UnLock();
	return s
}

//this process is very important
func (s *standAllocator) backAlloc(index uint8) ByteBuffer {
	s.l.Lock()
	defer s.l.Unlock()
	result := make([]sByteBuffeImpl,s.load)
	for i,_ := range result {
		result[i].next = result[i+1]
	}
	s.div[index].first = result
	s.div[index].last  = result + s.load
}


