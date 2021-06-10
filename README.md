# Jato
Manage network devices

## Requirements
Go ~v1.11+

## Supported Platforms
| Vendor  | Platform | SSH | Telnet |
|---------|----------|-----|--------|
| Cisco   | IOS      | :heavy_check_mark: | :heavy_check_mark: |
| Juniper | Junos    | :heavy_check_mark: | :x: |
| Arista  | EOS      | :heavy_check_mark: | :x: |

* :heavy_check_mark: - Implemented
* :x: - Not Implemented, but planned
* :red_circle: - Not planned

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
  "commands": [
    "show version",
    "show ip interface brief",
    "show ip arp",
    "show cdp neighbors",
    "show running-config"
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
  -a    Ask for user password
  -c string
        Commands to run file (default "commands.json")
  -d string
        Devices inventory file (default "devices.json")
  -noop
        Don't execute job against devices
  -u string
        Username to connect to devices with
  -v    Jato version
```

Run a series of commands against N number of devices.
```
./jato -d test/devices/cisco.json -c test/commands/cisco_ios.json
```
The results will be saved to a `output/` directory. The result of the 
commands run will be saved to a time stamped files. Once with the 
raw output and another with a json array of the command to value 
hash.
