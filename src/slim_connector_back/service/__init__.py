import os

import dotenv
from fastapi import Depends



class EnvService:
    def __init__(self):
        dotenv.load_dotenv("./.env.local")
        dotenv.load_dotenv()
        self.cors_list = self.get_str_list_or_empty("CORS_LIST")
        self.secret_key = os.getenv("SECRET_KEY")
        self.algorithm = self.get_or("ALGORITHM", "HS256")
        self.access_token_expire_minutes = self.get_int_or("ACCESS_TOKEN_EXPIRE_MINUTES", 60 * 24 * 14)
        self.refresh_token_expire_minutes = self.get_int_or("REFRESH_TOKEN_EXPIRE_MINUTES", 15)

    # noinspection PyMethodMayBeStatic
    def get_or_none(self, key: str) -> str | None:
        return os.getenv(key)

    def get_str(self, key: str) -> str:
        value = self.get_or_none(key)
        if value is None:
            raise ValueError(f"{key} is None")
        return value

    def get_or(self, key: str, default: str) -> str:
        value = self.get_or_none(key)
        if value is None:
            return default
        return value

    def get_int_or_none(self, key: str) -> int | None:
        value = self.get_or_none(key)
        if value is None:
            return None
        return int(value)

    def get_int_or(self, key: str, default: int) -> int:
        value = self.get_int_or_none(key)
        if value is None:
            return default
        return value

    def get_str_list_or_empty(self, key: str) -> list[str] | None:
        value = self.get_or_none(key)
        if value is None:
            value = ""
        return value.split(",")


class DbEnvService:

    def __init__(self, env: EnvService = Depends(EnvService)):
        db_user = env.get_str("DB_USER")
        db_pass = env.get_str("DB_PASS")
        db_port = env.get_int_or("DB_PORT", 27017)
        db_host = env.get_or("DB_HOST", "localhost")
        host_port = f"{db_host}:{db_port}"
        self.db_name = env.get_or("DB_NAME", "slim-connector")
        self.db_url = env.get_or("DB_URL", f"mongodb://{db_user}:{db_pass}@{host_port}/")


from starlette.middleware.cors import CORSMiddleware


class CorsService:
    def __init__(self, env: EnvService = Depends(EnvService)):
        from slim_connector_back import app
        app.add_middleware(
            CORSMiddleware,
            allow_origins=env.cors_list,
            allow_credentials=True,
            allow_methods=["*"],
            allow_headers=["*"],
        )


class DbService:
    pass
