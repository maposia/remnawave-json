# Environment Variables

This document provides a description of the environment variables used in the application.

## Environment Variables

1. **REMNAWAWE_TOKEN**  
   _Description:_ The token used for authentication or access to the REMNAWAWE service.  
   _Example:_ `REMNAWAWE_TOKEN=token`

2. **REMNAWAWE_URL**  
   _Description:_ The base URL for the REMNAWAWE service.  
   _Example:_ `REMNAWAWE_URL=domain`

3. **APP_PORT**  
   _Description:_ The port on which the application will run.  
   _Example:_ `APP_PORT=4000`

4. **V2RAY_TEMPLATE_PATH**  
   _Description:_ The file path to the default V2Ray configuration template.  
   _Example:_ `V2RAY_TEMPLATE_PATH=/app/templates/v2ray/default.json`

# How to run

1. Clone repo
```git clone https://github.com/Jolymmiles/remnawawe-json```
2. Go to cloned repo
``` cd remnawawe-json```
3. Run docker compose
```docker compose up -d```