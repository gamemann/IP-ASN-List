# IP/ASN List
## Description
A Go program that acts as a web server. This program loops through all `.json` files in the `lists/` directory and parses each list. Each list can contain an array of prefixes and ASNs. There are lookups performed on the ASNs via [BPFView](https://bgpview.docs.apiary.io/#reference/0/asn-prefixes/view-asn-prefixes)'s API. Afterwards, it outputs all prefixes to a file in the `public/` directory where clients can view them assuming they have the correct authorization header set from the config.

## Why Did I Make This?
After implementing my own custom filters into [Compressor V1](https://github.com/Dreae/compressor/) to mitigate (D)DoS attacks against GFL's Anycast network, I decided to implement whitelisting functionality for specific services. However, I didn't want to update these files locally on each POP each time I had to add a prefix. Therefore, I decided to create this program and have each POP cURL each list for specific services (with the correct authorization header set) and save them locally every five - ten minutes.

## Config
You can change settings in the `settings.conf` file which is in JSON format. Here's the default config:

```
{"token": "CHANGEME", "port": 7030, "updatetime": 15}
```

* `token` => The authorization header that must be set when accessing the list publicly.
* `port` => The port the web server binds to.
* `updatetime` => How often to update all lists.

## Lists
Each list config needs to have a file with the format `lists/<name>.json`. Here's an example of the test list (`lists/test.json`):

```
{
    "ASN": [
        398129,
        32590
    ],
    "Prefix": [
        "192.168.90.1/32",
        "192.168.90.2/32"
    ]
}
```

The above configuration will output all prefixes from ASN's 398129 and 32590 along with the additional prefixes `192.168.90.1/32` and `192.168.90.2/32` to `public/test.txt`.

## Building
You may use the following commands to build this project.

```
cd src/
go build -o ../iplist
```

## Credits
* [Christian Deacon](https://www.linkedin.com/in/christian-deacon-902042186/) - Creator.
* [BGPView](https://bgpview.docs.apiary.io/#reference/0/asn-prefixes/view-asn-prefixes) - ASN lookup API.