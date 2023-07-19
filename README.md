# Guarde

Guarde is a TCP and UDP reverse proxy written in Golang. It was created in an attempt to mitigate the potential of a 
DNS amplification attack by limiting the IPs to those that has specific WHOIS properties, this was created more towards 
protecting my local AdGuard Home instance.

## How does this protect Home DNS servers?

The philosophy behind Guarde is that we only allow addresses that come specifically those that you know such as 
ISPs that shouldn't be used by an enterprise (e.g. PLDT ISP for PH). It is designed this way due to the fact that 
most home routers and devices do not have a static IP address (or so is the case in the Philippines).

## How to configure Guarde

Guarde has a very simple configuration:
- `proxy`: contains all the proxy configurations
  - `udp`: contains all the udp configuration details
  - `tcp`: contains all the tcp configuration details
    - `[udp/tcp].forward`: forwards all requests to the given port to this address.
    - `[udp/tcp].port`: the port that Guarde should listen on.
  - `ruleset`: contains all the WHOIS rulesets that matter, please refer to [`Rulesets`](#rulesets).
  - `verbose`: whether to show all the request body and response body.

An example of this configuration would be:
```yaml
proxy:
  udp:
    forward: 172.17.0.1:1053
    port: 53
  tcp:
    forward: 172.17.0.1:1053
    port: 53
ruleset:
  - person: '%PLDT'
    e-mail: '%pldt.com.ph'
    country: 'PH'
verbose: true
```

You have to save the file under the `%HOME%/.guarde/config.yml` file. If you are using Docker, you can link it to the 
`/root/.guarde/config.yml` path since the Dockerfile here uses Alpine Linux.

## Rulesets

To know what ruleset properties are available, you can query your IP Address using [`ARIN Lookup`](https://mxtoolbox.com/arin.aspx). We 
recommend using the CLI tool, or using the link we provided since it provides the raw output, for example, to make sure that the organization 
id matches, you can add this field:
```yaml
OrgId: 'AT-88-Z'
```

If you want to make sure that the field **contains** a specific value, then you can add the `%` prefix at the value and it will 
only ensure that the field contains that specific value.