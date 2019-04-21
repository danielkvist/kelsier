# Kelsier

> This project has been inspired by [this article](https://dev.to/healeycodes/build-a-python-bot-to-find-your-website-s-dead-links-563c).

Kelsier is a simple CLI to check for dead links on a specified URL.

## Example

> The default URL is https://www.google.com

```bash
kelsier
200 - https://www.google.com/intl/es/policies/terms/
200 - https://www.google.com/intl/es/policies/privacy/
200 - https://www.google.com/preferences?hl=es
200 - https://www.google.com/advanced_search?hl=es&authuser=0
...
```

Or

```bash
kelsier -url https://dkvist.com
200 - https://dkvist.com/blog/index.xml
200 - https://www.dkvist.com/css/
200 - https://www.dkvist.com/favicon.png
200 - https://www.dkvist.com#contact
200 - https://www.dkvist.com/blog/
...
```

## Install

```bash
go get github.com/danielkvist/kelsier
```

## Get Docker Image

```bash
docker image pull danielkvist/kelsier
```