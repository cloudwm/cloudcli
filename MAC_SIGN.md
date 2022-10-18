# Signing the binary for Mac

Following procedure should be used to allow signing the binary for Mac

* create an AWS dedicated host type mac1.metal
* create an AWS instance:
  * macOS Monterey 12.5.1 AMI built by Amazon Web Services
  * mac1.metal
  * 100GB storage
  * Tenancy: dedicated host
  * target host by Host ID: the one you created above
* SSH to the instance using the key you selected and username ec2-user
  * enable access via vnc:
```
sudo /System/Library/CoreServices/RemoteManagement/ARDAgent.app/Contents/Resources/kickstart \
    -activate -configure -access -on \
    -configure -allowAccessFor -specifiedUsers \
    -configure -users ec2-user \
    -configure -restart -agent -privs -all
sudo /System/Library/CoreServices/RemoteManagement/ARDAgent.app/Contents/Resources/kickstart \
    -configure -access -on -privs -all -users ec2-user
sudo passwd ec2-user
```
* login to the machine using VNC with username ec2-user and the password you set above
* login to the mac os using the same password
* Acquiring a Developer ID Certificate
  * login to https://developer.apple.com
  * Navigate to the certificates page.
  * Click the “+” icon, select “Developer ID Application” and follow the steps.
  * Download the certificate so it's available on the mac instance
* in the mac instance VNC double-click the certificate to import it into the system keychain
* To verify you did this correctly, you can inspect your keychain:
```
$ security find-identity -v
1) 4194587FE60D93D416CF3F4669FF913C7BBA4271 "Developer ID Application: Your Name (GK80BB2A7)"
   1 valid identities found
```
* SSH to the mac instance and create the following file at the home directory `gon-config.json`:
```
{
    "source" : ["./cloudcli"],
    "bundle_id" : "com.kamatera.cloudcli",
    "apple_id": {
        "username": "@env:AC_USERNAME",
        "password":  "@env:AC_PASSWORD"
    },
    "sign" :{
        "application_identity" : "4194587FE60D93D416CF3F4669FF913C7BBA4271"
    },
    "zip" :{
        "output_path" : "./cloudcli.zip"
    }
}
```
* Install gon: `brew install mitchellh/gon/gon`
