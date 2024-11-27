import os

import dotenv

from slim_connector_back.util import urls

dotenv.load_dotenv("./.env.local")
dotenv.load_dotenv()


class Token:
    refresh_token_expire_minutes: int | float
    access_token_expire_minutes: int | float
    secret_key = os.getenv("SECRET_KEY")
    algorithm: str

    def __init__(self):
        refresh_token_expire_minutes = os.getenv("REFRESH_TOKEN_EXPIRE_MINUTES")
        if refresh_token_expire_minutes is None:
            self.refresh_token_expire_minutes = 60 * 24 * 14
        else:
            self.refresh_token_expire_minutes = float(refresh_token_expire_minutes)

        access_token_expire_minutes = os.getenv("ACCESS_TOKEN_EXPIRE_MINUTES")
        if access_token_expire_minutes is None:
            self.access_token_expire_minutes = 15
        else:
            self.access_token_expire_minutes = float(access_token_expire_minutes)

        algorithm = os.getenv("ALGORITHM")
        if algorithm is None:
            self.algorithm = "HS256"
        else:
            self.algorithm = algorithm

        secret_key = os.getenv("SECRET_KEY")
        if secret_key is None:
            raise ValueError
        else:
            self.secret_key = secret_key


class Database:
    db_url = os.getenv("DB_URL")
    db_user = os.getenv("DB_USER")
    db_pass = os.getenv("DB_PASS")
    db_host = os.getenv("DB_HOST")
    db_port = os.getenv("DB_PORT")
    db_name = os.getenv("DB_NAME")

    def __init__(self):
        if self.db_url is not None:
            return
        if self.db_user is None:
            raise ValueError("DB_USER must not be None")
        if self.db_pass is None:
            raise ValueError("DB_PASS must not be None")
        if self.db_host is None:
            self.db_host = "localhost"
        if self.db_port is not None:
            self.db_host = f"{self.db_host}:{self.db_port}"
        if self.db_name is None:
            self.db_name = "hew-dev"
        if self.db_url is None:
            self.db_url = f"mongodb://{self.db_user}:{self.db_pass}@{self.db_host}/"



class Env:
    cors_list = os.getenv("CORS_LIST")
    img_url = urls.URL.by_str(os.getenv("IMG_URL"))
    token = Token()
    database = Database()


ENV = Env()
