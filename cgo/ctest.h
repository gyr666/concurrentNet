#include<stdio.h>
void EQUAL(int val0, int val1,char *v){
	if(val0!=val1){
		dprintf(2,"%s\t:\033[31;1m [ERROR ] \033[0m %d not equal %d\n",v,val0,val1);
		return;
	}
	dprintf(1,"%s\t:\033[32;1m [PASSED] \033[0m\n",v);
}
	
