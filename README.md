# drone-monitor
A hodge-podge of scripts that were a part of the setup I ran on my drone for each flight and afterwards for data processing as a part of my MEng project. 

Before launching either the drone setup or the server setup, build all modules. 

# Server flight setup
To launch the data collection setup on the server, go to ```scripts/startup/server``` and run 

```
go build
./server
```

# Drone flight setup
To launch the data collection setup on the drone, go to ```scripts/startup/drone``` and run

```
go build
./drone -tcp
```

# Uploading data
The data processing setup is in ```scripts/upload```. Here, all the data collected on the server is processed and uploaded to Dropbox. 

To use this, install dropbox uploading script from https://github.com/andreafabrizi/Dropbox-Uploader and run 

```
go build 
./upload [args]
```

The uploader supports processing and uploading different types of log files collected and allows for uploading batches of data at once. 
