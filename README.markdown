# pushprovider

This is a simple client for [Apple’s new APNs Provider API](https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/APNsProviderAPI.html) for sending push notifications. It's designed to be easily scripted. You can test it from a terminal like this:

```
> pushprovider cert.p12
{"token":"4b62…","payload":{"aps":{"alert":"Hey there!"}}}
{"token":"3d9f…","payload":{"aps":{"alert":"Hey to you, too."}}}
^D
```

For each line of JSON you send, pushprovider delivers the notification to Apple and responds with a line like this:

```
{"status":200}
```

Or this:

```
{"status":400,"body":{"reason":"BadDeviceToken"}}
```

You can find more information what status codes and bodies to expect in Apple’s documentation, linked above.

When you're finished, close stdin and pushprovider will exit.

You must provide the path to a `.p12` file containing a push certificate and private key. You can generate this by selecting the private key and and certificate in Keychain Access and exporting them (make sure you’re in “All Items” view, when you do this, *not* “My Certificates” — or else Keychain Access appears to include extra stuff in the export which pushprovider won’t know what to do with). 

You can also pass an optional command line flag, `-d`, which connects to the development push environment instead of production.

