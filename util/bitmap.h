// NOTHING
#define X2 1
#define X4 2
#define X8 3
#define SETBIT(x,y) (*x)|=(1<<y)
#define CLRBIT(x,y) (*x)&=!(1<<y)
#define BITVAL(x,y) ((x)|=(1<<y))
#include <string.h>
#include <stdlib.h>
typedef unsigned long ul;
struct BitMap {
	void (*MakeItAs)(ul position,char v);
	ul length;
	ul count;
	ul data[];
};
