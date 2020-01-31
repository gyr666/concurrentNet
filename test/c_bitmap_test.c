#include <stdio.h>
#include "../util/bitmap.h"

int main(char l,char ** ls){
	struct BitMap *b = BitMapInit(20);
	if (!b){
		printf("seg error\n");
	}
	printf("-----------%d----------\n",AcquirePosition(b->data,6));
	MakeItAs(b->data,6,(char)1);
	printf("-----------%d----------\n",AcquirePosition(b->data,6));
	MakeItAs(b->data,6,(char)0);
	printf("-----------%d----------\n",AcquirePosition(b->data,6));
return 0;
}
