#include "bitmap.h"
void MakeItAs(ul* data,ul p,char v){
	if (v){
		SETBIT(&data[p/(sizeof(ul)<<X8)],p%(sizeof(ul)<<X8));
	}else{
		CLRBIT(&data[p/(sizeof(ul)<<X8)],p%(sizeof(ul)<<X8));
	}
}
void MakeItAsArea(ul* data,ul start,ul end,char v){
	for (ul i = start;i < end ;i++){
		MakeItAs(data,i,v);
	}
}
void MakeItAsAreal(ul* data,ul start,ul end,char* v){
	for (ul i = start;i < end ;i++){
		MakeItAs(data,i,v[i]);
	}
}
char AcquirePosition(ul* data,ul p){
	return BITVAL(data[p/(sizeof(ul)<<X8)],p)?1:0;
}
ul Size(void *data){
	return *((ul *)(data-sizeof(ul)));
}
ul Length(void *data){
	return *((ul *)(data-(sizeof(ul)<<X2)));
}
struct BitMap * BitMapInit(ul size){
	struct BitMap *map = NULL;
	ul count = size%sizeof(ul)==0?size/(sizeof(ul)<<X8):size/(sizeof(ul)<<X8)+1;
	if ((map = malloc(sizeof(struct BitMap)+count*sizeof(ul))) == NULL){
		return map;
	}
	map->count = count;
	map->length = size;
	memset(map->data,0,count*sizeof(ul));
	return map;
}

