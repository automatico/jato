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
Inspect the options available
```
# Linux
./jato -h

# Windows
./jato.exe -h

Usage of jato:
  -c string
        Commands to run file (default "commands.json")
  -d string
        Devices inventory file (default "devices.json")
  -p    Ask for user password
  -u string
        Username to connect to devices with (default "JATO_SSH_USER environment variable")
```

Run a series of commands against N number of devices.
```
./jato -d test/devices/cisco.json -c test/commands/cisco_ios.json
```
The results will be saved to a `data/` directory. The result of the 
commands run will be saved to a time stamped files. Once with the 
raw output and another with a json array of the command to value 
hash.
