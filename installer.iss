[Setup]
AppName=Aharsuchi Printer Agent
AppVersion=1.1.0
DefaultDirName={pf}\Aharsuchi
DefaultGroupName=Aharsuchi
PrivilegesRequired=admin
OutputBaseFilename=AharsuchiPrinterAgentSetup
Compression=lzma
SolidCompression=yes

[Files]
Source: "printer-agent.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "certs\cert.pem"; DestDir: "{app}\certs"
Source: "certs\cert-key.pem"; DestDir: "{app}\certs"
Source: "certs\rootCA.pem"; DestDir: "{tmp}"

[Run]
; Install trusted root certificate
Filename: "certutil.exe"; \
  Parameters: "-addstore Root ""{tmp}\rootCA.pem"""; \
  Flags: runhidden runascurrentuser

; Stop old service if upgrading (ignore errors)
Filename: "{app}\printer-agent.exe"; \
  Parameters: "stop"; \
  Flags: runhidden; StatusMsg: "Stopping existing service..."

; Uninstall old service if upgrading (ignore errors)
Filename: "{app}\printer-agent.exe"; \
  Parameters: "uninstall"; \
  Flags: runhidden; StatusMsg: "Removing old service registration..."

; Register as a Windows Service (Automatic start)
Filename: "{app}\printer-agent.exe"; \
  Parameters: "install"; \
  Flags: runhidden; StatusMsg: "Installing Windows Service..."

; Configure automatic restart on failure:
;   - Reset failure counter after 24 hours (86400 seconds)
;   - On 1st, 2nd, and subsequent failures: restart after 10 seconds (10000 ms)
Filename: "sc.exe"; \
  Parameters: "failure AharsuchiPrinterAgent reset= 86400 actions= restart/10000/restart/10000/restart/10000"; \
  Flags: runhidden; StatusMsg: "Configuring auto-restart..."

; Start the service
Filename: "{app}\printer-agent.exe"; \
  Parameters: "start"; \
  Flags: runhidden; StatusMsg: "Starting Printer Agent service..."

[UninstallRun]
; Stop and remove the Windows Service on uninstall
Filename: "{app}\printer-agent.exe"; \
  Parameters: "stop"; \
  Flags: runhidden

Filename: "{app}\printer-agent.exe"; \
  Parameters: "uninstall"; \
  Flags: runhidden

[UninstallDelete]
Type: files; Name: "{app}\printer-agent.exe"
Type: dirifempty; Name: "{app}\certs"
Type: dirifempty; Name: "{app}"