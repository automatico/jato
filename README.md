# Jato
Manage network devices

## Setup

### Environment variables
Create environment variables

### Linux
Place this in your environments `rc` file
```
export JATO_USERNAME="USERNAME"
export JATO_PASSWORD="PASSWORD"
```

### Windows
In an elevated Powershell / Command prompt
```
setx JATO_USERNAME "USERNAME"
setx JATO_PASSWORD "PASSWORD"
```

### Commands
Create a `commands.json` file with the list of commands to run
```json
{
  "command_expect": [
    {"command": "terminal length 0", "expecting": "#", "timeout": 5},
    {"command": "show version", "expecting": "#", "timeout": 5},
    {"command": "show ip interface brief", "expecting": "#", "timeout": 5},
    {"command": "show ip arp", "expecting": "#", "timeout": 5},
    {"command": "show cdp neighbors", "expecting": "#", "timeout": 5},
    {"command": "show running-config", "expecting": "#", "timeout": 5},
    {"command": "exit", "expecting": "#", "timeout": 5}
  ]
}
```

### Devices
Create a `devices.json` file with a list of devices to run against
```json
{
  "devices": [
    {"name": "iosv-1", "ip": "192.168.255.150", "vendor": "cisco", "platform": "ios", "connector": "telnet"},
    {"name": "iosv-4", "ip": "192.168.255.154", "vendor": "cisco", "platform": "ios", "connector": "ssh"},
    {"name": "iosv-5", "ip": "192.168.255.155", "vendor": "cisco", "platform": "ios", "connector": "ssh"},
    {"name": "iosv-6", "ip": "192.168.255.156", "vendor": "cisco", "platform": "ios", "connector": "ssh"},
    {"name": "iosv-7", "ip": "192.168.255.157", "vendor": "cisco", "platform": "ios", "connector": "ssh"},
    {"name": "iosv-8", "ip": "192.168.255.158", "vendor": "cisco", "platform": "ios", "connector": "telnet"},
    {"name": "iosv-9", "ip": "192.168.255.159", "vendor": "cisco", "platform": "ios", "connector": "telnet"},
    {"name": "iosv-10", "ip": "192.168.255.160", "vendor": "cisco", "platform": "ios", "connector": "telnet"}
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
        Username to connect to devices with (default "JATO_USERNAME environment variable")
```

Run a series of commands against N number of devices.
```
./jato -d test/devices/cisco.json -c test/commands/cisco_ios.json
```
The results will be saved to a `data/` directory. The result of the 
commands run will be saved to a time stamped files. Once with the 
raw output and another with a json array of the command to value 
hash.
