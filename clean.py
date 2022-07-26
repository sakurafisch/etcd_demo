import os
from sys import platform

def main() -> None:
    dir_to_delete: str = "default.etcd"
    if platform == "linux" or platform == "linux2" or platform == "darwin":
        os.system(f"rm -rf {dir_to_delete}")
        return
    if platform == "win32":
        os.system(f"echo Y | rmdir /s {dir_to_delete}")
        return
    os.system(f"rm -rf {dir_to_delete}")

if __name__ == '__main__':
    main()
