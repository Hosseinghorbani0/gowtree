; gowtree Windows Installer — Inno Setup 6+
; Build: .\scripts\build-installer.ps1

#define MyAppName "gowtree"
#define MyAppVersion "1.4.0"
#define MyAppPublisher "Hossein Ghorbani"
#define MyAppURL "https://github.com/hosseinghorbani0/gowtree"
#define MyAppExeName "gowtree.exe"

[Setup]
AppId={{A7B3C9D1-E5F2-4A8B-9C0D-1E2F3A4B5C6D}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}/issues
AppUpdatesURL={#MyAppURL}/releases
DefaultDirName={localappdata}\Programs\{#MyAppName}
DefaultGroupName={#MyAppName}
AllowNoIcons=yes
OutputDir=.\output
OutputBaseFilename=gowtree-setup-{#MyAppVersion}
SetupIconFile=..\assets\icon.ico
WizardImageFile=..\assets\wizard-large.bmp
WizardSmallImageFile=..\assets\wizard-small.bmp
WizardStyle=modern
WizardSizePercent=115
PrivilegesRequired=lowest
DisableProgramGroupPage=yes
Compression=lzma2/ultra64
SolidCompression=yes
MinVersion=10.0.22000
UninstallDisplayIcon={app}\{#MyAppExeName}
LicenseFile=..\LICENSE
InfoBeforeFile=welcome.txt
SetupLogging=yes
VersionInfoVersion=1.4.0.0
VersionInfoCompany={#MyAppPublisher}
VersionInfoDescription={#MyAppName} Setup — Directory Tree for Windows
VersionInfoCopyright=Copyright (C) 2026 {#MyAppPublisher}
VersionInfoProductName={#MyAppName}
VersionInfoProductVersion={#MyAppVersion}

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[CustomMessages]
english.WelcomeLabel2=Install [name/ver] — a modern directory tree command for Windows.%n%nClick Next to begin. If Windows SmartScreen appears, choose More info, then Run anyway.%n%nNo admin rights required.

[Tasks]
Name: "addToPath"; Description: "Add gowtree to PATH (recommended)"; GroupDescription: "Options:"
Name: "desktopIcon"; Description: "Create a desktop shortcut"; GroupDescription: "Options:"

[Files]
Source: "..\build\gowtree.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENSE"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Comment: "Directory tree viewer"
Name: "{group}\Uninstall {#MyAppName}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopIcon

[Registry]
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Tasks: addToPath; Check: NeedsAddPath(ExpandConstant('{app}'))

[Run]
Filename: "{app}\{#MyAppExeName}"; Parameters: "-h"; Description: "Launch gowtree help"; Flags: nowait postinstall skipifsilent shellexec

[Code]
function NeedsAddPath(Param: string): Boolean;
var OrigPath: string;
begin
  if not RegQueryStringValue(HKCU, 'Environment', 'Path', OrigPath) then
    Result := True
  else
    Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;

procedure RemoveFromPath(AppDir: string);
var OrigPath, NewPath: string; P: Integer;
begin
  if not RegQueryStringValue(HKCU, 'Environment', 'Path', OrigPath) then Exit;
  P := Pos(';' + AppDir, OrigPath);
  if P > 0 then
    NewPath := Copy(OrigPath, 1, P - 1) + Copy(OrigPath, P + Length(';' + AppDir), MaxInt)
  else if Pos(AppDir + ';', OrigPath) = 1 then
    NewPath := Copy(OrigPath, Length(AppDir) + 2, MaxInt)
  else Exit;
  RegWriteExpandStringValue(HKCU, 'Environment', 'Path', NewPath);
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then
    RemoveFromPath(ExpandConstant('{app}'));
end;

procedure InitializeWizard();
begin
  if WizardForm.TasksList.Items.Count > 0 then
    WizardForm.TasksList.Checked[0] := True;
end;

procedure CurPageChanged(CurPageID: Integer);
begin
  if CurPageID = wpFinished then
    WizardForm.FinishedLabel.Caption :=
      'gowtree installed successfully!' + #13#10 +
      'Open a NEW terminal and type:  gowtree' + #13#10#13#10 +
      'Try:  gowtree -a -s -L 2 --icons';
end;

[Messages]
FinishedLabel=Click Finish to close the installer.
