# chatrabbit

A openai proxy program writen by Golang.



This prgoram provides a api proxy for [OpenAI API](https://platform.openai.com/). 

We support:
* ChatGPT
* GPT-3, GPT-4

## setup

copy conf/config.yaml.example to conf/config.yaml

### launch

```bash
go run cmd/proxy/main.go -config conf/config.yaml
```

对于VSCode，可以设置环境变量：

```json
    "env": {
        "CHAT_RABBIT_CONF": "../../conf/config.yaml"
    }
```