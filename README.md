# dns-resolver


# What is a DNS resolver?

The part that deals with how your system translates the host part of the URL to an IP address, i.e. from this: dns.google.com to 8.8.8.8 or 8.8.4.4. In other words, domain name resolution; turning the hostname that your browser has extracted from the URL into an IP address.

To do this your browser contacts a DNS Resolver. The DNS Resolve may have the answer in it’s cache in which case it can return it immediately, if not then it will have to look it up. To look it up it contacts an authoritative name server. To do that it will first consult it’s cache for an authoritative name server and if it doesn’t have one it will contact a root name server to get one.

## What is an Authoritative DNS server? 

Authoritative DNS server provide answers to DNS resolvers with the correct IP addresses and records.These servers covers different regions, like a company, the local area or your country. No matter what region it covers, an authoritative DNS server performs two important task. First, it stores lists of domain names and their associated IP addresses. Second, it responds to requests from a DNS resolver about the correct IP address assigned to a domain name. It does not simply return cached responses from another name server, but rather provides answers to queries about domain names that are configured in its system.

There are three main levels of authoritative DNS servers. There are the *root nameservers*, the *top level domain (or TLD) nameservers* and the authoritative nameservers  for the domains.

1) The root nameservers store the IP addresses of the TLD name servers. There are 13 logical root nameservers, each one with a single IP address assigned to it, although actually they have many physical servers behind each IP address and we can get routed to the one closest to us.

2) The TLD nameservers store the IP addresses of all the domain above them. Some of them are very popular, like .com, .org, .br, .uk. 

3) The others are known by the TLD based on their region. Some authoritative name servers are very popular, from cloud companies that also register nameservers like AWS and Cloudfare.

The hierarchical design makes DNS highly decentralized and robust.

*DNS propagation* is slow because there is a TTL on each DNS record and some of the default TTLs are pretty long.

# DNS message format

  * A header.
  * A questions section.
  * An answer section.
  * An authority section.
  * An additional section.

  Traditional DNS messages are limited to 512 bytes in size when sent over UDP [RFC1035].

## Big endian byte order*

  * In computing, endianness is the order in which bytes within a word of digital data are stored in computer memory or transmitted over a data communication medium. Endianness is primarily expressed as big-endian (BE) or little-endian (LE)

  *A big-endian system* stores the most significant byte of a word at the smallest memory address and the least significant byte at the largest. *A little-endian system*, in contrast, stores the least-significant byte at the smallest address. Of the two, big-endian is thus closer to the way the digits of numbers are written left-to-right in English, comparing digits to bytes

## Question section

  [RFC 4.1.2](https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1)

  *QName*: Domain name represented as a sequence of labels. Where each label consists of a length octet followed by that number of octets
    -> Each element of a domain name separated by [.] is called a “label.” The maximum length of each label is 63 characters, and a full domain name can have a maximum of 253 characters. ([Reference](https://www.nic.ad.jp/timeline/en/20th/appendix1.html#:~:text=Format%20of%20a%20domain%20name,a%20maximum%20of%20253%20characters.))

  *QType*: 16 bits code that specifies the type of the query. [RFC Values](https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.2)

  *QClass*: 16 bits code that specifies the class of the query. Ex: The class field is IN for the Internet. [RFC Values](https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.4)

## What is UDP ?

  In computer networking, the [User Datagram Protocol](https://en.wikipedia.org/wiki/User_Datagram_Protocol) (UDP) is one of the core communication protocols of the Internet protocol suite used to send messages (transported as datagrams in packets) to other hosts on an Internet Protocol (IP) network. Within an IP network, *UDP does not require prior communication to set up communication channels or data paths*.

# Study Resources

[Everything You Need to Know About DNS: Crash Course System Design #4](https://www.youtube.com/watch?v=27r4Bzuj5NQ)

[What is the difference between authoritative and recursive DNS nameservers?](https://umbrella.cisco.com/blog/what-is-the-difference-between-authoritative-and-recursive-dns-nameservers)

[What is Authoritative DNS server?](https://www.cloudns.net/blog/authoritative-dns-server/#:~:text=The%20authoritative%20DNS%20server%20is,information%20cached%20in%20its%20memory.)
