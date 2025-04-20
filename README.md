# Container from scratch

Naive and dummy implementation of a linux container runtime from scratch:
- File system namespacing
- Hostname isolation
- Chroot
- Process namespacing

## Build

```shell
make
```

## Run a command in a container

```shell
sudo ./bin/container-from-scratch run alpine /bin/sh
```

```
sudo ./bin/container-from-scratch run alpine /bin/sh
CHILD PID: 1
CHILD Hostname: container
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /proc/self/exe child alpine /bin/sh
    7 root      0:00 /bin/sh
    8 root      0:00 ps
/ # 
```