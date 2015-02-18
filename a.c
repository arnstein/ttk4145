#include <stdio.h>
#include <pthread.h>

int count = 0;
pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;

void* add( void* arg){
	printf("adder da\n");
	int i;
	for(i = 0; i < 1000 * 1000; i++){
		pthread_mutex_lock(&mutex);
		count++;
		pthread_mutex_unlock(&mutex);
	}
}
void* sub( void* arg){
	printf("sub da\n");
	int i;
	for(i = 0; i < 1000 * 1000; i++){
		pthread_mutex_lock(&mutex);
		count--;
		pthread_mutex_unlock(&mutex);
	}
}

int main(){

	pthread_t t1, t2;

	int err = 0;

	err |= pthread_mutex_init(&mutex, NULL);

	err |= pthread_create(&t1, NULL, &add, NULL);
	err |= pthread_create(&t2, NULL, &sub, NULL);

	if(err){
		printf("some init went wrong\n");
		return -1;
	}

	pthread_join(t1, NULL);
	pthread_join(t2, NULL);

	printf("count = %d\n", count);

	return 0;
}
