# Netclip

Self-contained pastebin for local networks, with no dependencies. Runs as a service on Windows, macOS, and Linux. Can also run privately on a Tailscale network (tailnet).

### Why

I need to copy text between computers, and using a public pastebin isn't always an option. I need multiple computers
to be able to paste things and see things pasted.

### Limitations

- It's not multi-user. So don't paste things you don't want others to see.
- You're responsible for your own security, firewalling, etc.
- There is no permanent persistence. Restarting the service clears the database.


## Install

Download the binary for your OS.

### Run

You can run the program with defaults with

```
netclip
```

This runs the server on `0.0.0.0` port `9999`. Override the port with `-port 4000`.

You can use the following command-line options, which you can see with `--help`:

```
  -port string
        Port to use (default: 9999)
  -cert string
        Path to SSL certificate file
  -key string
        Path to SSL private key file
  -service string
        install/restart/start/stop/uninstall
  -tailscale
        Enable Tailscale networking
  -tailscale-hostname string
        Tailscale hostname (default: netclip)
  -tailscale-tls
        Use HTTPS with Tailscale certificates
  -v    Prints current app version.
```

You can also specify options in a `netclip.yml` file:

```yaml
port: "4000"
cert_file: "netclip.crt"
key_file: "netclip.key"
```

The application searches for `netclip.yml` in the following locations (in priority order):

1. **Executable directory** - Same directory as the netclip binary.
2. **Current working directory** - Where you ran the command from.
3. **User home directory** - `~/.netclip.yml` or `~/netclip.yml`
4. **User config directory** - `~/.config/netclip/netclip.yml` or `$XDG_CONFIG_HOME/netclip/netclip.yml`
5. **System-wide locations**:
   - **Linux**: `/etc/netclip/netclip.yml`, `/usr/local/etc/netclip/netclip.yml`, `/opt/netclip/netclip.yml`
   - **macOS**: `/etc/netclip/netclip.yml`, `/usr/local/etc/netclip/netclip.yml`, `/opt/netclip/netclip.yml`, `/Library/Application Support/netclip/netclip.yml`
   - **Windows**: `%PROGRAMDATA%\netclip\netclip.yml`, `%PROGRAMFILES%\netclip\netclip.yml`


### Run as a service

This supports running as a service on Windows, macOS, and Linux.

```
netclip -service install
netclip -service start
```

Stop the service with

```
netclip -service stop
```

Uninstall  the service with

```
netclip -service uninstall
```

This works on Windows 7 and above, and Windows Server 2008 and up.

#### Service Configuration Requirements

**IMPORTANT**: When running as a system service, you need to use a configuration file because the service doesn't accept command line arguments.
Place the config file in one of these locations:

**Linux (systemd service)**
- `/etc/netclip/netclip.yml` (recommended)
- `/usr/local/etc/netclip/netclip.yml`

```bash
sudo mkdir -p /etc/netclip
sudo cp netclip.yml /etc/netclip/
```

**macOS (launchd service)**
- `/usr/local/etc/netclip/netclip.yml` (recommended)
- `/Library/Application Support/netclip/netclip.yml`

```bash
sudo mkdir -p /usr/local/etc/netclip
sudo cp netclip.yml /usr/local/etc/netclip/
```

**Windows Service**
- `C:\ProgramData\netclip\netclip.yml` (recommended)
- `C:\Program Files\netclip\netclip.yml`

PowerShell as Administrator:

```powershell
New-Item -ItemType Directory -Path "C:\ProgramData\netclip" -Force
Copy-Item "netclip.yml" "C:\ProgramData\netclip\"
```

## SSL support

Copying to the clipboard with JavaScript requires a secure connection. Run this behind a front-end with HTTPS and a reverse proxy or use self-signed certs.

### Using self-signed certs

Generate self-signed cert that's good for a year.

```
openssl genrsa -out netclip.key 2048
openssl req -new -x509 -sha256 -key netclip.key -out netclip.crt -days 365
```

The first command generates a 2048-bit RSA private key and saves it to a file named netclip.key.
The second command creates a new self-signed X.509 certificate called `netclip.crt` using the generated private key
that's valid for 365 days.

Then create the file `netclip.yml` and give it the paths to those keys:


```yaml
port: "4000"
cert_file: "netclip.crt"
key_file: "netclip.key"
```

You can also set these with flags, which override config file settings:

```
netclip -port 8080 -cert server.crt -key server.key
```


These certs aren't signed by an authority so your browser will prevent you from using the site unless you allow it, which is only temporary.

### Permanently add self-signed certs:

macOS (for Safari, Chrome, and Edge):

* Double-click on the server.crt file to open it in the "Keychain Access" application.
* Choose "System" or "login" under "Keychains" and click "Add."
* If prompted, enter your macOS user password to confirm the action.
* Find the imported certificate in the list, double-click on it, and expand the "Trust" section.
* Set "Secure Sockets Layer (SSL)" to "Always Trust" and close the window.
* If prompted, enter your macOS user password to confirm the action.
* Restart your browser for the changes to take effect.

Windows (for Chrome and Edge):

* Open the site using the self-signed certificate in Chrome or Edge. You'll see a warning message.
* Click on the lock icon in the address bar, then click on "Certificate."
* Go to the "Details" tab and click on "Copy to File" to start the Certificate Export Wizard.
* Export the certificate using the default settings (DER encoded binary X.509).
* Press Win + R, type certmgr.msc, and press Enter to open the Certificate Manager.
* Go to "Trusted Root Certification Authorities" > "Certificates" and right-click on it.
* Choose "All Tasks" > "Import" to start the Certificate Import Wizard.
* Import the exported certificate file and complete the wizard.
* Restart your browser for the changes to take effect.

## Tailscale support

You can run netclip on your Tailscale network instead of a regular network. This way you can access it from any device on your tailnet without exposing ports or dealing with firewalls.

You should get an auth key and tag for this service. Once you have it, set the `TS_AUTHKEY` environment variable.

```
export TS_AUTHKEY="tskey-auth-xxxxx"
netclip -tailscale -tailscale-hostname my-netclip
```

If you don't have a Tailscale auth key set up, it'll walk you through logging in.

```
netclip -tailscale -tailscale-hostname my-netclip
```

Netclip will be available at `https://my-netclip.your-tailnet.ts.net` to all devices on your tailnet. Add `-tailscale-tls` to use HTTPS with automatic Tailscale certificates.

You can configure Tailscale support in the config file as well, but your auth key needs to come from the environment.

```yaml
tailscale:
  enabled: true
  hostname: "my-netclip"
  use_tls: true
```

## Changelog

### 0.6.1 - 2025-06-24

- Configuration file lookup which makes it easier to run as a service.

### 0.6.0 - 2025-06-22

- Exposed cert and key as config options
- Use netclip on your Tailscale network with the `-tailscale` flag.
- Flags now consistently override config file options.
- Improved internals to allow for better testing.

### 0.5.0 - 2023-03-18

- single binary with https support

### 0.4.0 2023-03-17

- refactor binaries to use templates and css

### 0.3.0 2023-03-17
- initial public test

### 0.2.0 2023-03-17
- service support

### 0.1.0 2023-03-17
- it's alive!

## License

MIT
