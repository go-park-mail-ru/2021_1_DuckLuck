import threading
import subprocess


class BuildThread(threading.Thread):
    def __init__(self, target: str):
        threading.Thread.__init__(self)
        self.target = target

    def run(self):
        print(f'\33[33mStarting build {self.target}\033[0m')
        build_image(self.target)
        print(f'\33[32mExiting build {self.target}\033[0m')


def build_image(image):
    result = subprocess.run(['docker', 'build', '-t', f'{image}:local', '--target', image[image.rindex('/') + 1:], '.'], stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    if result.stderr:
        print(result.stderr.decode("utf-8"))


class PullThread(threading.Thread):
    def __init__(self, target: str):
        threading.Thread.__init__(self)
        self.target = target

    def run(self):
        print(f'\33[33Starting pull {self.target}\033[0m')
        pull_image(self.target)
        print(f'\33[32mExiting pull {self.target}\033[0m')


def pull_image(image):
    result = subprocess.run(['docker', 'pull', image], stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    if result.stderr:
        print(result.stderr.decode("utf-8"))


class RebuildThread(threading.Thread):
    def __init__(self, target: str):
        threading.Thread.__init__(self)
        self.target = target

    def run(self):
        print(f'\33[33Starting rebuild {self.target}\033[0m')
        pull_image(self.target)
        print(f'\33[32mExiting rebuild {self.target}\033[0m')


def rebuild_image(target):
    result = subprocess.run(['docker-compose', 'rm', '-vfs', target], stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    if result.stderr:
        print(result.stderr.decode("utf-8"))
    print(f'\33[92mBuilding {target} from local \033[0m')
    result = subprocess.run(['docker', 'build', '-t', f'{target}:local', '--target', target], stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    if result.stderr:
        print(result.stderr.decode("utf-8"))