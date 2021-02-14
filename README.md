# Jato
Manage network devices

## Setup

### Environment variables
Create environment variables

### Linux
Place this in your environments `rc` file
```
export JATO_SSH_USER="USERNAME"
export JATO_SSH_PASS="PASSWORD"
```

### Windows
In an elevated Powershell / Command prompt
```
setx JATO_SSH_USER "USERNAME"
setx JATO_SSH_PASS "PASSWORD"
```

### Commands
Create a `commands.json` file with the list of commands to run
```json
{
  "commands": [
    "terminal length 0",
    "show version",
    "show ip interface brief",
    "show ip arp",
    "show cdp neighbors",
    "show running-config",
    "exit"
  ]
}
```

### Devices
Create a `devices.json` file with a list of devices to run against
```json
{
  "devices": [
    {"name": "192.168.255.150", "vendor": "cisco", "platform": "ios"},
    {"name": "192.168.255.154", "vendor": "cisco", "platform": "ios"},
    {"name": "192.168.255.155", "vendor": "cisco", "platform": "ios"},
    {"name": "192.168.255.156", "vendor": "cisco", "platform": "ios"},
    {"name": "192.168.255.157", "vendor": "cisco", "platform": "ios"}
  ]
}
```

## Run
```
# Linux
./jato

# Windows
./jato.exe
```
The results will be saved to a `data/` directory. The result of the 
commands run will be saved to a time stamped files. Once with the 
raw output and another with a json array of the command to value 
hash.
