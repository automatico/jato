# Jato
Manage network devices

## Requirements
Go ~v1.11+

## Supported Platforms
| Vendor  | Platform | SSH | Telnet |
|---------|----------|-----|--------|
| Arista  | EOS      | :heavy_check_mark: | :x: |
| Aruba   | AOS-CX   | :heavy_check_mark: | :red_circle: |
| Cisco   | AireOS   | :heavy_check_mark: | :red_circle: |
| Cisco   | IOS      | :heavy_check_mark: | :heavy_check_mark: |
| Cisco   | IOS-XR   | :heavy_check_mark: | :red_circle: |
| Cisco   | NXOS     | :heavy_check_mark: | :red_circle: |
| Cisco   | SMB      | :heavy_check_mark: | :red_circle: |
| Juniper | Junos    | :heavy_check_mark: | :x: |


* :heavy_check_mark: - Supported
* :x: - Support planned, but not yet implemented
* :red_circle: - Not Supported

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
    {"name": "eos-1", "ip": "192.168.255.152", "vendor": "arista", "platform": "eos", "connector": "ssh"},
    {"name": "aoscx-1", "ip": "192.168.255.161", "vendor": "aruba", "platform": "aoscx", "connector": "ssh"},
    {"name": "aireos-1", "ip": "192.168.255.164", "vendor": "cisco", "platform": "aireos", "connector": "ssh"},
    {"name": "iosv-1", "ip": "192.168.255.150", "vendor": "cisco", "platform": "ios", "connector": "telnet"},
    {"name": "iosxr-1", "ip": "192.168.255.162", "vendor": "cisco", "platform": "iosxr", "connector": "ssh"},
    {"name": "nsox-1", "ip": "192.168.255.163", "vendor": "cisco", "platform": "nxos", "connector": "ssh"},
    {"name": "smb-1", "ip": "192.168.255.154", "vendor": "cisco", "platform": "smb", "connector": "ssh"},
    {"name": "vmx-1", "ip": "192.168.255.151", "vendor": "juniper", "platform": "junos", "connector": "ssh"}
  ]
}
```
### Configuration Parameters
| vendor  | platform | connector   |
|---------|----------|-------------|
| arista  | eos      | ssh         |
| aruba   | aoscx    | ssh         |
| cisco   | aireos   | ssh         |
| cisco   | ios      | ssh, telnet |
| cisco   | iosxr    | ssh         |
| cisco   | nxos     | ssh         |
| cisco   | smb      | ssh         |
| juniper | junos    | ssh         |

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
