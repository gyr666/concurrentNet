#include <stdio.h>
#include "../util/bitmap.h"
#include "ctest.h"

int main(char l,char ** ls){
	struct BitMap *b = BitMapInit(20);
	if (!b){
		perror("seg error\n");
	}
	EQUAL(AcquirePosition(b->data,6),0,"init test");
	MakeItAs(b->data,6,(char)1);
	EQUAL(AcquirePosition(b->data,6),1,"set test");
	MakeItAs(b->data,6,(char)0);
	EQUAL(AcquirePosition(b->data,6),0,"clear test");
	MakeItAsArea(b->data,10,15,(char)1);
	EQUAL(AcquirePosition(b->data,12),1,"area test");
	EQUAL(b->count,1,"size test");
	EQUAL(b->length,20,"count test");
	EQUAL(Length(b->data),20,"size test");
	EQUAL(Size(b->data),1,"count test");
return 0;
}
