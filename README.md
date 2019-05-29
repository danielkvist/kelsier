# Kelsier

> This project has been inspired by [this article](https://dev.to/healeycodes/build-a-python-bot-to-find-your-website-s-dead-links-563c).

Kelsier is a simple terminal tool that allows you to find broken links on one or more Web pages.

## Example

```bash
kelsier https://google.com
200 - https://google.com/intl/es/policies/terms/
200 - https://google.com/intl/es/policies/privacy/
200 - https://google.com/preferences?hl=es
200 - https://google.com/advanced_search?hl=es&authuser=0
...
```

Or

```bash
kelsier https://dkvist.com https://github.com
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

Or

```bash
go install github.com/danielkvist/kelsier
```

## Docker Image

```bash
docker image pull danielkvist/kelsier
```

## Help

If you think there's something that can be improved. Please don't hesitate to let me know.