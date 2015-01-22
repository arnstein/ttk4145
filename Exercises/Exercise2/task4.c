#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t countMutex;

void* countUp()
{
    for (int j = 0; j < 1000000; j++) {
        pthread_mutex_lock(&countMutex);
        i += 1;
        pthread_mutex_unlock(&countMutex);
    }
    return NULL;
}

void* countDown()
{
    for (int j = 0; j < 1000001; j++) {
        pthread_mutex_lock(&countMutex);
        i -= 1;
        pthread_mutex_unlock(&countMutex);
    }
    return NULL;
}


int main() {
    pthread_t thread_1;
    pthread_t thread_2;
    pthread_create(&thread_1, NULL, countUp, NULL);
    pthread_create(&thread_2, NULL, countDown, NULL);
    pthread_join(thread_1, NULL);
    pthread_join(thread_2, NULL);
    printf("%d\n", i);
}
