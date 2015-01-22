import threading

i = 0
useCounter = threading.Lock()

def countUp():
    global i
    for j in range(0, 1000000):
        useCounter.acquire(True)
        i += 1
        useCounter.release()

def countDown():
    global i
    for k in range(0, 1000001):
        useCounter.acquire(True)
        i -= 1
        useCounter.release()

thread1 = threading.Thread(target = countUp, args = (),)
thread2 = threading.Thread(target = countDown, args = (),)
thread1.start()
thread2.start()

thread1.join()
thread2.join()

print i
