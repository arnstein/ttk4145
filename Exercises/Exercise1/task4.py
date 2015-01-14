from threading import Thread

i = 0

def countUp():
    global i
    for j in range(0, 1000000):
        i += 1

def countDown():
    global i
    for k in range(0, 1000000):
        i -= 1

thread1 = Thread(target = countUp, args = (),)
thread2 = Thread(target = countDown, args = (),)
thread1.start()
thread2.start()

thread1.join()
thread2.join()

print i
