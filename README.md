# Notifsaver
## Config path: 
./config.yml

## Config structure:
```
# config for s3
s3:

 id: ""
 key: ""
 token: ""
 region: ""
 bucket: ""
# config for kafka
consumer:
 addr: "kafka:9092"
 notif_topic: "mt-notifcation"
 client_id: "go-notification-consumer"
