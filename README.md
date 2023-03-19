# Netclip

Self-contained pastebin for local networks, with no dependencies. Runs as a service on Windows too.

### Why

I need to copy text between computers, and using a public pastebin isn't always an option. I need multiple computers 
to be able to paste things and see things pasted.

And once this is running it's quicker than transferring files.

## Install

Download the binary for your OS.

### Run

```
netclip
```

This runs the server on 0.0.0.0 port `9999`. Override the port with `-p 4000`.

### Run as a service

This supports running as a service on Windows, macOS, and Linux.

```
netclip -service install
netclip -service start
```

Stop the service  with

```
netclip -service stop
```

Uninstall  the service  with

```
netclip -service uninstall
```

This works on Windows 7 and above, and Windows Server 2008 and up. 

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


## History

2023-03-18

0.0.5 - single binary with https support

2023-03-17

0.0.4 - refactor binaries to use templates and css
0.0.3 - initial public test
0.0.2 - service support
0.0.1 - it's alive!

## License

MIT
