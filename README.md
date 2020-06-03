# hostingDeUpdateIP

_hostingDeUpdateIP_ updates DNS entries for dynamically changed IPv4 and IPv6 addresses. It works with the API of `hosting.de`

The client reads out the external IPv4 address and the assigned IPv6 address, based on a _postfix_.

## Command line arguments

```bash
-cf string
  	ConfigFile (default "config.json")
-ec
  	Show Example Config
-h	Show Help
-v	Show Version
```



## Configuration

### Example Configuration

The example shows only the mandatory elements. A full example can be printed with
`hostingDeUpdateIP -ec`

```json
{
   "HostingDe": {
      "Api": {
         "AuthToken": "NOT-KEY-SET"
      }
   },
   "Domains": [
      {
         "Host": "domain1.tld",
         "SetHostToo": true,
         "Subs": [
            "www",
            "www2"
         ]
      },
      {
         "Host": "domain2.tld",
         "SetHostToo": false,
         "Subs": [
            "blog"
         ]
      }
   ],
   "ZoneUpdateDefault": {
      "ZoneConfig": {
         "EmailAddress": "admin@domain.tld"
      }
   },
   "Ipv6Postfix": "::1234:5678",
   "LogLevel": "Debug"
}
```

`AuthToken` holds your API-Token from `hosting.de`.

The array `Domains` holds the domain, named `Host` and the sub domains as array, named `Subs`. `SetHostToo` define if the domain also should be updated or only the sub domains. `Host` is necessary to identify the Zone of the sub domains.

`Ipv6Postfix` holds the IPv6 postfix. Your _ISP_ delegates a IPv6 prefix to your router which often changes after router or modem reboot. A _static_ IPv6 address can be assigned to an interface via _IPv6 DHCP_ service. Most routers allow the assignment of an IPv6 postfix.