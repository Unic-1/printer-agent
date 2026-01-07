[Setup]
AppName=Aharsuchi Printer Agent
AppVersion=1.0.0
DefaultDirName={pf}\Aharsuchi Printer Agent
DefaultGroupName=Aharsuchi
OutputDir=dist
OutputBaseFilename=aharsuchi-printer-agent-setup
Compression=lzma
SolidCompression=yes
PrivilegesRequired=admin

[Files]
Source: "build\printer-agent.exe"; DestDir: "{app}"; Flags: ignoreversion

[Registry]
; Custom protocol handler
Root: HKCR; Subkey: "aharsuchi-printer"; ValueType: string; ValueName: ""; ValueData: "URL:Aharsuchi Printer Protocol"
Root: HKCR; Subkey: "aharsuchi-printer"; ValueType: string; ValueName: "URL Protocol"; ValueData: ""
Root: HKCR; Subkey: "aharsuchi-printer\shell\open\command"; ValueType: string; ValueName: ""; ValueData: """{app}\printer-agent.exe"" ""%1"""

; Optional: auto-start on login
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Run"; \
  ValueType: string; ValueName: "AharsuchiPrinterAgent"; \
  ValueData: """{app}\printer-agent.exe"""

[Icons]
Name: "{group}\Aharsuchi Printer Agent"; Filename: "{app}\printer-agent.exe"
