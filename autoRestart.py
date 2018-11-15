import os,sys,time

if __name__== '__main__':
    pid = os.fork()
    if pid > 0:
        sys.exit(0)
    os.setsid()
    pid = os.fork()
    if pid > 0:
        sys.exit(0)

    while True:
        time.sleep(10 * 60)
        res = os.popen('ps -Af | grep bilibiliVideoDataCrawler | grep -v "grep" | wc -l')
        cnt = int(res.read())
        if cnt == 0:
            os.system("./bilibiliVideoDataCrawler")

