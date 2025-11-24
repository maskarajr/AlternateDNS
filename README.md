<p align="center"><img src="https://github.com/MaxIsJoe/AlternateDNS/blob/master/logo.webp"></p>

<p align="center">
  <strong>Enhanced Fork with GUI, Smart DNS Switching, and DNS Testing</strong><br>
  <small>Original project by <a href="https://github.com/MaxIsJoe">MaxIsJoe</a> | <a href="#credits">Credits</a></small>
</p>

# AlternateDNS (Enhanced Fork)

A tool that changes your DNS settings periodically every couple of hours based on your configuration.

> **⚠️ Important:** This is an enhanced fork of the original AlternateDNS project. See [Credits](#credits) section below.


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


## Features

### Original Features
- Periodic DNS rotation based on configuration
- Cross-platform support (Windows, Linux, macOS)
- System tray integration
- Startup registration
- Notification support

### Enhanced Features (This Fork)
- **Modern GUI**: Full-featured desktop application with Fyne framework
- **Smart DNS Switching**: Automatically tests and compares DNS latency before switching
- **DNS Tester**: Built-in tool to test and benchmark DNS servers
- **Portable Build**: Single executable with embedded resources
- **Enhanced Logging**: Real-time logs with better visibility
- **Manual Override**: Force DNS change option that bypasses latency testing

## How do I run this?

### Quick Start (Portable)

**Windows:**
1. Download or build `AlternateDNS.exe`
2. Place `config.yaml` in the same directory (or let it auto-generate)
3. Run `AlternateDNS.exe` (requires admin privileges for DNS changes)

**Linux/macOS:**
1. Clone the repo
2. Install Go and dependencies
3. Run `./build.sh` or `go build -o AlternateDNS`
4. Run with `sudo ./AlternateDNS` (requires root for DNS changes)

### Building from Source

1. Clone the repo
2. Install Go 1.22.3 or later
3. Install dependencies: `go mod download`
4. Build:
   - Windows: `build.bat` or `go build -o AlternateDNS.exe`
   - Linux/macOS: `./build.sh` or `go build -o AlternateDNS`
5. Run the executable (with appropriate privileges)

**Note:** For Windows GUI builds, you'll need MinGW-w64 (GCC) for CGO support.

This will generally work on Windows and Linux, though I'm not sure about Mac because I don't own one.

## Configuration

The application uses `config.yaml` for configuration. A default configuration is embedded in the executable and will be created on first run if not present.

Example `config.yaml`:
```yaml
dns_addresses:
  - 9.9.9.9
  - 149.112.112.112
  - 1.1.1.1
  - 1.0.0.1
run_on_startup: true
change_interval_hours: 6
notify_user: true
test_domains:
  - google.com
  - cloudflare.com
  - github.com
  - microsoft.com
  - amazon.com
```

## Contributing

We welcome contributions! This is a community-maintained fork, and we'd love your help improving it.

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to this project.

**Note:** Contributions are made to this fork, not the original repository. All contributions will maintain proper attribution to the original author.

## Credits

**Original Author:** [MaxIsJoe](https://github.com/MaxIsJoe)  
**Original Repository:** https://github.com/MaxIsJoe/AlternateDNS

This repository is a **community-maintained fork** with additional features. All original functionality and code structure are credited to [MaxIsJoe](https://github.com/MaxIsJoe). We maintain this fork to add GUI features, smart DNS switching, and DNS testing capabilities while preserving the original author's work and giving proper attribution.

**Why a fork?** The original repository doesn't have contribution guidelines, so we've created this maintained fork to:
- Add new features and improvements
- Accept contributions from the community
- Keep the original author's work properly credited
- Provide ongoing maintenance and updates

See [AUTHORS.md](AUTHORS.md) for detailed attribution information.

## License

Please refer to the original repository for license information.

## Support Original Author

Support the original creator:
https://www.maxisjoe.xyz/maxfund
