$dest = "$env:APPDATA\terraform.d\plugins\local\taskmanager\taskmanager\0.1.0\windows_amd64"
New-Item -ItemType Directory -Force -Path $dest

Move-Item -Path ".\terraform-provider-taskmanager.exe" -Destination "$dest\terraform-provider-taskmanager.exe" 

C:\Users\ShivrajTelange\AppData\Roaming\terraform.d\plugins\example.com\shivraj\taskmanager\0.1.0\windows_amd64

Remove-Item -Path "$env:APPDATA\terraform.d\plugins\example.com\shivraj\taskmanager\0.1.0\windows_amd64" -Recurse -Force


