### 📖 V2Ray Json Subscription Proxy with User-Agent Routing

This application serves as a **proxy** for generating and serving **V2Ray subscription configurations** based on
`User-Agent`. It dynamically provides the correct configuration format depending on the client version, ensuring
seamless compatibility across different applications.

Work with <https://remna.st>

For use your own web page template, use {{.MetaTitle}} {{.MetaDescription}} {{.PanelData}} to get panel data and other.

## 🇷🇺 [Happ Routing](https://github.com/hydraponique/roscomvpn-happ-routing/tree/main)

## ✨ Features

- **🔀 User-Agent-Based Routing**
    - Automatically detects and serves the correct subscription format for supported clients:
        - **Streisand** (JSON)
        - **Happ** (JSON)
- **Web page template**
    - Supported web page template.
- **🌍 Direct Proxy Fallback**
    - If `User-Agent` is unsupported, the server provides a **default proxy response**.

---

## ⚙️ Configuration

Modify `.env.sample` to adjust the application settings:

```
REMNAWAVE_URL=sub_domain
APP_PORT=4000
# V2RAY_TEMPLATE_PATH=/app/templates/v2ray/default.json
# V2RAY_MUX_ENABLED=true
# V2RAY_MUX_TEMPLATE_PATH=/app/templates/v2ray/mux_default.json
# WEB_PAGE_TEMPLATE_PATH=/app/templates/subscription/index.html
META_TITLE=Zalupa
META_DESCRIPTION=Pupa
REMNAWAVE_TOKEN=
MODE=
```

After modifying execute this

```bash
mv .env.sample .env
```

---

# How to Run

1. Clone the repo

```bash
git clone https://github.com/Jolymmiles/remnawave-json
```

2. Go to the cloned repo

```bash
cd remnawave-json
```

3. Configure .env

4. Run Docker Compose

```bash
docker compose up -d
```

---

# How to update

1. Go to directory with docker-compose.yaml

```bash
remnawave-json
```

2. Update image

```bash
docker compose pull
```

3. Restart container

```bash
docker compose down --remove-orphans && docker compose up -d
```

---

## 🌿 **Environment Variables**

| Variable Name           | Description                                                          | Example Value                            |
|-------------------------|----------------------------------------------------------------------|------------------------------------------|
| REMNAWAVE_URL           | The base URL for the subdomain                                       | `https://panel.com`                      |
| APP_PORT                | The port on which the application will run                           | `4000`                                   |
| V2RAY_TEMPLATE_PATH     | The file path to the default V2Ray configuration template            | `/app/templates/v2ray/default.json`      |
| V2RAY_MUX_ENABLED       | A flag to enable or disable the V2Ray Mux feature                    | `true`                                   |
| V2RAY_MUX_TEMPLATE_PATH | The file path to the V2Ray Mux configuration template                | `/app/templates/v2ray/mux_default.json`  |
| WEB_PAGE_TEMPLATE_PATH  | The file path to the subscription template                           | `/app/templates/subscription/index.html` |
| HAPP_JSON_ENABLED       | A flag to enable or disable JSON output for Happ                     | `false`                                  |
| HAPP_ROUTING            | The routing path for Happ connections                                | `happ://routing/...`                     |
| HAPP_ANNOUNCEMENTS      | Announcement text in plain text                                      | `zalupa`                                 |
| RU_OUTBOUND_NAME        | RU outbound name                                                     | `RU`                                     |
| RU_USER_HOST            | RU user host                                                         | `Россия`                                 |
| REMNAWAVE_TOKEN         | REMNAWAVE token                                                      | `zalupa`                                 |
| MetaDescription         | MetaDescription                                                      | `Zalupa`                                 |
| MetaTitle               | MetaTitle                                                            | `Zalupa`                                 |
| MODE                    | Set if using remnawave:3000                                          | `local`                                  |
| EXCEPT_RU_RULES_USERS   | Set subscription short uuid for exclude routing via RU_OUTBOUND_NAME | `c11JfduMqrkBZrTZ`                       |

---

## Nginx example

```nginx configuration
server
{
        listen 443 ssl;
        listen [::]:443 ssl;
        http2 on;

        #from .env remnawave SUB_PUBLIC_DOMAIN
        server_name sub_domain ;


        # you certs
        ssl_certificate /root/.acme.sh/example.com_ecc/example.com.cer;
        ssl_certificate_key /root/.acme.sh/example.com_ecc/example.com.key;
        ssl_dhparam /etc/nginx/ssl/dhparams.pem;

        ssl_protocols TLSv1.3;
        ssl_ciphers TLS13-CHACHA20-POLY1305-SHA256:TLS13-AES-256-GCM-SHA384:TLS13-AES-128-GCM-SHA256:EECDH+CHACHA20:EECDH+AESGCM:EECDH+AES;
        ssl_prefer_server_ciphers on;

        # HSTS
        add_header Strict-Transport-Security "max-age=63072000" always;

        # OCSP stapling
        ssl_stapling on;
        ssl_stapling_verify on;
        ssl_trusted_certificate /root/.acme.sh/example.com_ecc/ca.cer;

        location / {
        proxy_http_version 1.1;
        # APP_PORT
        proxy_pass http://127.0.0.1:4000;
        proxy_set_header Host $host;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Forwarded-Port $server_port;

        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
 }
```

## 📜 License

This project is open-source and available under the **AGPL v2.0**.

🚀 _Enhance your V2Ray subscription management with automated User-Agent routing!_
