import threading
import thread

count = 0
lock = thread.allocate_lock()

def add():
	global count
	global lock
	global done1
	for i in range(1000 * 1000):
		with lock:
			count += 1

def sub():
	global count
	global lock
	global done2
	for i in range(1000 * 1000 + 1):
		with lock:
			count -= 1

t1 = threading.Thread(target = add, args = ())
t2 = threading.Thread(target = sub, args = ())

t1.start()
t2.start()

t1.join()
t2.join()

print "count = " + str(count)
