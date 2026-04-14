#!/usr/bin/env python3
import os
import subprocess
import socket
from datetime import timedelta

def get_uptime():
    with open('/proc/uptime', 'r') as f:
        uptime_seconds = float(f.readline().split()[0])
        return str(timedelta(seconds=int(uptime_seconds)))

def get_ip():
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        s.connect(("8.8.8.8", 80))
        ip = s.getsockname()[0]
        s.close()
        return ip
    except Exception:
        return "Not Connected"

def get_load():
    with open('/proc/loadavg', 'r') as f:
        return f.readline().split()[:3]

def get_users():
    try:
        output = subprocess.check_output(['who']).decode('utf-8')
        return len(output.strip().split('\n'))
    except:
        return 0

def get_ssh_sessions():
    try:
        # Count lines in 'who' output containing 'pts/'
        output = subprocess.check_output(['who']).decode('utf-8')
        return len([line for line in output.split('\n') if 'pts/' in line])
    except:
        return 0

def get_temp():
    # Try common thermal zones
    for i in range(5):
        path = f'/sys/class/thermal/thermal_zone{i}/temp'
        if os.path.exists(path):
            with open(path, 'r') as f:
                temp = int(f.readline().strip()) / 1000
                return f"{temp:.1f}°C"
    return "N/A (Virtual)"

def main():
    cyan = "\033[1;36m"
    yellow = "\033[1;33m"
    reset = "\033[0m"
    bold = "\033[1m"

    print(f"{cyan}╔════════════════════════════════════════════════════════════╗{reset}")
    print(f"{cyan}║{reset}  {bold}BLOCK HUB{reset} - System Dashboard for Block OS v1.0         {cyan}║{reset}")
    print(f"{cyan}╠════════════════════════════════════════════════════════════╣{reset}")
    
    # Row 1: Uptime & IP
    print(f"{cyan}║{reset}  {yellow}Uptime:{reset} {get_uptime():<18} {yellow}Local IP:{reset} {get_ip():<16} {cyan}║{reset}")
    
    # Row 2: Load & Temp
    load = " ".join(get_load())
    print(f"{cyan}║{reset}  {yellow}Load:{reset}   {load:<18} {yellow}Temp:{reset}     {get_temp():<16} {cyan}║{reset}")
    
    # Row 3: Users & SSH
    print(f"{cyan}║{reset}  {yellow}Users:{reset}  {get_users():<18} {yellow}SSH:{reset}      {get_ssh_sessions():<16} {cyan}║{reset}")
    
    print(f"{cyan}╚════════════════════════════════════════════════════════════╝{reset}")

if __name__ == "__main__":
    main()
