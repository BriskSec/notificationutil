# notificationutil

Package notificationutil implements utility that should be used across all BriskSec products when it's required to
send a notification to the central monitoring panel. This could include informational details, new detections, or
any abnormal or error conditions that should be alerted to the administration team.

Install: 

```
go get github.com/brisksec/notificationutil
```

Usage for non-error conditions: 
```go
notificationutil.Notify("Title for the notification", notificationutil.AbnormalCondition, "Contents of the message", nil)
```

Usage for error conditions: 
```go
notificationutil.Notify("Title for the notification", notificationutil.AbnormalCondition, "Contents of the message", err)
```
