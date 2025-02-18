# Environment Variables

This document provides a description of the environment variables used in the application.

## Environment Variables

1. **REMNAWAWE_URL**  
   _Description:_ The base URL for the REMNAWAWE service.  
   _Example:_ `REMNAWAWE_URL=domain`

2. **APP_PORT**  
   _Description:_ The port on which the application will run.  
   _Example:_ `APP_PORT=4000`

3. **V2RAY_TEMPLATE_PATH**  
   _Description:_ The file path to the default V2Ray configuration template.  
   _Example:_ `V2RAY_TEMPLATE_PATH=/app/templates/v2ray/default.json`

4. **V2RAY_MUX_ENABLED**  
   _Description:_ A flag to enable or disable the V2Ray Mux feature. Set to `true` to enable Mux.  
   _Example:_ `V2RAY_MUX_ENABLED=true`

5. **V2RAY_MUX_TEMPLATE_PATH**  
   _Description:_ The file path to the V2Ray Mux configuration template.  
   _Example:_ `V2RAY_MUX_TEMPLATE_PATH=/app/templates/v2ray/mux_default.json`

---

# How to Run

1. Clone the repo
```bash
git clone https://github.com/Jolymmiles/remnawawe-json
```

2. Go to the cloned repo
```bash
cd remnawawe-json
```

3. Run Docker Compose
```bash
docker compose up -d
```

