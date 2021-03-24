#!/usr/bin/env python3

import socket
import random

HOST = '127.0.0.1'  # Standard loopback interface address (localhost)
PORT = 4000         # Port to listen on (non-privileged ports are > 1023)

random.seed(1337)   # Everyone sends the same numbers
count = 0

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((HOST, PORT))
    try:
        while (count < 100):
                value = f"{random.randint(0, 999999999):09d}\n"
                s.sendall(bytes(value.encode()))
                print(value, end='')
                count = count + 1
    finally:
        s.close()
