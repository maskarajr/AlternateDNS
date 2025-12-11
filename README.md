<p align="center"><img src="https://github.com/MaxIsJoe/AlternateDNS/blob/master/logo.webp"></p>

# AlternateDNS

A tool that changes your DNS settings periodically every couple of hours based on your configuration.

## Huh? Why?

I don't know. I woke up one day with internet issues where I couldn't access a lot of websites because they wouldn't get resolved; and quickly found out that the only solution is to just change my DNS settings manually every time this happened. After weeks of having to go through this, I decided to automate the process while I'm learning Go. 

I've done something similar like this in Python waaaay back on my Intel Core 2 duo laptop. So this proved to be a good learning project since it is a familiar task for me.

## Are you associated with Alternate DNS, the DNS service provider?

No.

AlternateDNS and Alternate DNS are two completely different projects, one is a computer program designed to change your DNS settings periodically, the other provides you a DNS address to use.

Though I heavily discourage people from using Alternate DNS as their privacy policy seems quite fishy: https://web.archive.org/web/20231207072356/https://alternate-dns.com/privacy.html

If you want DNS recommendations from me, here:

| DNS Provider   | Address                                                     | My Rating | Reason                                                                                                                                                                                                                                                                                         |
|----------------|-------------------------------------------------------------|-----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Mullvad DNS    | https://mullvad.net/en/help/dns-over-https-and-dns-over-tls | 10/10     | I love Mullvad. They respect your privacy, and their whole setup is great. Never had a single issue with them and their product.                                                                                                                                                                          |
| Quad9          | https://quad9.net/                                          | 8.5/10    | Great service for privacy and security. Though, it lacks behind in performance from my experience.                                                                                                                                                                                             |
| Cloudflare DNS | https://one.one.one.one/                                    | 7/10      | I generally do not trust Cloudflare due to a lot of their controversies.  However, they provide a lot of great services; and their DNS service has been always reliable for me.                                                                                                                                                                    |
| 8.8.8.8        | https://dns.google/                                         | -1/10     | It's Google, what more do you want me to say? I am only recommending this because in some scenarios, you might not be able to use any of the above services. Google's DNS is reliable and works well for everyone, but is censored in some areas around the world, and not privacy respecting. |


## How do I run this?

1. Clone the repo
2. Install Go 1.22.3 or later
3. Install dependencies: `go mod download`
4. Build:
   - Windows: `build.bat` or `go build -o AlternateDNS.exe` (requires MinGW-w64/GCC for GUI)
   - Linux/macOS: `./build.sh` or `go build -o AlternateDNS`
5. Run the executable (requires admin/root privileges for DNS changes)

**Note:** The application now includes a modern GUI interface. For Windows GUI builds, you'll need MinGW-w64 (GCC) for CGO support. The application will create a `config.yaml` file on first run if one doesn't exist.

This will generally work on Windows and Linux, though I'm not sure about Mac because I don't own one.

## Features

- Periodic DNS rotation based on configuration
- Cross-platform support (Windows, Linux, macOS)
- System tray integration
- Startup registration
- Notification support
- Modern GUI with Fyne framework
- Smart DNS switching with latency comparison
- Built-in DNS tester for benchmarking servers
- Portable single-executable builds

## Support me:

https://www.maxisjoe.xyz/maxfund
