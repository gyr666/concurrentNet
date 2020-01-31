#include "bitmap.h"
void MakeItAs(ul* data,ul p,char v){
	if (v){
		SETBIT(&data[p/(sizeof(ul)<<X8)],p%(sizeof(ul)<<8));
	}else{
		CLRBIT(&data[p/(sizeof(ul)<<8)],p%(sizeof(ul)<<8));
	}
}
char AcquirePosition(ul* data,ul p){
	return BITVAL(data[p/(sizeof(ul)<<8)],p);
}
struct BitMap * BitMapInit(ul size){
	struct BitMap *map = NULL;
	ul count = size%sizeof(ul)==0?size/(sizeof(ul)<<X8):size/(sizeof(ul)<<X8)+1;
	if ((map = malloc(sizeof(struct BitMap)+count*sizeof(ul))) == NULL){
		return map;
	}
	map->count = count;
	map->length = size;
	map->MakeItAs = MakeItAs;
	memset(map->data,0,count*sizeof(ul));
	return map;
}

