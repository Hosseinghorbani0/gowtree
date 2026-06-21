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
WizardStyle=modern
WizardSizePercent=110
PrivilegesRequired=lowest
DisableProgramGroupPage=yes
Compression=lzma2/ultra64
SolidCompression=yes
MinVersion=10.0
UninstallDisplayIcon={app}\{#MyAppExeName}
LicenseFile=..\LICENSE
InfoBeforeFile=..\installer\welcome.txt
SetupLogging=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "farsi"; MessagesFile: "compiler:Translations\Farsi.isl"

[CustomMessages]
english.WelcomeLabel2=This wizard will install [name/ver] on your computer.%n%n🌳 A modern tree command for Windows — colors, icons, JSON/Markdown/HTML export, and more.%n%nNo admin rights required.
farsi.WelcomeLabel2=این برنامه [name/ver] را روی سیستم شما نصب می‌کند.%n%n🌳 یک دستور tree مدرن برای ویندوز — رنگ، آیکون، خروجی JSON/Markdown/HTML و بیشتر.%n%nنیازی به دسترسی Administrator نیست.

[Tasks]
Name: "addToPath"; Description: "Add gowtree to PATH (recommended)"; GroupDescription: "Environment:"; Flags: checked
Name: "desktopIcon"; Description: "Create a desktop shortcut"; GroupDescription: "Shortcuts:"; Flags: unchecked

[Files]
Source: "..\gowtree.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\README.md"; DestDir: "{app}"; Flags: ignoreversion isreadme
Source: "..\LICENSE"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Parameters: "-h"; Comment: "Directory tree viewer"
Name: "{group}\Uninstall {#MyAppName}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopIcon

[Registry]
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Tasks: addToPath; Check: NeedsAddPath(ExpandConstant('{app}'))

[Run]
Filename: "{app}\{#MyAppExeName}"; Parameters: "-h"; Description: "Show gowtree help"; Flags: nowait postinstall skipifsilent shellexec

[UninstallRun]
Filename: "cmd.exe"; Parameters: "/c echo PATH updated — open a new terminal"; Flags: runhidden

[Code]
function NeedsAddPath(Param: string): Boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKCU, 'Environment', 'Path', OrigPath) then
    Result := True
  else
    Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;

procedure RemoveFromPath(AppDir: string);
var
  OrigPath, NewPath: string;
  P: Integer;
begin
  if not RegQueryStringValue(HKCU, 'Environment', 'Path', OrigPath) then
    Exit;
  P := Pos(';' + AppDir, OrigPath);
  if P > 0 then
    NewPath := Copy(OrigPath, 1, P - 1) + Copy(OrigPath, P + Length(';' + AppDir), MaxInt)
  else if Pos(AppDir + ';', OrigPath) = 1 then
    NewPath := Copy(OrigPath, Length(AppDir) + 2, MaxInt)
  else
    Exit;
  RegWriteExpandStringValue(HKCU, 'Environment', 'Path', NewPath);
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then
    RemoveFromPath(ExpandConstant('{app}'));
end;

procedure CurPageChanged(CurPageID: Integer);
begin
  if CurPageID = wpFinished then
    WizardForm.FinishedLabel.Caption :=
      '🌳 gowtree is ready!' + #13#10 +
      'Open a new terminal and run: gowtree' + #13#10#13#10 +
      'Tip: try  gowtree -a -s -L 2 --icons';
end;

function InitializeSetup(): Boolean;
begin
  Result := True;
end;

function NextButtonClick(CurPageID: Integer): Boolean;
begin
  Result := True;
end;

[Messages]
FinishedLabel=Click Finish to close the installer.
