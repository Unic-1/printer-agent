[Setup]
AppName=Aharsuchi Printer Agent
AppVersion=1.0.0
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

; Start agent silently
Filename: "{app}\printer-agent.exe"; Flags: nowait runhidden

[Registry]
; Run agent at startup
Root: HKLM; Subkey: "Software\Microsoft\Windows\CurrentVersion\Run"; \
  ValueType: string; ValueName: "AharsuchiPrinterAgent"; \
  ValueData: """{app}\printer-agent.exe"""; Flags: uninsdeletevalue
