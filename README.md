# Crony

Simple cron replacement that calls URL

## Usage

Create `crony.json` in same folder with crony binary.

```
{
  "tasks": [
    {
      "cron": "@every 1m",
      "url": "http://google.com"
    }
  ]
}

```  


## License

MIT