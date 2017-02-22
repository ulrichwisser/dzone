# dzone
Take statistics of IIS Anycast and save to InfluxDB


## Installation

```
$ go get -u github.com/ulrichwisser/dzone
```

## Configuration
Dzone will read first $HOME/.dzone then it will read ./.dzone. Next the config file given at the command line (if any) will be read and finally the command line arguments will be parsed. All these configurations will be joined together. Information which is read later overwrites any information from previous configuration.

The configuration files have to be in YAML format.
```
apiuser: apiusername
apisecret: apisecret
apipasswd: password
serverroot: https://api.anycast.iis.se
influxserver: https://127.0.0.1:8086
influxdb: databasename
influxuser: username
influxpasswd: password
```

## Command Line Parameters
```
--dryrun                     run all statistics but do not write to InfluxDB (write data to STDOUT instead)
--conf <filename>            file to read configuration
--serverRoot                 the base URL of the RestAPI
--influxServer <server>      name or ip of the server running InfluxDB
--influxDB <dbname>          name of the database to save statistics to
--influxUser <username>      username for authorization to InfluxDB
--influxPasswd <password>    password for authorization to InfluxDB
```
## Limitations
- Untested for registrars and zone owners 
