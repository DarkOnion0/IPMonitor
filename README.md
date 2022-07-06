# IPMonitor

[![Go](https://github.com/DarkOnion0/IPMonitor/actions/workflows/go.yml/badge.svg)](https://github.com/DarkOnion0/IPMonitor/actions/workflows/go.yml) [![Container](https://github.com/DarkOnion0/IPMonitor/actions/workflows/container.yml/badge.svg)](https://github.com/DarkOnion0/IPMonitor/actions/workflows/container.yml) [![Latest release](https://shields.io/github/v/release/DarkOnion0/IPMonitor?display_name=tag&include_prereleases&label=%F0%9F%93%A6%20Latest%20release)](https://shields.io/github/v/release/DarkOnion0/IPMonitor?display_name=tag&include_prereleases&label=%F0%9F%93%A6%20Latest%20release)

IPMonitor is a little app that warn you when your public IP changed (the IP of the internet connection where the app is installed). It can be very useful when you don't use DNS record and [Dynamic DNS](https://en.wikipedia.org/wiki/Dynamic_DNS), just a raw IP to connect to your selfhosted app. But by using just an IP, if it changes and that your not at home/can't have a way to see your new IP, you just can't connect to your app anymore. Here come the real job of IPMonitor, every 10 minutes (the value can be changed) the app check your IP and print the result in the console, now you just need to use a log collector like Promtail, FluentBIT... or just a bash script to parse the log and trigger an alert that will send you a message with your new IP.

## 🚀 Main Features

- Monitor IP
- Has multiple modes to monitor IP (Cron/API)
- Fast
- Self-contained

## 📖 Usage

### 📦️ Providers

<details>
  <summary>
    Docker (recommended)
  </summary>
  <p>
  > **⚠️ The `latest` tag follow the master branch so it may be unstable or just not-working, USE IT AT YOUR OWN RISKS ⚠️**
  1. Download the container from GitHub
  2. Run it, further configuration can be done, see the corresponding sections bellow.

    docker pull ghcr.io/darkonion0/ipmonitor:latest
    docker run ghcr.io/darkonion0/ipmonitor:latest

  </p>
</details>

<details>
  <summary>
    Binary
  </summary>
  <p>
  
  1. Download the binary from the release page
  2. Run it, further configuration can be done, see the corresponding sections below.  
     Also don't forget to change the binary name according to the plateforme where you're one and the selected version

```bash
 ❯ ./ipmonitor-linux-amd64-latest
```

  </p>
</details>

### 🔄 Modes

> NOTE: both modes can be run at the same time without any problem

#### API Mode

Just run the command and replace `localhost:8080` according to the server's ip where you deployed the app and its port

```bash
❯ curl localhost:8080 | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                  Dload  Upload   Total   Spent    Left  Speed
100    76  100    76    0     0   7428      0 --:--:-- --:--:-- --:--:--  7600
{
  "CurrentIP": "1.1.1.1",
  "PreviousIP": "1.1.1.1",
  "IpChanged": false
}
```

#### Cron Mode

Run the container/binary then just wait, and see the magic (by default the cron job is running every 10 minutes but is also triggered at startup time).

## 🧰 Configuration

<details>
  <summary>
    Docker (recommended)
  </summary>
  <p>
  
  Since the app doesn't write log to a file, Docker will allow you to easily capture log with any log collector (FluentD, Logstach...)  
  All the config is done trough the following env var

```dockerfile
DEBUG "false" # Enable debug mode, make the log very verbose and pretty (but not fast 😉)
CRON_ENABLE true # Enable/Disable cron mode
API_ENABLE true # Enable/Disable API mode

# All the settings for the cron, they are split up since Docker messed them when it pass them to the binary inside the container
MINUTES "*/10"
HOURS "*"
MONTH_DAY "*"
MONTH "*"
WEEK_DAY "*"
```

  </p>
</details>

<details>
  <summary>
    Binary
  </summary>
  <p>
  
  All the configuration of the binary are done trough flags, here is the list:
  
  ```
  -api-enable string
      This flag enable the API mode, it can be disable to run it in cron mode only (default "true")
  -api-port string
      Set a custom api's listen port (default "8080")
  -cron string
      Set a custom cron scheduled to run the IP check, run every 15 minutes by default (default "*/15 * * * *")
  -cron-enable string
      This flag enable the cron mode, it can be disable to it in API mode only (default "true")
  -debug string
      Sets log level to debug (default "false")
  ```

  </p>
</details>
