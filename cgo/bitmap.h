// NOTHING
// no border check
#define ULONG
#define X2 1
#define X4 2
#define X8 3
#define SETBIT(x,y) (*x)|=(1<<y)
#define CLRBIT(x,y) (*x)&=!(1<<y)
#define BITVAL(x,y) ((x)&=(1<<y))
#include <string.h>
#include <stdlib.h>
#ifdef ULONG
typedef unsigned long ul;
#elif
typedef unsigned int ul;
#endif
struct BitMap {
	ul length;
	ul count;
	ul data[];
};

struct BitMap * BitMapInit(ul size);
void MakeItAs(ul* data,ul p,char v);
char AcquirePosition(ul* data,ul p);
void MakeItAsArea(ul* data,ul start,ul end,char v);
void MakeItAsAreal(ul* data,ul start,ul end,char* v);
ul Size(void *data);
ul Length(void *data);
