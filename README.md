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
./jato -d test/devices/cisco_iosxr.json -c test/commands/cisco_iosxr.json
```
The results will be saved to a `output/` directory. The result of the 
commands run will be saved to a time stamped files. Once with the 
raw output and another with a json array of the command to value 
hash.

### Example JSON ouput
```json
{
 "device": "iosxr-1",
 "ok": true,
 "error": null,
 "timestamp": 1623629887,
 "commandOutputs": [
  {
   "command": "show version",
   "output": "\rMon Jun 14 00:18:06.005 UTC\r\n\r\nCisco IOS XR Software, Version 6.1.3[Default]\r\nCopyright (c) 2017 by Cisco Systems, Inc.\r\n\r\nROM: GRUB, Version 1.99(0), DEV RELEASE\r\n\r\nios uptime is 3 minutes\r\nSystem image file is \"bootflash:disk0/xrvr-os-mbi-6.1.3/mbixrvr-rp.vm\"\r\n\r\ncisco IOS XRv Series (Pentium Celeron Stepping 3) processor with 3145215K bytes of memory.\r\nPentium Celeron Stepping 3 processor at 2792MHz, Revision 2.174\r\nIOS XRv Chassis\r\n\r\n3 GigabitEthernet\r\n1 Management Ethernet\r\n97070k bytes of non-volatile configuration memory.\r\n866M bytes of hard disk.\r\n2321392k bytes of disk0: (Sector size 512 bytes).\r\n\r\nConfiguration register on node 0/0/CPU0 is 0x2102\r\nBoot device on node 0/0/CPU0 is disk0:\r\nPackage active on node 0/0/CPU0:\r\niosxr-infra, V 6.1.3[Default], Cisco Systems, at disk0:iosxr-infra-6.1.3\r\n    Built on Mon Feb 13 15:01:56 UTC 2017\r\n    By iox-lnx-005 in /auto/srcarchive14/production/6.1.3/xrvr/workspace for pie\r\n\r\niosxr-fwding, V 6.1.3[Default], Cisco Systems, at disk0:iosxr-fwding-6.1.3\r\n    Built on Mon Feb 13 15:01:56 UTC 2017\r\n    By iox-lnx-005 in /auto/srcarchive14/production/6.1.3/xrvr/workspace for pie\r\n<snip>"
  },
  {
   "command": "show ip interface brief",
   "output": "\rMon Jun 14 00:18:07.274 UTC\r\n\r\nInterface                      IP-Address      Status          Protocol Vrf-Name\r\nMgmtEth0/0/CPU0/0              192.168.255.162 Up              Up       default \r\nGigabitEthernet0/0/0/0         unassigned      Shutdown        Down     default \r\nGigabitEthernet0/0/0/1         unassigned      Shutdown        Down     default \r\nGigabitEthernet0/0/0/2         unassigned      Shutdown        Down     default "
  },
  {
   "command": "show arp",
   "output": "\rMon Jun 14 00:18:07.624 UTC\r\n\r\n-------------------------------------------------------------------------------\r\n0/0/CPU0\r\n-------------------------------------------------------------------------------\r\nAddress         Age        Hardware Addr   State      Type  Interface\r\n192.168.255.51  00:00:04   64bc.58e3.4fc7  Dynamic    ARPA  MgmtEth0/0/CPU0/0\r\n192.168.255.81  00:00:29   788a.20c5.975f  Dynamic    ARPA  MgmtEth0/0/CPU0/0\r\n192.168.255.85  00:00:23   788a.2080.2be6  Dynamic    ARPA  MgmtEth0/0/CPU0/0\r\n192.168.255.87  00:00:09   788a.20d0.f5bd  Dynamic    ARPA  MgmtEth0/0/CPU0/0\r\n192.168.255.88  00:00:31   c8d0.83c6.0f26  Dynamic    ARPA  MgmtEth0/0/CPU0/0\r\n192.168.255.90  00:00:14   788a.205c.d7fd  Dynamic    ARPA  MgmtEth0/0/CPU0/0\r\n192.168.255.91  00:00:23   788a.2046.905c  Dynamic    ARPA  MgmtEth0/0/CPU0/0\r\n192.168.255.162 -          5000.000c.0000  Interface  ARPA  MgmtEth0/0/CPU0/0"
  },
  {
   "command": "show cdp neighbors",
   "output": "\rMon Jun 14 00:18:07.974 UTC\r\n% CDP is not enabled"
  },
  {
   "command": "show running-config",
   "output": "\rMon Jun 14 00:18:08.334 UTC\r\nBuilding configuration...\r\n!! IOS XR Configuration 6.1.3\r\n!! Last configuration change at Sun Jun 13 04:23:17 2021 by cisco\r\n!\r\ninterface MgmtEth0/0/CPU0/0\r\n ipv4 address 192.168.255.162 255.255.255.0\r\n!\r\ninterface GigabitEthernet0/0/0/0\r\n shutdown\r\n!\r\ninterface GigabitEthernet0/0/0/1\r\n shutdown\r\n!\r\ninterface GigabitEthernet0/0/0/2\r\n shutdown\r\n!\r\nssh server v2\r\nend\r\n"
  }
 ]
}
```

### Example RAW ouput
```
! Device:    iosxr-1
! Timestamp: 1623629887
! OK:        true
! Error:     %!s(<nil>)
!----------------------------------------------------------!
!                       show version                       !
!----------------------------------------------------------!

Mon Jun 14 00:18:06.005 UTC

Cisco IOS XR Software, Version 6.1.3[Default]
Copyright (c) 2017 by Cisco Systems, Inc.

ROM: GRUB, Version 1.99(0), DEV RELEASE

ios uptime is 3 minutes
System image file is "bootflash:disk0/xrvr-os-mbi-6.1.3/mbixrvr-rp.vm"

cisco IOS XRv Series (Pentium Celeron Stepping 3) processor with 3145215K bytes of memory.
Pentium Celeron Stepping 3 processor at 2792MHz, Revision 2.174
IOS XRv Chassis

3 GigabitEthernet
1 Management Ethernet
97070k bytes of non-volatile configuration memory.
866M bytes of hard disk.
2321392k bytes of disk0: (Sector size 512 bytes).

Configuration register on node 0/0/CPU0 is 0x2102
Boot device on node 0/0/CPU0 is disk0:
Package active on node 0/0/CPU0:
iosxr-infra, V 6.1.3[Default], Cisco Systems, at disk0:iosxr-infra-6.1.3
    Built on Mon Feb 13 15:01:56 UTC 2017
    By iox-lnx-005 in /auto/srcarchive14/production/6.1.3/xrvr/workspace for pie

<snip>
!----------------------------------------------------------!
!                 show ip interface brief                  !
!----------------------------------------------------------!

Mon Jun 14 00:18:07.274 UTC

Interface                      IP-Address      Status          Protocol Vrf-Name
MgmtEth0/0/CPU0/0              192.168.255.162 Up              Up       default 
GigabitEthernet0/0/0/0         unassigned      Shutdown        Down     default 
GigabitEthernet0/0/0/1         unassigned      Shutdown        Down     default 
GigabitEthernet0/0/0/2         unassigned      Shutdown        Down     default 
!----------------------------------------------------------!
!                         show arp                         !
!----------------------------------------------------------!

Mon Jun 14 00:18:07.624 UTC

-------------------------------------------------------------------------------
0/0/CPU0
-------------------------------------------------------------------------------
Address         Age        Hardware Addr   State      Type  Interface
192.168.255.51  00:00:04   64bc.58e3.4fc7  Dynamic    ARPA  MgmtEth0/0/CPU0/0
192.168.255.81  00:00:29   788a.20c5.975f  Dynamic    ARPA  MgmtEth0/0/CPU0/0
192.168.255.85  00:00:23   788a.2080.2be6  Dynamic    ARPA  MgmtEth0/0/CPU0/0
192.168.255.87  00:00:09   788a.20d0.f5bd  Dynamic    ARPA  MgmtEth0/0/CPU0/0
192.168.255.88  00:00:31   c8d0.83c6.0f26  Dynamic    ARPA  MgmtEth0/0/CPU0/0
192.168.255.90  00:00:14   788a.205c.d7fd  Dynamic    ARPA  MgmtEth0/0/CPU0/0
192.168.255.91  00:00:23   788a.2046.905c  Dynamic    ARPA  MgmtEth0/0/CPU0/0
192.168.255.162 -          5000.000c.0000  Interface  ARPA  MgmtEth0/0/CPU0/0
!----------------------------------------------------------!
!                    show cdp neighbors                    !
!----------------------------------------------------------!

Mon Jun 14 00:18:07.974 UTC
% CDP is not enabled
!----------------------------------------------------------!
!                   show running-config                    !
!----------------------------------------------------------!

Mon Jun 14 00:18:08.334 UTC
Building configuration...
!! IOS XR Configuration 6.1.3
!! Last configuration change at Sun Jun 13 04:23:17 2021 by cisco
!
interface MgmtEth0/0/CPU0/0
 ipv4 address 192.168.255.162 255.255.255.0
!
interface GigabitEthernet0/0/0/0
 shutdown
!
interface GigabitEthernet0/0/0/1
 shutdown
!
interface GigabitEthernet0/0/0/2
 shutdown
!
ssh server v2
end
```