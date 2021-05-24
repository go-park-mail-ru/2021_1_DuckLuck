import os
import time
import datetime
import argparse
import subprocess
import constants
import workers


THREADS = []


def rebuild(targets):
    for target in targets:
        print(f'\33[92mRemoving {target} from local \033[0m')
        thread = workers.RebuildThread(target)
        THREADS.append(thread)
        thread.start()


def check_images():
    result = subprocess.run(['docker', 'images'], stdout=subprocess.PIPE)
    images = [image.split()[0] for image in result.stdout.decode("utf-8").split('\n')[1:-1]]
    for image in constants.IMAGES_FROM_DOCKERHUB:
        if image not in images:
            print(f'\033[91mImage {image} not found \033[0m')
            print(f'\33[92mPulling {image} from DockerHub \033[0m')
            thread = workers.PullThread(image)
            THREADS.append(thread)
            thread.start()
        else:
            print(f'\33[92mImage {image} found local \033[0m')
    for image in constants.LOCAL_UP_TARGETS_WITH_BUILD:
        if image not in images:
            print(f'\033[91mImage {image} not found \033[0m')
            print(f'\33[92mBuilding {image} from local \033[0m')
            thread = workers.BuildThread(image)
            THREADS.append(thread)
            thread.start()
        else:
            print(f'\33[92mImage {image} found local \033[0m')


def up_local():
    check_images()
    for thread in THREADS:
        thread.join()
    result = subprocess.run(['docker-compose', 'up', '-d'], stderr=subprocess.PIPE)
    if result.stderr:
        print(result.stderr.decode("utf-8"))
    else:
        print(f'\33[102mAll containers started \033[0m')


if __name__ == '__main__':
    begin = time.time()
    parser = argparse.ArgumentParser()
    parser.add_argument("--rebuild_targets", default=None, help="Targets for rebuild")
    parser.add_argument("--target", default=None, help="Target for script")
    args = parser.parse_args()
    target = args.target
    if target == 'rebuild':
        rebuild(args.rebuild_targets)
    elif target == 'up_local':
        up_local()
    else:
        os.system('echo skip target')
    end = time.time()
    print('\n'*5)
    print(f'Total time: {datetime.timedelta(seconds=(end - begin))}')