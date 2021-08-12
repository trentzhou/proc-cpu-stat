#!/usr/bin/env python3
import os
import sys
import time
import subprocess

def cat(filename):
    content = open(filename).read()
    return str(content)

def parse_stat(filename):
    stat = cat(filename).strip()
    items = stat.split(" ")
    return {
        "utime": items[13],
        "stime": items[14],
        "cutime": items[15],
        "cstime": items[16]
    }

def get_cmd_output(cmd):
    return subprocess.check_output(cmd, shell=True)

class Thread(object):
    def __init__(self, pid, tid):
        super().__init__()
        self.pid = pid
        self.tid = tid
        self.comm = ""
        self.stat = None
        self.update()

    def update(self):
        task_folder = f"/proc/{self.pid}/task/{self.tid}"
        self.comm = cat(f"{task_folder}/comm").strip()
        self.stat = parse_stat(f"{task_folder}/stat")
        
    @property
    def utime(self):
        return int(self.stat['utime'])

    @property
    def stime(self):
        return int(self.stat['stime'])

    @property
    def id(self):
        return self.tid

    @property
    def name(self):
        return self.comm

    def __repr__(self):
        return f"Thread {self.tid} comm {self.comm} utime {self.stat['utime']} stime {self.stat['stime']} cutime {self.stat['cutime']} cstime {self.stat['cstime']}"

class Process(object):
    def __init__(self, pid):
        super().__init__()
        self.pid = pid
        self.cmdline = ""
        self.comm = ""
        self.threads = []
        self.stat = None
        self.update()

    def update(self):
        proc_folder = f"/proc/{self.pid}"
        self.comm = cat(f"{proc_folder}/comm").strip()
        self.cmdline = cat(f"{proc_folder}/cmdline")
        threads = os.listdir(f"{proc_folder}/task")
        self.threads = [Thread(self.pid, x) for x in threads if x.isdigit()]
        self.stat = parse_stat(f"{proc_folder}/stat")

    def cat(self, filename):
        content = open(filename).read()
        return str(content)

    @property
    def utime(self):
        return int(self.stat['utime'])

    @property
    def stime(self):
        return int(self.stat['stime'])

    @property
    def id(self):
        return self.pid

    @property
    def name(self):
        return self.comm

    def threads(self):
        pass

    def __repr__(self):
        lines = []
        lines.append(f"Pid {self.pid} comm {self.comm} utime {self.stat['utime']} stime {self.stat['stime']} cutime {self.stat['cutime']} cstime {self.stat['cstime']}")
        for t in self.threads:
            lines.append(str(t))
        return "\n".join(lines)


class ProcStat(object):
    def __init__(self):
        self.last_process = None
        self.last_stat = {}
        self.last_time = 0
        self.hz = int(subprocess.check_output("getconf CLK_TCK", shell=True).strip())
        self.report_lines = []
        self.cursor_up = 0

    def report_thread(self, thread_old, thread_new, time_delta):
        rate = (thread_new.utime + thread_new.stime - thread_old.utime - thread_old.stime) * 100 / self.hz / time_delta
        s = "{:.8s} {:18s} {:6.2f} {:s}".format(
            thread_new.id,
            thread_new.name,
            rate,
            "#" * int(rate + 0.5)
        )
        self.report_lines.append(s)

    def show_report(self):
        print("\033[2K\r\033[F" * self.cursor_up)
        for x in self.report_lines:
            print("\033[2K\r" + x)
        self.cursor_up = 1 + len(self.report_lines)

    def update(self, process: Process):
        now = time.time()
        self.report_lines.clear()

        if len(self.last_stat) > 0:
            time_delta = now - self.last_time
            self.report_thread(self.last_process, process, time_delta)
            self.report_lines.append("-------------------------------------")

            for t in process.threads:
                thread_old = self.last_stat.get(t.tid, None)
                if thread_old:
                    self.report_thread(thread_old, t, time_delta)
        self.last_time = now
        self.last_stat.clear()
        self.last_process = process
        for t in process.threads:
            self.last_stat[t.tid] = t
        self.show_report()


if __name__ == '__main__':
    pid = sys.argv[1]
    p = ProcStat()
    while True:
        try:
            p.update(Process(pid))
            time.sleep(1)
        except KeyboardInterrupt:
            break
