# teleport: "прозрачное прокси"

## Генерируем сертификаты если необходимо использовать ssl

```
openssl req \
    -x509 \
    -nodes \
    -newkey rsa:2048 \
    -keyout server.key \
    -out server.crt \
    -days 3650 \
    -subj "/C=RU/ST=Moscow/L=Moscow/O=Global Security/OU=IT Department/CN=*"
```

## Запуск

```bash
./teleport \
    --proxy socks5://user:password@ip:port \
    --listen-http 127.0.0.127:80 \
    --listen-https 127.0.0.127:443 \
    --ssl-key ./example/server.key \
    --ssl-crt ./example/server.crt
```

## Прописываем в hosts

```bash
echo "api.telegram.org 127.0.0.127" >> /etc/hosts
echo "web.telegram.org 127.0.0.127" >> /etc/hosts
```
